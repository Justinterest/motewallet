package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	dtoreq "motewallet/internal/dto/request"
	dtoresp "motewallet/internal/dto/response"
	"motewallet/internal/model"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/kun"
	kundto "motewallet/internal/pkg/kun/dto"
	"motewallet/internal/pkg/utils"
	"motewallet/internal/repository"
)

type TransferService struct {
	db                    *gorm.DB
	merchantRepo          repository.MerchantRepository
	walletSvc             *WalletService
	transferOrderRepo     repository.TransferOrderRepository
	transactionRecordRepo repository.TransactionRecordRepository
	kunClient             kun.KUNClient
	currencyConfigSvc     *CurrencyConfigService
}

func NewTransferService(
	db *gorm.DB,
	merchantRepo repository.MerchantRepository,
	walletSvc *WalletService,
	transferOrderRepo repository.TransferOrderRepository,
	transactionRecordRepo repository.TransactionRecordRepository,
	kunClient kun.KUNClient,
	currencyConfigSvc *CurrencyConfigService,
) *TransferService {
	return &TransferService{
		db:                    db,
		merchantRepo:          merchantRepo,
		walletSvc:             walletSvc,
		transferOrderRepo:     transferOrderRepo,
		transactionRecordRepo: transactionRecordRepo,
		kunClient:             kunClient,
		currencyConfigSvc:     currencyConfigSvc,
	}
}

func (s *TransferService) Transfer(ctx context.Context, merchantID uint64, req *dtoreq.TransferReq) (uint64, error) {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, bizerrors.ErrNotFoundError
		}
		return 0, bizerrors.ErrInternalError
	}

	if merchant.KunSubCustomerNo == nil {
		return 0, bizerrors.ErrMerchantNotRegisteredE
	}

	if err := s.currencyConfigSvc.EnsureCurrencySupported(ctx, merchant, req.Currency); err != nil {
		return 0, err
	}

	if err := validateTransferAccountType(req.FromAccountType); err != nil {
		return 0, err
	}
	if err := validateTransferAccountType(req.ToAccountType); err != nil {
		return 0, err
	}

	if req.FromAccountType == req.ToAccountType {
		return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "from and to account types must be different")
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "invalid amount")
	}

	fromAcc, err := kun.PlatformAccountToKUN(req.FromAccountType)
	if err != nil {
		return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "invalid from account type")
	}
	toAcc, err := kun.PlatformAccountToKUN(req.ToAccountType)
	if err != nil {
		return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "invalid to account type")
	}

	requestNo := kun.GenerateRequestNo()
	subType := "INTERNAL"
	var orderID uint64
	var txRecordID uint64

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.walletSvc.FreezeBalance(ctx, tx, merchantID, req.FromAccountType, req.Currency, amount); err != nil {
			return err
		}

		platformOrderID := utils.GeneratePlatformOrderID("TR")
		txRecord := &model.TransactionRecord{
			PlatformOrderID: platformOrderID,
			MerchantID:      merchantID,
			Type:            "TRANSFER",
			SubType:         &subType,
			Amount:          amount,
			Currency:        req.Currency,
			PlatformFee:     decimal.Zero,
			Status:          "PROCESSING",
		}
		if err := tx.WithContext(ctx).Create(txRecord).Error; err != nil {
			return err
		}

		transferOrder := &model.TransferOrder{
			TransactionRecordID: txRecord.ID,
			MerchantID:          merchantID,
			FromAccountType:     req.FromAccountType,
			ToAccountType:       req.ToAccountType,
			KunRequestNo:        &requestNo,
		}
		if err := tx.WithContext(ctx).Create(transferOrder).Error; err != nil {
			return err
		}

		orderID = transferOrder.ID
		txRecordID = txRecord.ID
		return nil
	})
	if err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			return 0, bizErr
		}
		slog.Error("create transfer order failed", slog.Any("error", err))
		return 0, bizerrors.ErrInternalError
	}

	var kunResp kundto.FundTransferResp
	err = s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, kun.FundTransferPath, &kundto.FundTransferReq{
		RequestNo: requestNo,
		FromAcc:   fromAcc,
		ToAcc:     toAcc,
		Currency:  req.Currency,
		Amount:    req.Amount,
	}, &kunResp)
	if err != nil {
		slog.Error("KUN fund transfer failed", slog.Any("error", err))
		s.rollbackTransfer(ctx, txRecordID, merchantID, req.FromAccountType, req.Currency, amount)
		return 0, bizerrors.ErrKUNAPIFailedE
	}

	if err := s.transferOrderRepo.UpdateFields(ctx, orderID, map[string]interface{}{
		"kun_order_id": kunResp.OrderId,
	}); err != nil {
		slog.Error("update transfer kun order id failed", slog.Any("error", err))
	}

	if kunResp.ResolvedStatus() == "SUCCESS" {
		if err := s.SettleFromWebhook(ctx, requestNo, "SUCCESS"); err != nil {
			slog.Error("settle transfer on sync success failed", slog.Any("error", err))
		}
	}

	return orderID, nil
}

func (s *TransferService) SettleFromWebhook(ctx context.Context, requestNo, orderStatus string) error {
	order, err := s.transferOrderRepo.FindByKunRequestNo(ctx, requestNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	txRecord, err := s.transactionRecordRepo.FindByID(ctx, order.TransactionRecordID)
	if err != nil {
		return err
	}

	if txRecord.Status == "COMPLETED" || txRecord.Status == "FAILED" {
		return nil
	}

	amount := txRecord.Amount

	switch orderStatus {
	case "SUCCESS":
		return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := s.walletSvc.DeductFrozen(ctx, tx, order.MerchantID, order.FromAccountType, txRecord.Currency, amount); err != nil {
				return err
			}
			if err := s.walletSvc.CreditBalance(ctx, tx, order.MerchantID, order.ToAccountType, txRecord.Currency, amount); err != nil {
				return err
			}
			return tx.WithContext(ctx).Model(&model.TransactionRecord{}).
				Where("id = ?", txRecord.ID).
				Update("status", "COMPLETED").Error
		})
	case "FAIL":
		return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := s.walletSvc.UnfreezeBalance(ctx, tx, order.MerchantID, order.FromAccountType, txRecord.Currency, amount); err != nil {
				return err
			}
			return tx.WithContext(ctx).Model(&model.TransactionRecord{}).
				Where("id = ?", txRecord.ID).
				Update("status", "FAILED").Error
		})
	default:
		slog.Info("fund transfer status update", slog.String("status", orderStatus))
	}

	return nil
}

func (s *TransferService) rollbackTransfer(ctx context.Context, txRecordID, merchantID uint64, fromAccountType, currency string, amount decimal.Decimal) {
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.walletSvc.UnfreezeBalance(ctx, tx, merchantID, fromAccountType, currency, amount); err != nil {
			return err
		}
		return tx.WithContext(ctx).Model(&model.TransactionRecord{}).
			Where("id = ?", txRecordID).
			Update("status", "FAILED").Error
	}); err != nil {
		slog.Error("rollback transfer failed", slog.Any("error", err), slog.Uint64("tx_record_id", txRecordID))
	}
}

func validateTransferAccountType(accountType string) error {
	switch accountType {
	case "FUNDING", "TRADING":
		return nil
	default:
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "invalid account type")
	}
}

func (s *TransferService) ListTransfers(ctx context.Context, merchantID uint64, page, pageSize int) (*dtoresp.TransferOrderListResp, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	orders, total, err := s.transferOrderRepo.ListByMerchant(ctx, merchantID, page, pageSize)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	var resp []dtoresp.TransferOrderResp
	for _, o := range orders {
		txRecord, _ := s.transactionRecordRepo.FindByID(ctx, o.TransactionRecordID)
		item := dtoresp.TransferOrderResp{
			ID:              o.ID,
			FromAccountType: o.FromAccountType,
			ToAccountType:   o.ToAccountType,
			CreatedAt:       o.CreatedAt,
		}
		if txRecord != nil {
			item.Currency = txRecord.Currency
			item.Amount = txRecord.Amount.String()
			item.Status = txRecord.Status
		}
		resp = append(resp, item)
	}

	return &dtoresp.TransferOrderListResp{
		Orders: resp,
		Total:  total,
	}, nil
}

func (s *TransferService) SyncOrderStatus(ctx context.Context, orderID uint64) (*dtoresp.AdminTransferSyncResp, error) {
	order, err := s.transferOrderRepo.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	txRecord, err := s.transactionRecordRepo.FindByID(ctx, order.TransactionRecordID)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	resp := &dtoresp.AdminTransferSyncResp{
		OrderID: order.ID,
		Status:  txRecord.Status,
		Updated: false,
	}

	if txRecord.Status == "COMPLETED" {
		return resp, nil
	}

	if order.KunRequestNo == nil || strings.TrimSpace(*order.KunRequestNo) == "" {
		return nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "transfer order has no kun request no")
	}

	merchant, err := s.merchantRepo.FindByID(ctx, order.MerchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}
	if merchant.KunSubCustomerNo == nil {
		return nil, bizerrors.ErrMerchantNotRegisteredE
	}

	listReq := kundto.FundTransferListReq{
		RequestNo:         kun.GenerateRequestNo(),
		OriginalRequestNo: strings.TrimSpace(*order.KunRequestNo),
		PageNo:            1,
		PageSize:          20,
	}
	if order.KunOrderID != nil && strings.TrimSpace(*order.KunOrderID) != "" {
		listReq.OrderId = strings.TrimSpace(*order.KunOrderID)
	}

	var kunResp kundto.FundTransferListResp
	err = s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, kun.FundTransferListPath, listReq, &kunResp)
	if err != nil {
		slog.Error("KUN fund transfer list query failed", slog.Any("error", err), slog.Uint64("order_id", orderID))
		return nil, bizerrors.ErrKUNAPIFailedE
	}

	matched := findKUNTransferRecord(kunResp.Records, order, txRecord)
	if matched == nil {
		resp.KunStatus = "NOT_FOUND"
		return resp, nil
	}

	kunStatus := strings.ToUpper(strings.TrimSpace(matched.Status))
	if kunStatus == "" && strings.TrimSpace(matched.TransferTime) != "" {
		kunStatus = "SUCCESS"
	}
	resp.KunStatus = kunStatus

	switch kunStatus {
	case "SUCCESS":
		prevStatus := txRecord.Status
		if err := s.SettleFromWebhook(ctx, *order.KunRequestNo, "SUCCESS"); err != nil {
			return nil, err
		}
		updatedRecord, err := s.transactionRecordRepo.FindByID(ctx, txRecord.ID)
		if err != nil {
			return nil, bizerrors.ErrInternalError
		}
		resp.Status = updatedRecord.Status
		resp.Updated = prevStatus != updatedRecord.Status
	case "FAIL":
		prevStatus := txRecord.Status
		if err := s.SettleFromWebhook(ctx, *order.KunRequestNo, "FAIL"); err != nil {
			return nil, err
		}
		updatedRecord, err := s.transactionRecordRepo.FindByID(ctx, txRecord.ID)
		if err != nil {
			return nil, bizerrors.ErrInternalError
		}
		resp.Status = updatedRecord.Status
		resp.Updated = prevStatus != updatedRecord.Status
	}

	return resp, nil
}

func findKUNTransferRecord(items []kundto.FundTransferListItem, order *model.TransferOrder, txRecord *model.TransactionRecord) *kundto.FundTransferListItem {
	for i := range items {
		item := &items[i]
		if order.KunOrderID != nil && strings.TrimSpace(*order.KunOrderID) != "" {
			if item.OrderId == strings.TrimSpace(*order.KunOrderID) {
				return item
			}
			continue
		}
		if item.Currency != txRecord.Currency {
			continue
		}
		itemAmount, err1 := decimal.NewFromString(item.Amount)
		recordAmount := txRecord.Amount
		if err1 == nil && itemAmount.Equal(recordAmount) {
			return item
		}
	}
	return nil
}
