package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	dtoreq "motewallet/internal/dto/request"
	dtoresp "motewallet/internal/dto/response"
	"motewallet/internal/model"
	bizerrors "motewallet/internal/pkg/errors"
	exchangepkg "motewallet/internal/pkg/exchange"
	"motewallet/internal/pkg/kun"
	kundto "motewallet/internal/pkg/kun/dto"
	"motewallet/internal/pkg/utils"
	"motewallet/internal/repository"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ExchangeService struct {
	db                    *gorm.DB
	merchantRepo          repository.MerchantRepository
	walletSvc             *WalletService
	exchangeOrderRepo     repository.ExchangeOrderRepository
	transactionRecordRepo repository.TransactionRecordRepository
	exchangeItemRepo      repository.FeeTemplateExchangeItemRepository
	kunClient             kun.KUNClient
	currencyConfigSvc     *CurrencyConfigService
}

func NewExchangeService(
	db *gorm.DB,
	merchantRepo repository.MerchantRepository,
	walletSvc *WalletService,
	exchangeOrderRepo repository.ExchangeOrderRepository,
	transactionRecordRepo repository.TransactionRecordRepository,
	exchangeItemRepo repository.FeeTemplateExchangeItemRepository,
	kunClient kun.KUNClient,
	currencyConfigSvc *CurrencyConfigService,
) *ExchangeService {
	return &ExchangeService{
		db:                    db,
		merchantRepo:          merchantRepo,
		walletSvc:             walletSvc,
		exchangeOrderRepo:     exchangeOrderRepo,
		transactionRecordRepo: transactionRecordRepo,
		exchangeItemRepo:      exchangeItemRepo,
		kunClient:             kunClient,
		currencyConfigSvc:     currencyConfigSvc,
	}
}

func (s *ExchangeService) PreviewExchange(ctx context.Context, merchantID uint64, req *dtoreq.ExchangePreviewReq) (*dtoresp.ExchangePreviewResp, error) {
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

	if err := s.validateExchangeRequest(ctx, merchant, req.FromCurrency, req.ToCurrency, req.FromAmount); err != nil {
		return nil, err
	}

	fromAmount, _ := decimal.NewFromString(req.FromAmount)
	platformFee := s.calculateExchangeFee(ctx, merchant.FeeTemplateID, req.FromCurrency, req.ToCurrency, fromAmount)

	// 暂时使用 1:1 兑换，询价接口保留待恢复
	// quote, err := s.requestKUNQuote(ctx, *merchant.KunSubCustomerNo, req.FromCurrency, req.ToCurrency, req.FromAmount)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// toAmount, err := decimal.NewFromString(quote.ToAmount)
	// if err != nil || toAmount.LessThanOrEqual(decimal.Zero) {
	// 	slog.Error("invalid KUN quote to amount", slog.String("to_amount", quote.ToAmount))
	// 	return nil, bizerrors.ErrKUNAPIFailedE
	// }
	//
	// exchangeRate := quote.ExchangeRate
	// if strings.TrimSpace(exchangeRate) == "" {
	// 	exchangeRate = toAmount.Div(fromAmount).String()
	// }
	//
	// return &dtoresp.ExchangePreviewResp{
	// 	FromCurrency:   req.FromCurrency,
	// 	ToCurrency:     req.ToCurrency,
	// 	FromAmount:     fromAmount.String(),
	// 	ToAmount:       toAmount.String(),
	// 	ExchangeRate:   exchangeRate,
	// 	QuoteID:        quote.QuoteId,
	// 	ExpireTime:     quote.ExpireTime,
	// 	PlatformFee:    platformFee.String(),
	// 	FeeCurrency:    req.FromCurrency,
	// 	NetToAmount:    toAmount.String(),
	// 	TotalDeduction: fromAmount.Add(platformFee).String(),
	// }, nil

	return &dtoresp.ExchangePreviewResp{
		FromCurrency:   req.FromCurrency,
		ToCurrency:     req.ToCurrency,
		FromAmount:     fromAmount.String(),
		ToAmount:       fromAmount.String(),
		ExchangeRate:   "1",
		PlatformFee:    platformFee.String(),
		FeeCurrency:    req.FromCurrency,
		NetToAmount:    fromAmount.String(),
		TotalDeduction: fromAmount.Add(platformFee).String(),
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

	if err := s.validateExchangeRequest(ctx, merchant, req.FromCurrency, req.ToCurrency, req.FromAmount); err != nil {
		return 0, err
	}

	fromAmount, _ := decimal.NewFromString(req.FromAmount)
	platformFee := s.calculateExchangeFee(ctx, merchant.FeeTemplateID, req.FromCurrency, req.ToCurrency, fromAmount)
	totalFreeze := fromAmount.Add(platformFee)

	requestNo := kun.GenerateRequestNo()
	autoTransfer := "YES"
	subType := "1TO1"
	exchangeRate := decimal.NewFromInt(1)
	var orderID uint64
	var txRecordID uint64

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
			ToAmount:            &fromAmount,
			ExchangeRate:        &exchangeRate,
			AutoTransfer:        &autoTransfer,
			KunRequestNo:        &requestNo,
		}
		if err := tx.WithContext(ctx).Create(exchangeOrder).Error; err != nil {
			return err
		}

		ref := WalletChangeRef{TransactionRecordID: txRecord.ID, BizType: "EXCHANGE"}
		if err := s.walletSvc.FreezeBalance(ctx, tx, merchantID, "FUNDING", req.FromCurrency, totalFreeze, ref); err != nil {
			return err
		}

		orderID = exchangeOrder.ID
		txRecordID = txRecord.ID
		return nil
	})
	if err != nil {
		slog.Error("create exchange order failed", slog.Any("error", err))
		return 0, err
	}

	var kunResp kundto.InnerMatchCreateResp
	err = s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, kun.InnerMatchCreatePath, &kundto.InnerMatchCreateReq{
		RequestNo:    requestNo,
		FromCurrency: req.FromCurrency,
		OrderAmount:  req.FromAmount,
		ToCurrency:   req.ToCurrency,
		AutoTransfer: autoTransfer,
	}, &kunResp)
	if err != nil {
		slog.Error("KUN 1:1 exchange order failed", slog.Any("error", err))
		s.rollbackExchangeOrder(ctx, txRecordID, merchantID, req.FromCurrency, totalFreeze)
		return 0, bizerrors.ErrKUNAPIFailedE
	}

	// 现货询价下单暂时停用
	// if strings.TrimSpace(req.QuoteID) == "" {
	// 	return 0, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "quote_id is required")
	// }
	// quoteID := strings.TrimSpace(req.QuoteID)
	// var kunResp kundto.ExchangeOrderResp
	// err = s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, kun.ExchangeOrderPath, &kundto.ExchangeOrderReq{
	// 	SubCustomerNo: *merchant.KunSubCustomerNo,
	// 	RequestNo:     requestNo,
	// 	QuoteId:       quoteID,
	// 	FromCurrency:  req.FromCurrency,
	// 	ToCurrency:    req.ToCurrency,
	// 	FromAmount:    req.FromAmount,
	// }, &kunResp)
	// if err != nil {
	// 	slog.Error("KUN spot exchange order failed", slog.Any("error", err))
	// 	s.rollbackExchangeOrder(ctx, txRecordID, merchantID, req.FromCurrency, totalFreeze)
	// 	return 0, bizerrors.ErrKUNAPIFailedE
	// }

	if err := s.exchangeOrderRepo.UpdateFields(ctx, orderID, map[string]interface{}{
		"kun_order_id": kunResp.OrderId,
	}); err != nil {
		slog.Error("update exchange kun order id failed", slog.Any("error", err))
	}

	return orderID, nil
}

func (s *ExchangeService) QueryKUNExchangeOrder(ctx context.Context, merchantID, orderID uint64) (*kundto.ExchangeOrderQueryResp, error) {
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

	order, err := s.exchangeOrderRepo.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}
	if order.MerchantID != merchantID {
		return nil, bizerrors.ErrForbiddenError
	}
	if order.KunOrderID == nil || *order.KunOrderID == "" {
		return nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "exchange order has no kun order id")
	}

	kunResp, err := s.queryKUNExchangeOrder(ctx, merchant, order)
	if err != nil {
		return nil, err
	}

	return kunResp, nil
}

func (s *ExchangeService) SettleFromWebhook(ctx context.Context, data *kundto.ExchangeData) error {
	order, err := s.exchangeOrderRepo.FindByKunRequestNo(ctx, data.RequestNo)
	if err != nil {
		return err
	}

	txRecord, err := s.transactionRecordRepo.FindByID(ctx, order.TransactionRecordID)
	if err != nil {
		return err
	}

	return s.applyKUNExchangeStatus(ctx, order, txRecord, data.OrderStatus, data.ToAmount, data.ExchangeRate, data.TradeFee, resolveExchangeFailReason(data.RejectReason, data.FailReason))
}

func (s *ExchangeService) SyncOrderStatus(ctx context.Context, orderID uint64) (*dtoresp.AdminExchangeSyncResp, error) {
	order, err := s.exchangeOrderRepo.FindByID(ctx, orderID)
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

	resp := &dtoresp.AdminExchangeSyncResp{
		OrderID:     order.ID,
		Status:      txRecord.Status,
		PlatformFee: txRecord.PlatformFee.String(),
		Updated:     false,
	}
	if order.ToAmount != nil {
		resp.ToAmount = order.ToAmount.String()
	}
	if order.ExchangeRate != nil {
		resp.ExchangeRate = order.ExchangeRate.String()
	}

	if txRecord.Status == "COMPLETED" {
		if order.FailReason != nil {
			resp.FailReason = *order.FailReason
		}
		return resp, nil
	}

	if order.KunOrderID == nil || strings.TrimSpace(*order.KunOrderID) == "" {
		return nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "exchange order has no kun order id")
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

	var kunResp kundto.ExchangeOrderQueryResp
	err = s.queryKUNExchangeOrderInto(ctx, merchant, order, &kunResp)
	if err != nil {
		return nil, err
	}

	resp.KunStatus = kunResp.OrderStatus
	toAmount := kunResp.ToAmount
	exchangeRate := kunResp.ExchangeRate
	if strings.TrimSpace(exchangeRate) == "" {
		exchangeRate = "1"
	}

	beforeStatus := txRecord.Status
	beforeFailReason := ""
	if order.FailReason != nil {
		beforeFailReason = *order.FailReason
	}
	if err := s.applyKUNExchangeStatus(ctx, order, txRecord, kunResp.OrderStatus, toAmount, exchangeRate, kunResp.TradeFee, ""); err != nil {
		return nil, err
	}

	updatedTxRecord, err := s.transactionRecordRepo.FindByID(ctx, order.TransactionRecordID)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}
	updatedOrder, err := s.exchangeOrderRepo.FindByID(ctx, order.ID)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	resp.Status = updatedTxRecord.Status
	afterFailReason := ""
	if updatedOrder.FailReason != nil {
		afterFailReason = *updatedOrder.FailReason
	}
	resp.Updated = beforeStatus != updatedTxRecord.Status || beforeFailReason != afterFailReason
	resp.PlatformFee = updatedTxRecord.PlatformFee.String()
	if updatedOrder.ToAmount != nil {
		resp.ToAmount = updatedOrder.ToAmount.String()
	}
	if updatedOrder.ExchangeRate != nil {
		resp.ExchangeRate = updatedOrder.ExchangeRate.String()
	}
	if updatedOrder.FailReason != nil {
		resp.FailReason = *updatedOrder.FailReason
	}

	return resp, nil
}

func (s *ExchangeService) applyKUNExchangeStatus(
	ctx context.Context,
	order *model.ExchangeOrder,
	txRecord *model.TransactionRecord,
	kunStatus, toAmountStr, exchangeRateStr, tradeFeeStr, failReason string,
) error {
	if txRecord.Status == "COMPLETED" {
		return nil
	}
	if txRecord.Status == "FAILED" {
		return s.refreshFailedExchangeFromKUN(ctx, order, kunStatus, failReason)
	}

	switch strings.ToUpper(kunStatus) {
	case "SUCCESS":
		toAmount, err := decimal.NewFromString(toAmountStr)
		if err != nil || toAmount.LessThanOrEqual(decimal.Zero) {
			toAmount = order.FromAmount
		}
		exchangeRate, err := decimal.NewFromString(exchangeRateStr)
		if err != nil || exchangeRate.LessThanOrEqual(decimal.Zero) {
			exchangeRate = decimal.NewFromInt(1)
		}
		kunFee, _ := decimal.NewFromString(tradeFeeStr)
		totalFrozen := order.FromAmount.Add(txRecord.PlatformFee)
		now := time.Now()

		ref := WalletChangeRef{TransactionRecordID: txRecord.ID, BizType: "EXCHANGE"}
		return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := s.walletSvc.DeductFrozen(ctx, tx, order.MerchantID, "FUNDING", order.FromCurrency, totalFrozen, ref); err != nil {
				return err
			}
			if err := s.walletSvc.CreditBalance(ctx, tx, order.MerchantID, "FUNDING", order.ToCurrency, toAmount, ref); err != nil {
				return err
			}
			if err := tx.WithContext(ctx).Model(&model.TransactionRecord{}).
				Where("id = ?", txRecord.ID).
				Update("status", "COMPLETED").Error; err != nil {
				return err
			}
			return tx.WithContext(ctx).Model(&model.ExchangeOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
				"to_amount":     toAmount,
				"exchange_rate": exchangeRate,
				"kun_fee":       kunFee,
				"completed_at":  now,
			}).Error
		})
	case "FAIL":
		totalFrozen := order.FromAmount.Add(txRecord.PlatformFee)
		reason := strings.TrimSpace(failReason)
		ref := WalletChangeRef{TransactionRecordID: txRecord.ID, BizType: "EXCHANGE"}
		return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := s.walletSvc.UnfreezeBalance(ctx, tx, order.MerchantID, "FUNDING", order.FromCurrency, totalFrozen, ref); err != nil {
				return err
			}
			if err := tx.WithContext(ctx).Model(&model.TransactionRecord{}).
				Where("id = ?", txRecord.ID).
				Update("status", "FAILED").Error; err != nil {
				return err
			}
			updates := map[string]interface{}{}
			if reason != "" {
				updates["fail_reason"] = reason
			}
			if len(updates) == 0 {
				return nil
			}
			return tx.WithContext(ctx).Model(&model.ExchangeOrder{}).Where("id = ?", order.ID).Updates(updates).Error
		})
	case "PROCESS", "PROCESSING":
		slog.Info("exchange still processing", slog.Uint64("order_id", order.ID), slog.String("kun_status", kunStatus))
		return nil
	default:
		slog.Info("exchange status update", slog.String("status", kunStatus), slog.Uint64("order_id", order.ID))
		return nil
	}
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
		if o.FailReason != nil {
			item.FailReason = *o.FailReason
		}
		resp = append(resp, item)
	}

	return &dtoresp.ExchangeOrderListResp{
		Orders: resp,
		Total:  total,
	}, nil
}

func (s *ExchangeService) validateExchangeRequest(ctx context.Context, merchant *model.Merchant, fromCurrency, toCurrency, fromAmount string) error {
	if err := s.currencyConfigSvc.EnsureCurrencySupported(ctx, merchant, fromCurrency); err != nil {
		return err
	}
	if err := s.currencyConfigSvc.EnsureCurrencySupported(ctx, merchant, toCurrency); err != nil {
		return err
	}
	if err := exchangepkg.EnsureSupportedPair(fromCurrency, toCurrency); err != nil {
		return err
	}

	amount, err := decimal.NewFromString(fromAmount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "invalid amount")
	}

	return nil
}

func (s *ExchangeService) rollbackExchangeOrder(ctx context.Context, txRecordID, merchantID uint64, currency string, totalFreeze decimal.Decimal) {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ref := WalletChangeRef{TransactionRecordID: txRecordID, BizType: "EXCHANGE", Remark: "rollback after KUN failure"}
		if err := s.walletSvc.UnfreezeBalance(ctx, tx, merchantID, "FUNDING", currency, totalFreeze, ref); err != nil {
			return err
		}
		return tx.WithContext(ctx).Model(&model.TransactionRecord{}).
			Where("id = ?", txRecordID).
			Update("status", "FAILED").Error
	})
	if err != nil {
		slog.Error("rollback exchange order failed",
			slog.Any("error", err),
			slog.Uint64("tx_record_id", txRecordID),
		)
	}
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

// requestKUNQuote 暂时停用，恢复实时询价时取消注释。
// func (s *ExchangeService) requestKUNQuote(
// 	ctx context.Context,
// 	subCustomerNo, fromCurrency, toCurrency, fromAmount string,
// ) (*kundto.ExchangeQuoteResp, error) {
// 	var quote kundto.ExchangeQuoteResp
// 	err := s.kunClient.PostAsCustomer(ctx, subCustomerNo, kun.ExchangeQuoteRequestPath, &kundto.ExchangeQuoteReq{
// 		RequestNo:      kun.GenerateRequestNo(),
// 		Amount:         fromAmount,
// 		QuoteCurrency:  fromCurrency,
// 		QuotedCurrency: toCurrency,
// 	}, &quote)
// 	if err != nil {
// 		slog.Error("KUN exchange quote failed", slog.Any("error", err))
// 		return nil, bizerrors.ErrKUNAPIFailedE
// 	}
// 	if strings.TrimSpace(quote.QuoteId) == "" {
// 		slog.Error("KUN exchange quote returned empty quote id")
// 		return nil, bizerrors.ErrKUNAPIFailedE
// 	}
// 	return &quote, nil
// }

func (s *ExchangeService) queryKUNExchangeOrder(
	ctx context.Context,
	merchant *model.Merchant,
	order *model.ExchangeOrder,
) (*kundto.ExchangeOrderQueryResp, error) {
	var kunResp kundto.ExchangeOrderQueryResp
	if err := s.queryKUNExchangeOrderInto(ctx, merchant, order, &kunResp); err != nil {
		return nil, err
	}
	return &kunResp, nil
}

func (s *ExchangeService) queryKUNExchangeOrderInto(
	ctx context.Context,
	merchant *model.Merchant,
	order *model.ExchangeOrder,
	resp *kundto.ExchangeOrderQueryResp,
) error {
	if merchant.KunSubCustomerNo == nil {
		return bizerrors.ErrMerchantNotRegisteredE
	}
	if order.KunOrderID == nil || strings.TrimSpace(*order.KunOrderID) == "" {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "exchange order has no kun order id")
	}

	if order.ExchangeType == "1TO1" {
		var innerResp kundto.InnerMatchQueryResp
		err := s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, kun.InnerMatchQueryPath, &kundto.InnerMatchQueryReq{
			RequestNo: kun.GenerateRequestNo(),
			OrderId:   *order.KunOrderID,
		}, &innerResp)
		if err != nil {
			slog.Error("KUN 1:1 exchange query failed", slog.Any("error", err), slog.Uint64("order_id", order.ID))
			return bizerrors.ErrKUNAPIFailedE
		}
		*resp = kundto.ExchangeOrderQueryResp{
			OrderId:      innerResp.OrderId,
			OrderStatus:  innerResp.OrderStatus,
			FromCurrency: innerResp.FromCurrency,
			ToCurrency:   innerResp.ToCurrency,
			FromAmount:   innerResp.OrderAmount,
			ToAmount:     firstNonEmpty(innerResp.ReceiveAmount, innerResp.OrderAmount),
			ExchangeRate: innerResp.ExchangeRate,
			TradeFee:     innerResp.TradeFee,
			FeeCurrency:  innerResp.TradeFeeCurrency,
		}
		return nil
	}

	err := s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, kun.ExchangeOrderQueryPath, &kundto.ExchangeOrderQueryReq{
		SubCustomerNo: *merchant.KunSubCustomerNo,
		OrderId:       *order.KunOrderID,
	}, resp)
	if err != nil {
		slog.Error("KUN spot exchange query failed", slog.Any("error", err), slog.Uint64("order_id", order.ID))
		return bizerrors.ErrKUNAPIFailedE
	}
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func resolveExchangeFailReason(reasons ...string) string {
	for _, reason := range reasons {
		if trimmed := strings.TrimSpace(reason); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func (s *ExchangeService) refreshFailedExchangeFromKUN(
	ctx context.Context,
	order *model.ExchangeOrder,
	kunStatus, failReason string,
) error {
	reason := strings.TrimSpace(failReason)
	updates := map[string]interface{}{}
	if reason != "" {
		current := ""
		if order.FailReason != nil {
			current = *order.FailReason
		}
		if reason != current {
			updates["fail_reason"] = reason
		}
	}
	if len(updates) == 0 {
		slog.Info("re-sync failed exchange: no metadata changes",
			slog.Uint64("order_id", order.ID),
			slog.String("kun_status", kunStatus),
		)
		return nil
	}
	return s.exchangeOrderRepo.UpdateFields(ctx, order.ID, updates)
}
