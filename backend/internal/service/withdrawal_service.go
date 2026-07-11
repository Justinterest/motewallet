package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

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

type WithdrawalService struct {
	db                       *gorm.DB
	merchantRepo             repository.MerchantRepository
	walletSvc                *WalletService
	withdrawalOrderRepo      repository.WithdrawalOrderRepository
	transactionRecordRepo    repository.TransactionRecordRepository
	cryptoWithdrawalItemRepo repository.FeeTemplateCryptoWithdrawalItemRepository
	fiatWithdrawalItemRepo   repository.FeeTemplateFiatWithdrawalItemRepository
	bankAccountRepo          repository.BankAccountRepository
	cryptoAddressRepo        repository.CryptoAddressRepository
	kunClient                kun.KUNClient
	currencyConfigSvc        *CurrencyConfigService
}

func NewWithdrawalService(
	db *gorm.DB,
	merchantRepo repository.MerchantRepository,
	walletSvc *WalletService,
	withdrawalOrderRepo repository.WithdrawalOrderRepository,
	transactionRecordRepo repository.TransactionRecordRepository,
	cryptoWithdrawalItemRepo repository.FeeTemplateCryptoWithdrawalItemRepository,
	fiatWithdrawalItemRepo repository.FeeTemplateFiatWithdrawalItemRepository,
	bankAccountRepo repository.BankAccountRepository,
	cryptoAddressRepo repository.CryptoAddressRepository,
	kunClient kun.KUNClient,
	currencyConfigSvc *CurrencyConfigService,
) *WithdrawalService {
	return &WithdrawalService{
		db:                       db,
		merchantRepo:             merchantRepo,
		walletSvc:                walletSvc,
		withdrawalOrderRepo:      withdrawalOrderRepo,
		transactionRecordRepo:    transactionRecordRepo,
		cryptoWithdrawalItemRepo: cryptoWithdrawalItemRepo,
		fiatWithdrawalItemRepo:   fiatWithdrawalItemRepo,
		bankAccountRepo:          bankAccountRepo,
		cryptoAddressRepo:        cryptoAddressRepo,
		kunClient:                kunClient,
		currencyConfigSvc:        currencyConfigSvc,
	}
}

func (s *WithdrawalService) SubmitCryptoWithdrawal(ctx context.Context, merchantID uint64, req *dtoreq.SubmitCryptoWithdrawalReq) (uint64, error) {
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

	cryptoAddress, err := s.cryptoAddressRepo.FindByID(ctx, req.CryptoAddressID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, bizerrors.ErrNotFoundError
		}
		return 0, bizerrors.ErrInternalError
	}
	if cryptoAddress.MerchantID != merchantID {
		return 0, bizerrors.ErrForbiddenError
	}
	if cryptoAddress.Status != "ACTIVE" {
		return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "crypto address is not active")
	}
	if cryptoAddress.Currency != req.Currency {
		return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "crypto address currency does not match withdrawal currency")
	}
	if err := s.currencyConfigSvc.EnsureChainSupported(ctx, merchant, cryptoAddress.Currency, cryptoAddress.Chain); err != nil {
		return 0, err
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "invalid amount")
	}

	platformFee := s.calculateCryptoFee(ctx, merchant.FeeTemplateID, cryptoAddress.Currency, cryptoAddress.Chain, amount)
	totalDeduction := amount.Add(platformFee)

	subType := "CRYPTO"
	var orderID uint64

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.walletSvc.FreezeBalance(ctx, tx, merchantID, "FUNDING", req.Currency, totalDeduction); err != nil {
			return err
		}

		platformOrderID := utils.GeneratePlatformOrderID("WD")
		txRecord := &model.TransactionRecord{
			PlatformOrderID: platformOrderID,
			MerchantID:      merchantID,
			Type:            "WITHDRAWAL",
			SubType:         &subType,
			Amount:          amount,
			Currency:        req.Currency,
			PlatformFee:     platformFee,
			Status:          "PENDING",
		}
		if err := tx.WithContext(ctx).Create(txRecord).Error; err != nil {
			return err
		}

		withdrawalOrder := &model.WithdrawalOrder{
			TransactionRecordID: txRecord.ID,
			MerchantID:          merchantID,
			WithdrawalType:      "CRYPTO",
			CryptoAddressID:     &req.CryptoAddressID,
			ToAddress:           &cryptoAddress.Address,
			Chain:               &cryptoAddress.Chain,
			ReviewStatus:        "PENDING_REVIEW",
		}
		if err := tx.WithContext(ctx).Create(withdrawalOrder).Error; err != nil {
			return err
		}

		orderID = withdrawalOrder.ID
		return nil
	})

	if err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			return 0, bizErr
		}
		slog.Error("submit crypto withdrawal failed", slog.Any("error", err))
		return 0, bizerrors.ErrInternalError
	}

	return orderID, nil
}

func (s *WithdrawalService) PreviewWithdrawalFee(ctx context.Context, merchantID uint64, req *dtoreq.WithdrawalFeePreviewReq) (*dtoresp.WithdrawalFeePreviewResp, error) {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	if err := s.currencyConfigSvc.EnsureCurrencySupported(ctx, merchant, req.Currency); err != nil {
		return nil, err
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "invalid amount")
	}

	var platformFee decimal.Decimal

	switch strings.ToUpper(req.Type) {
	case "CRYPTO":
		if req.CryptoAddressID == 0 {
			return nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "crypto_address_id is required")
		}
		cryptoAddress, err := s.cryptoAddressRepo.FindByID(ctx, req.CryptoAddressID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, bizerrors.ErrNotFoundError
			}
			return nil, bizerrors.ErrInternalError
		}
		if cryptoAddress.MerchantID != merchantID {
			return nil, bizerrors.ErrForbiddenError
		}
		if cryptoAddress.Status != "ACTIVE" {
			return nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "crypto address is not active")
		}
		if cryptoAddress.Currency != req.Currency {
			return nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "crypto address currency does not match withdrawal currency")
		}
		platformFee = s.calculateCryptoFee(ctx, merchant.FeeTemplateID, cryptoAddress.Currency, cryptoAddress.Chain, amount)
	case "FIAT":
		if req.BankAccountID == 0 {
			return nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "bank_account_id is required")
		}
		bankAccount, err := s.bankAccountRepo.FindByID(ctx, req.BankAccountID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, bizerrors.ErrNotFoundError
			}
			return nil, bizerrors.ErrInternalError
		}
		if bankAccount.MerchantID != merchantID {
			return nil, bizerrors.ErrForbiddenError
		}
		if bankAccount.Status != "ACTIVE" {
			return nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "bank account is not active")
		}
		if bankAccount.CurrencyList != req.Currency {
			return nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "bank account currency does not match withdrawal currency")
		}
		platformFee = s.calculateFiatFee(ctx, merchant.FeeTemplateID, req.Currency, bankAccount.TransferType, amount)
	default:
		return nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "invalid withdrawal type")
	}

	totalDeduction := amount.Add(platformFee)

	return &dtoresp.WithdrawalFeePreviewResp{
		Currency:       req.Currency,
		Amount:         amount.String(),
		PlatformFee:    platformFee.String(),
		TotalDeduction: totalDeduction.String(),
		NetAmount:      amount.String(),
	}, nil
}

func (s *WithdrawalService) SubmitFiatWithdrawal(ctx context.Context, merchantID uint64, req *dtoreq.SubmitFiatWithdrawalReq) (uint64, error) {
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

	bankAccount, err := s.bankAccountRepo.FindByID(ctx, req.BankAccountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, bizerrors.ErrNotFoundError
		}
		return 0, bizerrors.ErrInternalError
	}
	if bankAccount.MerchantID != merchantID {
		return 0, bizerrors.ErrForbiddenError
	}
	if bankAccount.Status != "ACTIVE" {
		return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "bank account is not active")
	}
	if bankAccount.CurrencyList != req.Currency {
		return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "bank account currency does not match withdrawal currency")
	}

	if err := validateFiatWithdrawalPurpose(req.Purpose); err != nil {
		return 0, err
	}
	postscript := strings.TrimSpace(req.Postscript)
	if postscript == "" {
		return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "postscript is required")
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "invalid amount")
	}

	platformFee := s.calculateFiatFee(ctx, merchant.FeeTemplateID, req.Currency, bankAccount.TransferType, amount)
	totalDeduction := amount.Add(platformFee)
	purpose := strings.ToUpper(strings.TrimSpace(req.Purpose))

	subType := "FIAT"
	var orderID uint64

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.walletSvc.FreezeBalance(ctx, tx, merchantID, "FUNDING", req.Currency, totalDeduction); err != nil {
			return err
		}

		platformOrderID := utils.GeneratePlatformOrderID("WD")
		txRecord := &model.TransactionRecord{
			PlatformOrderID: platformOrderID,
			MerchantID:      merchantID,
			Type:            "WITHDRAWAL",
			SubType:         &subType,
			Amount:          amount,
			Currency:        req.Currency,
			PlatformFee:     platformFee,
			Remark:          &postscript,
			Status:          "PENDING",
		}
		if err := tx.WithContext(ctx).Create(txRecord).Error; err != nil {
			return err
		}

		bankAccountID := req.BankAccountID
		withdrawalOrder := &model.WithdrawalOrder{
			TransactionRecordID: txRecord.ID,
			MerchantID:          merchantID,
			WithdrawalType:      "FIAT",
			BankAccountID:       &bankAccountID,
			TransferType:        &bankAccount.TransferType,
			Purpose:             &purpose,
			ReviewStatus:        "PENDING_REVIEW",
		}
		if err := tx.WithContext(ctx).Create(withdrawalOrder).Error; err != nil {
			return err
		}

		orderID = withdrawalOrder.ID
		return nil
	})

	if err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			return 0, bizErr
		}
		slog.Error("submit fiat withdrawal failed", slog.Any("error", err))
		return 0, bizerrors.ErrInternalError
	}

	return orderID, nil
}

func (s *WithdrawalService) ApproveWithdrawal(ctx context.Context, adminID uint64, orderID uint64) error {
	order, err := s.withdrawalOrderRepo.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrWithdrawalNotFoundE
		}
		return bizerrors.ErrInternalError
	}

	if order.ReviewStatus != "PENDING_REVIEW" {
		return bizerrors.ErrWithdrawalNotPendingE
	}

	merchant, err := s.merchantRepo.FindByID(ctx, order.MerchantID)
	if err != nil {
		return bizerrors.ErrInternalError
	}
	if merchant.KunSubCustomerNo == nil {
		return bizerrors.ErrMerchantNotRegisteredE
	}

	txRecord, err := s.transactionRecordRepo.FindByID(ctx, order.TransactionRecordID)
	if err != nil {
		return bizerrors.ErrInternalError
	}

	requestNo := kun.GenerateRequestNo()
	var kunOrderID string

	switch order.WithdrawalType {
	case "CRYPTO":
		if order.ToAddress == nil || order.Chain == nil {
			return bizerrors.ErrInternalError
		}
		var kunOrderResp kundto.CryptoWithdrawalResp
		err = s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, kun.CryptoWithdrawalPath, &kundto.CryptoWithdrawalReq{
			RequestNo:  requestNo,
			Amount:     txRecord.Amount.String(),
			Chain:      *order.Chain,
			Currency:   txRecord.Currency,
			Address:    *order.ToAddress,
			RegionCode: s.kunClient.GetRegionCode(),
		}, &kunOrderResp)
		if err != nil {
			slog.Error("KUN crypto withdrawal failed", slog.Any("error", err))
			return bizerrors.ErrKUNAPIFailedE
		}
		kunOrderID = string(kunOrderResp)
	case "FIAT":
		bankAccount, err := s.bankAccountRepo.FindByID(ctx, *order.BankAccountID)
		if err != nil {
			return bizerrors.ErrInternalError
		}
		if bankAccount.KunAccountID == nil {
			return bizerrors.ErrKUNAPIFailedE
		}
		purpose := "OTHER"
		if order.Purpose != nil && strings.TrimSpace(*order.Purpose) != "" {
			purpose = strings.ToUpper(strings.TrimSpace(*order.Purpose))
		}
		postscript := ""
		if txRecord.Remark != nil {
			postscript = strings.TrimSpace(*txRecord.Remark)
		}
		if postscript == "" {
			postscript = "Withdrawal"
		}
		var resp kundto.FiatWithdrawalResp
		err = s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, kun.FiatWithdrawalPath, &kundto.FiatWithdrawalReq{
			RequestNo:  requestNo,
			AccountId:  *bankAccount.KunAccountID,
			Amount:     txRecord.Amount.String(),
			Currency:   txRecord.Currency,
			FeeMethod:  "SHA",
			PoboType:   "NO",
			Postscript: postscript,
			Purpose:    purpose,
		}, &resp)
		if err != nil {
			slog.Error("KUN fiat withdrawal failed", slog.Any("error", err))
			return bizerrors.ErrKUNAPIFailedE
		}
		kunOrderID = resp.OrderId
	}

	now := time.Now()
	reviewerType := "ADMIN"
	return s.withdrawalOrderRepo.UpdateFields(ctx, orderID, map[string]interface{}{
		"review_status":    "APPROVED",
		"reviewer_id":      adminID,
		"reviewer_type":    reviewerType,
		"reviewed_at":      now,
		"kun_order_id":     kunOrderID,
		"kun_request_no":   requestNo,
		"kun_submitted_at": now,
	})
}

func (s *WithdrawalService) RejectWithdrawal(ctx context.Context, adminID uint64, orderID uint64, reason string) error {
	order, err := s.withdrawalOrderRepo.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrWithdrawalNotFoundE
		}
		return bizerrors.ErrInternalError
	}

	if order.ReviewStatus != "PENDING_REVIEW" {
		return bizerrors.ErrWithdrawalNotPendingE
	}

	txRecord, err := s.transactionRecordRepo.FindByID(ctx, order.TransactionRecordID)
	if err != nil {
		return bizerrors.ErrInternalError
	}

	totalFrozen := txRecord.Amount.Add(txRecord.PlatformFee)

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.walletSvc.UnfreezeBalance(ctx, tx, order.MerchantID, "FUNDING", txRecord.Currency, totalFrozen); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		slog.Error("unfreeze balance failed during rejection", slog.Any("error", err))
		return bizerrors.ErrInternalError
	}

	now := time.Now()
	reviewerType := "ADMIN"
	if err := s.withdrawalOrderRepo.UpdateFields(ctx, orderID, map[string]interface{}{
		"review_status": "REJECTED",
		"reviewer_id":   adminID,
		"reviewer_type": reviewerType,
		"reviewed_at":   now,
		"review_remark": reason,
	}); err != nil {
		return bizerrors.ErrInternalError
	}

	return s.transactionRecordRepo.UpdateStatus(ctx, order.TransactionRecordID, "FAILED")
}

func (s *WithdrawalService) ListWithdrawals(ctx context.Context, merchantID uint64, page, pageSize int) (*dtoresp.WithdrawalOrderListResp, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	orders, total, err := s.withdrawalOrderRepo.ListByMerchant(ctx, merchantID, page, pageSize)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	var resp []dtoresp.WithdrawalOrderResp
	for _, o := range orders {
		txRecord, _ := s.transactionRecordRepo.FindByID(ctx, o.TransactionRecordID)
		item := dtoresp.WithdrawalOrderResp{
			ID:           o.ID,
			Type:         o.WithdrawalType,
			Chain:        o.Chain,
			ToAddress:    o.ToAddress,
			TxID:         o.TxID,
			ReviewStatus: o.ReviewStatus,
			CreatedAt:    o.CreatedAt,
		}
		if txRecord != nil {
			item.Currency = txRecord.Currency
			item.Amount = txRecord.Amount.String()
			item.PlatformFee = txRecord.PlatformFee.String()
			item.Status = txRecord.Status
		}
		resp = append(resp, item)
	}

	return &dtoresp.WithdrawalOrderListResp{
		Orders: resp,
		Total:  total,
	}, nil
}

func (s *WithdrawalService) GetWithdrawalDetail(ctx context.Context, merchantID uint64, orderID uint64) (*dtoresp.WithdrawalOrderResp, error) {
	order, err := s.withdrawalOrderRepo.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrWithdrawalNotFoundE
		}
		return nil, bizerrors.ErrInternalError
	}

	if order.MerchantID != merchantID {
		return nil, bizerrors.ErrForbiddenError
	}

	txRecord, err := s.transactionRecordRepo.FindByID(ctx, order.TransactionRecordID)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	return &dtoresp.WithdrawalOrderResp{
		ID:           order.ID,
		Type:         order.WithdrawalType,
		Currency:     txRecord.Currency,
		Chain:        order.Chain,
		Amount:       txRecord.Amount.String(),
		PlatformFee:  txRecord.PlatformFee.String(),
		Status:       txRecord.Status,
		ReviewStatus: order.ReviewStatus,
		ToAddress:    order.ToAddress,
		TxID:         order.TxID,
		CreatedAt:    order.CreatedAt,
	}, nil
}

func (s *WithdrawalService) ListPendingReviews(ctx context.Context, page, pageSize int) (*dtoresp.WithdrawalOrderListResp, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	orders, total, err := s.withdrawalOrderRepo.ListPendingReview(ctx, page, pageSize)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	var resp []dtoresp.WithdrawalOrderResp
	for _, o := range orders {
		txRecord, _ := s.transactionRecordRepo.FindByID(ctx, o.TransactionRecordID)
		item := dtoresp.WithdrawalOrderResp{
			ID:           o.ID,
			Type:         o.WithdrawalType,
			Chain:        o.Chain,
			ToAddress:    o.ToAddress,
			ReviewStatus: o.ReviewStatus,
			CreatedAt:    o.CreatedAt,
		}
		if txRecord != nil {
			item.Currency = txRecord.Currency
			item.Amount = txRecord.Amount.String()
			item.PlatformFee = txRecord.PlatformFee.String()
			item.Status = txRecord.Status
		}
		resp = append(resp, item)
	}

	return &dtoresp.WithdrawalOrderListResp{
		Orders: resp,
		Total:  total,
	}, nil
}

func (s *WithdrawalService) calculateCryptoFee(ctx context.Context, feeTemplateID *uint64, currency, chain string, amount decimal.Decimal) decimal.Decimal {
	if feeTemplateID == nil {
		return decimal.Zero
	}

	items, err := s.cryptoWithdrawalItemRepo.FindByTemplateID(ctx, *feeTemplateID)
	if err != nil {
		return decimal.Zero
	}

	for _, item := range items {
		if item.Currency == currency && item.Chain == chain {
			rateFee := amount.Mul(item.FeeRate)
			if item.FixedFee.GreaterThan(rateFee) {
				return item.FixedFee
			}
			return rateFee
		}
	}

	return decimal.Zero
}

func (s *WithdrawalService) calculateFiatFee(ctx context.Context, feeTemplateID *uint64, currency, transferType string, amount decimal.Decimal) decimal.Decimal {
	if feeTemplateID == nil {
		return decimal.Zero
	}

	items, err := s.fiatWithdrawalItemRepo.FindByTemplateID(ctx, *feeTemplateID)
	if err != nil {
		return decimal.Zero
	}

	for _, item := range items {
		if item.Currency == currency && item.TransferType == transferType {
			rateFee := amount.Mul(item.FeeRate)
			if item.FixedFee.GreaterThan(rateFee) {
				return item.FixedFee
			}
			return rateFee
		}
	}

	return decimal.Zero
}
