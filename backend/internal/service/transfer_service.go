package service

import (
	"context"
	"errors"
	"log/slog"

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
}

func NewTransferService(
	db *gorm.DB,
	merchantRepo repository.MerchantRepository,
	walletSvc *WalletService,
	transferOrderRepo repository.TransferOrderRepository,
	transactionRecordRepo repository.TransactionRecordRepository,
	kunClient kun.KUNClient,
) *TransferService {
	return &TransferService{
		db:                    db,
		merchantRepo:          merchantRepo,
		walletSvc:             walletSvc,
		transferOrderRepo:     transferOrderRepo,
		transactionRecordRepo: transactionRecordRepo,
		kunClient:             kunClient,
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

	if req.FromAccountType == req.ToAccountType {
		return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "from and to account types must be different")
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "invalid amount")
	}

	requestNo := kun.GenerateRequestNo()
	var kunResp kundto.FundTransferResp
	err = s.kunClient.Post(ctx, "/rest/v2.0/user/fund/transfer", &kundto.FundTransferReq{
		SubCustomerNo:   *merchant.KunSubCustomerNo,
		RequestNo:       requestNo,
		Currency:        req.Currency,
		Amount:          req.Amount,
		FromAccountType: req.FromAccountType,
		ToAccountType:   req.ToAccountType,
		RegionCode:      s.kunClient.GetRegionCode(),
	}, &kunResp)
	if err != nil {
		slog.Error("KUN fund transfer failed", slog.Any("error", err))
		return 0, bizerrors.ErrKUNAPIFailedE
	}

	subType := "INTERNAL"
	var orderID uint64

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
			KunOrderID:          &kunResp.OrderId,
			KunRequestNo:        &requestNo,
		}
		if err := tx.WithContext(ctx).Create(transferOrder).Error; err != nil {
			return err
		}

		orderID = transferOrder.ID
		return nil
	})

	if err != nil {
		slog.Error("create transfer order failed", slog.Any("error", err))
		return 0, bizerrors.ErrInternalError
	}

	return orderID, nil
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
