package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	dtoreq "motewallet-withdrawal/backend/internal/dto/request"
	dtoresp "motewallet-withdrawal/backend/internal/dto/response"
	"motewallet-withdrawal/backend/internal/model"
	bizerrors "motewallet-withdrawal/backend/internal/pkg/errors"
	"motewallet-withdrawal/backend/internal/pkg/kun"
	kundto "motewallet-withdrawal/backend/internal/pkg/kun/dto"
	"motewallet-withdrawal/backend/internal/pkg/utils"
	"motewallet-withdrawal/backend/internal/repository"
)

type ExchangeService struct {
	db                    *gorm.DB
	merchantRepo          repository.MerchantRepository
	walletSvc             *WalletService
	exchangeOrderRepo     repository.ExchangeOrderRepository
	transactionRecordRepo repository.TransactionRecordRepository
	exchangeItemRepo      repository.FeeTemplateExchangeItemRepository
	kunClient             kun.KUNClient
}

func NewExchangeService(
	db *gorm.DB,
	merchantRepo repository.MerchantRepository,
	walletSvc *WalletService,
	exchangeOrderRepo repository.ExchangeOrderRepository,
	transactionRecordRepo repository.TransactionRecordRepository,
	exchangeItemRepo repository.FeeTemplateExchangeItemRepository,
	kunClient kun.KUNClient,
) *ExchangeService {
	return &ExchangeService{
		db:                    db,
		merchantRepo:          merchantRepo,
		walletSvc:             walletSvc,
		exchangeOrderRepo:     exchangeOrderRepo,
		transactionRecordRepo: transactionRecordRepo,
		exchangeItemRepo:      exchangeItemRepo,
		kunClient:             kunClient,
	}
}

func (s *ExchangeService) GetQuote(ctx context.Context, merchantID uint64, req *dtoreq.GetExchangeQuoteReq) (*dtoresp.ExchangeQuoteResp, error) {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	if merchant.KunSubCustomerNo == nil {
		return nil, bizerrors.ErrMerchantNotRegisteredE
	}

	var quoteResp kundto.ExchangeQuoteResp
	err = s.kunClient.Post(ctx, "/rest/v2.0/trade/exchange/marketInquiry", &kundto.ExchangeQuoteReq{
		SubCustomerNo: *merchant.KunSubCustomerNo,
		FromCurrency:  req.FromCurrency,
		ToCurrency:    req.ToCurrency,
		FromAmount:    req.FromAmount,
	}, &quoteResp)
	if err != nil {
		slog.Error("KUN exchange quote failed", slog.Any("error", err))
		return nil, bizerrors.ErrKUNAPIFailedE
	}

	fromAmount, _ := decimal.NewFromString(req.FromAmount)
	platformFee := s.calculateExchangeFee(ctx, merchant.FeeTemplateID, req.FromCurrency, req.ToCurrency, fromAmount)

	return &dtoresp.ExchangeQuoteResp{
		QuoteId:      quoteResp.QuoteId,
		FromCurrency: quoteResp.FromCurrency,
		ToCurrency:   quoteResp.ToCurrency,
		FromAmount:   quoteResp.FromAmount,
		ToAmount:     quoteResp.ToAmount,
		ExchangeRate: quoteResp.ExchangeRate,
		PlatformFee:  platformFee.String(),
		FeeCurrency:  req.FromCurrency,
		ExpireTime:   quoteResp.ExpireTime,
	}, nil
}

func (s *ExchangeService) CreateExchangeOrder(ctx context.Context, merchantID uint64, req *dtoreq.CreateExchangeOrderReq) (uint64, error) {
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

	fromAmount, err := decimal.NewFromString(req.FromAmount)
	if err != nil || fromAmount.LessThanOrEqual(decimal.Zero) {
		return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "invalid amount")
	}

	platformFee := s.calculateExchangeFee(ctx, merchant.FeeTemplateID, req.FromCurrency, req.ToCurrency, fromAmount)

	requestNo := kun.GenerateRequestNo()
	var kunResp kundto.ExchangeOrderResp
	err = s.kunClient.Post(ctx, "/rest/v2.0/trade/exchange/order", &kundto.ExchangeOrderReq{
		SubCustomerNo: *merchant.KunSubCustomerNo,
		RequestNo:     requestNo,
		QuoteId:       req.QuoteId,
		FromCurrency:  req.FromCurrency,
		ToCurrency:    req.ToCurrency,
		FromAmount:    req.FromAmount,
	}, &kunResp)
	if err != nil {
		slog.Error("KUN exchange order failed", slog.Any("error", err))
		return 0, bizerrors.ErrKUNAPIFailedE
	}

	subType := "MARKET"
	var orderID uint64

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		platformOrderID := utils.GeneratePlatformOrderID("EX")
		txRecord := &model.TransactionRecord{
			PlatformOrderID: platformOrderID,
			MerchantID:      merchantID,
			Type:            "EXCHANGE",
			SubType:         &subType,
			Amount:          fromAmount,
			Currency:        req.FromCurrency,
			PlatformFee:     platformFee,
			Status:          "PROCESSING",
		}
		if err := tx.WithContext(ctx).Create(txRecord).Error; err != nil {
			return err
		}

		exchangeOrder := &model.ExchangeOrder{
			TransactionRecordID: txRecord.ID,
			MerchantID:          merchantID,
			ExchangeType:        "MARKET",
			FromCurrency:        req.FromCurrency,
			FromAmount:          fromAmount,
			ToCurrency:          req.ToCurrency,
			QuoteID:             &req.QuoteId,
			KunOrderID:          &kunResp.OrderId,
			KunRequestNo:        &requestNo,
		}
		if err := tx.WithContext(ctx).Create(exchangeOrder).Error; err != nil {
			return err
		}

		orderID = exchangeOrder.ID
		return nil
	})

	if err != nil {
		slog.Error("create exchange order failed", slog.Any("error", err))
		return 0, bizerrors.ErrInternalError
	}

	return orderID, nil
}

func (s *ExchangeService) Create1to1Order(ctx context.Context, merchantID uint64, req *dtoreq.Create1to1OrderReq) (uint64, error) {
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

	fromAmount, err := decimal.NewFromString(req.FromAmount)
	if err != nil || fromAmount.LessThanOrEqual(decimal.Zero) {
		return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "invalid amount")
	}

	platformFee := s.calculateExchangeFee(ctx, merchant.FeeTemplateID, req.FromCurrency, req.ToCurrency, fromAmount)

	requestNo := kun.GenerateRequestNo()
	var kunResp kundto.InnerMatchCreateResp
	err = s.kunClient.Post(ctx, "/rest/v2.0/trade/inner/match/create", &kundto.InnerMatchCreateReq{
		SubCustomerNo: *merchant.KunSubCustomerNo,
		RequestNo:     requestNo,
		FromCurrency:  req.FromCurrency,
		ToCurrency:    req.ToCurrency,
		FromAmount:    req.FromAmount,
	}, &kunResp)
	if err != nil {
		slog.Error("KUN 1:1 order failed", slog.Any("error", err))
		return 0, bizerrors.ErrKUNAPIFailedE
	}

	subType := "1TO1"
	var orderID uint64

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		platformOrderID := utils.GeneratePlatformOrderID("EX")
		txRecord := &model.TransactionRecord{
			PlatformOrderID: platformOrderID,
			MerchantID:      merchantID,
			Type:            "EXCHANGE",
			SubType:         &subType,
			Amount:          fromAmount,
			Currency:        req.FromCurrency,
			PlatformFee:     platformFee,
			Status:          "PROCESSING",
		}
		if err := tx.WithContext(ctx).Create(txRecord).Error; err != nil {
			return err
		}

		exchangeOrder := &model.ExchangeOrder{
			TransactionRecordID: txRecord.ID,
			MerchantID:          merchantID,
			ExchangeType:        "1TO1",
			FromCurrency:        req.FromCurrency,
			FromAmount:          fromAmount,
			ToCurrency:          req.ToCurrency,
			KunOrderID:          &kunResp.OrderId,
			KunRequestNo:        &requestNo,
		}
		if err := tx.WithContext(ctx).Create(exchangeOrder).Error; err != nil {
			return err
		}

		orderID = exchangeOrder.ID
		return nil
	})

	if err != nil {
		slog.Error("create 1to1 order failed", slog.Any("error", err))
		return 0, bizerrors.ErrInternalError
	}

	return orderID, nil
}

func (s *ExchangeService) ListExchangeOrders(ctx context.Context, merchantID uint64, page, pageSize int) (*dtoresp.ExchangeOrderListResp, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	orders, total, err := s.exchangeOrderRepo.ListByMerchant(ctx, merchantID, page, pageSize)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	var resp []dtoresp.ExchangeOrderResp
	for _, o := range orders {
		txRecord, _ := s.transactionRecordRepo.FindByID(ctx, o.TransactionRecordID)
		item := dtoresp.ExchangeOrderResp{
			ID:           o.ID,
			ExchangeType: o.ExchangeType,
			FromCurrency: o.FromCurrency,
			ToCurrency:   o.ToCurrency,
			FromAmount:   o.FromAmount.String(),
			CreatedAt:    o.CreatedAt,
		}
		if o.ToAmount != nil {
			item.ToAmount = o.ToAmount.String()
		}
		if o.ExchangeRate != nil {
			item.ExchangeRate = o.ExchangeRate.String()
		}
		if txRecord != nil {
			item.PlatformFee = txRecord.PlatformFee.String()
			item.Status = txRecord.Status
		}
		resp = append(resp, item)
	}

	return &dtoresp.ExchangeOrderListResp{
		Orders: resp,
		Total:  total,
	}, nil
}

func (s *ExchangeService) calculateExchangeFee(ctx context.Context, feeTemplateID *uint64, fromCurrency, toCurrency string, amount decimal.Decimal) decimal.Decimal {
	if feeTemplateID == nil {
		return decimal.Zero
	}

	items, err := s.exchangeItemRepo.FindByTemplateID(ctx, *feeTemplateID)
	if err != nil {
		return decimal.Zero
	}

	for _, item := range items {
		if item.FromCurrency == fromCurrency && item.ToCurrency == toCurrency {
			rateFee := amount.Mul(item.FeeRate)
			if item.MinFee.GreaterThan(rateFee) {
				return item.MinFee
			}
			return rateFee
		}
	}

	return decimal.Zero
}
