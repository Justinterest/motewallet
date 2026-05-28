package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"motewallet/internal/model"
	"motewallet/internal/pkg/utils"
	kundto "motewallet/internal/pkg/kun/dto"
	"motewallet/internal/repository"
)

type WebhookService struct {
	db                    *gorm.DB
	webhookLogRepo        repository.WebhookLogRepository
	merchantRepo          repository.MerchantRepository
	walletSvc             *WalletService
	transactionRecordRepo repository.TransactionRecordRepository
	depositOrderRepo      repository.DepositOrderRepository
	withdrawalOrderRepo   repository.WithdrawalOrderRepository
	exchangeOrderRepo     repository.ExchangeOrderRepository
	transferOrderRepo     repository.TransferOrderRepository
}

func NewWebhookService(
	db *gorm.DB,
	webhookLogRepo repository.WebhookLogRepository,
	merchantRepo repository.MerchantRepository,
	walletSvc *WalletService,
	transactionRecordRepo repository.TransactionRecordRepository,
	depositOrderRepo repository.DepositOrderRepository,
	withdrawalOrderRepo repository.WithdrawalOrderRepository,
	exchangeOrderRepo repository.ExchangeOrderRepository,
	transferOrderRepo repository.TransferOrderRepository,
) *WebhookService {
	return &WebhookService{
		db:                    db,
		webhookLogRepo:        webhookLogRepo,
		merchantRepo:          merchantRepo,
		walletSvc:             walletSvc,
		transactionRecordRepo: transactionRecordRepo,
		depositOrderRepo:      depositOrderRepo,
		withdrawalOrderRepo:   withdrawalOrderRepo,
		exchangeOrderRepo:     exchangeOrderRepo,
		transferOrderRepo:     transferOrderRepo,
	}
}

func (s *WebhookService) ProcessEvent(ctx context.Context, event *kundto.WebhookEvent) error {
	existing, err := s.webhookLogRepo.FindByEventIDAndTopic(ctx, event.EventID, event.EventTopic)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil && existing.ProcessStatus == "PROCESSED" {
		slog.Info("webhook event already processed", slog.String("event_id", event.EventID), slog.String("topic", event.EventTopic))
		return nil
	}

	rawData, _ := json.Marshal(event.Data)
	var customerNo *string
	if v, ok := event.Data["subCustomerNo"]; ok {
		s := v.(string)
		customerNo = &s
	}

	var logID uint64
	if existing != nil {
		logID = existing.ID
		_ = s.webhookLogRepo.UpdateProcessStatus(ctx, logID, "PROCESSING", nil)
	} else {
		webhookLog := &model.WebhookLog{
			EventID:       event.EventID,
			EventTopic:    event.EventTopic,
			EventType:     event.EventType,
			CustomerNo:    customerNo,
			RawData:       rawData,
			ProcessStatus: "PROCESSING",
			ReceivedAt:    time.Now(),
		}
		if err := s.webhookLogRepo.Create(ctx, webhookLog); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return nil
			}
			return err
		}
		logID = webhookLog.ID
	}

	var processErr error
	switch event.EventTopic {
	case "CUSTOMER_AUTH":
		processErr = s.handleCustomerAuth(ctx, event)
	case "DIGITAL_RECHARGE":
		processErr = s.handleCryptoDeposit(ctx, event)
	case "DIGITAL_WITHDRAWAL":
		processErr = s.handleCryptoWithdrawal(ctx, event)
	case "LEGAL_WITHDRAWAL":
		processErr = s.handleFiatWithdrawal(ctx, event)
	case "DIGITAL_EXCHANGE_BUY", "DIGITAL_EXCHANGE_SELL", "LEGAL_EXCHANGE_DIGITAL", "DIGITAL_EXCHANGE_LEGAL":
		processErr = s.handleExchange(ctx, event)
	case "MAKER_MATCH_CREATE", "MAKER_MATCH_CANCEL":
		processErr = s.handleExchange(ctx, event)
	case "FUND_TRANSFER":
		processErr = s.handleFundTransfer(ctx, event)
	default:
		slog.Warn("unhandled webhook topic", slog.String("topic", event.EventTopic))
	}

	if processErr != nil {
		errMsg := processErr.Error()
		_ = s.webhookLogRepo.UpdateProcessStatus(ctx, logID, "FAILED", &errMsg)
		return processErr
	}

	_ = s.webhookLogRepo.UpdateProcessStatus(ctx, logID, "PROCESSED", nil)
	return nil
}

func (s *WebhookService) handleCustomerAuth(ctx context.Context, event *kundto.WebhookEvent) error {
	dataBytes, _ := json.Marshal(event.Data)
	var authData kundto.CustomerAuthData
	if err := json.Unmarshal(dataBytes, &authData); err != nil {
		return err
	}

	merchant, err := s.merchantRepo.FindByKunSubCustomerNo(ctx, authData.SubCustomerNo)
	if err != nil {
		return err
	}

	now := time.Now()

	switch authData.AuthStatus {
	case "SUCCESS":
		if merchant.Status == "ACTIVE" {
			return nil
		}
		err = s.merchantRepo.UpdateFields(ctx, merchant.ID, map[string]interface{}{
			"status":           "ACTIVE",
			"kyc_status":       "AUTH_SUC",
			"kyc_completed_at": now,
		})
		if err != nil {
			return err
		}
		return s.walletSvc.InitializeWallets(ctx, merchant.ID)
	case "FAIL":
		return s.merchantRepo.UpdateFields(ctx, merchant.ID, map[string]interface{}{
			"kyc_status":       "AUTH_FAIL",
			"kyc_completed_at": now,
		})
	default:
		slog.Warn("unknown CUSTOMER_AUTH status", slog.String("status", authData.AuthStatus))
	}

	return nil
}

func (s *WebhookService) handleCryptoDeposit(ctx context.Context, event *kundto.WebhookEvent) error {
	dataBytes, _ := json.Marshal(event.Data)
	var data kundto.CryptoDepositData
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return err
	}

	if data.OrderStatus != "SUCCESS" {
		slog.Info("crypto deposit not SUCCESS, skipping", slog.String("status", data.OrderStatus))
		return nil
	}

	existing, err := s.depositOrderRepo.FindByKunOrderID(ctx, data.OrderId)
	if err == nil && existing != nil {
		return nil
	}

	merchant, err := s.merchantRepo.FindByKunSubCustomerNo(ctx, data.SubCustomerNo)
	if err != nil {
		return err
	}

	amount, _ := decimal.NewFromString(data.OrderAmount)

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		platformOrderID := utils.GeneratePlatformOrderID("DEP")
		subType := "CRYPTO"
		txRecord := &model.TransactionRecord{
			PlatformOrderID: platformOrderID,
			MerchantID:      merchant.ID,
			Type:            "DEPOSIT",
			SubType:         &subType,
			Amount:          amount,
			Currency:        data.Currency,
			PlatformFee:     decimal.Zero,
			Status:          "COMPLETED",
		}
		if err := tx.WithContext(ctx).Create(txRecord).Error; err != nil {
			return err
		}

		depositOrder := &model.DepositOrder{
			TransactionRecordID: txRecord.ID,
			MerchantID:          merchant.ID,
			Currency:            data.Currency,
			Chain:               data.Chain,
			ToAddress:           data.ToAddress,
			FromAddress:         &data.FromAddress,
			TxID:                &data.TxId,
			KunOrderID:          &data.OrderId,
		}
		if err := tx.WithContext(ctx).Create(depositOrder).Error; err != nil {
			return err
		}

		return s.walletSvc.CreditBalance(ctx, tx, merchant.ID, "FUNDING", data.Currency, amount)
	})
}

func (s *WebhookService) handleCryptoWithdrawal(ctx context.Context, event *kundto.WebhookEvent) error {
	dataBytes, _ := json.Marshal(event.Data)
	var data kundto.CryptoWithdrawalData
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return err
	}

	order, err := s.withdrawalOrderRepo.FindByKunRequestNo(ctx, data.RequestNo)
	if err != nil {
		return err
	}

	txRecord, err := s.transactionRecordRepo.FindByID(ctx, order.TransactionRecordID)
	if err != nil {
		return err
	}

	if txRecord.Status == "COMPLETED" || txRecord.Status == "FAILED" {
		return nil
	}

	kunFee, _ := decimal.NewFromString(data.FeeAmount)
	totalFrozen := txRecord.Amount.Add(txRecord.PlatformFee)

	switch data.OrderStatus {
	case "SUCCESS":
		return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := s.walletSvc.DeductFrozen(ctx, tx, order.MerchantID, "FUNDING", txRecord.Currency, totalFrozen); err != nil {
				return err
			}
			if err := s.transactionRecordRepo.UpdateStatus(ctx, txRecord.ID, "COMPLETED"); err != nil {
				return err
			}
			return s.withdrawalOrderRepo.UpdateFields(ctx, order.ID, map[string]interface{}{
				"tx_id":    data.TxId,
				"kun_fee":  kunFee,
			})
		})
	case "FAIL":
		return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := s.walletSvc.UnfreezeBalance(ctx, tx, order.MerchantID, "FUNDING", txRecord.Currency, totalFrozen); err != nil {
				return err
			}
			return s.transactionRecordRepo.UpdateStatus(ctx, txRecord.ID, "FAILED")
		})
	default:
		slog.Info("crypto withdrawal status update", slog.String("status", data.OrderStatus))
	}

	return nil
}

func (s *WebhookService) handleFiatWithdrawal(ctx context.Context, event *kundto.WebhookEvent) error {
	dataBytes, _ := json.Marshal(event.Data)
	var data kundto.FiatWithdrawalData
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return err
	}

	order, err := s.withdrawalOrderRepo.FindByKunRequestNo(ctx, data.RequestNo)
	if err != nil {
		return err
	}

	txRecord, err := s.transactionRecordRepo.FindByID(ctx, order.TransactionRecordID)
	if err != nil {
		return err
	}

	if txRecord.Status == "COMPLETED" || txRecord.Status == "FAILED" {
		return nil
	}

	kunFee, _ := decimal.NewFromString(data.FeeAmount)
	totalFrozen := txRecord.Amount.Add(txRecord.PlatformFee)

	switch data.OrderStatus {
	case "SUCCESS":
		return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := s.walletSvc.DeductFrozen(ctx, tx, order.MerchantID, "FUNDING", txRecord.Currency, totalFrozen); err != nil {
				return err
			}
			if err := s.transactionRecordRepo.UpdateStatus(ctx, txRecord.ID, "COMPLETED"); err != nil {
				return err
			}
			return s.withdrawalOrderRepo.UpdateFields(ctx, order.ID, map[string]interface{}{
				"kun_fee": kunFee,
			})
		})
	case "FAIL":
		return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := s.walletSvc.UnfreezeBalance(ctx, tx, order.MerchantID, "FUNDING", txRecord.Currency, totalFrozen); err != nil {
				return err
			}
			return s.transactionRecordRepo.UpdateStatus(ctx, txRecord.ID, "FAILED")
		})
	default:
		slog.Info("fiat withdrawal status update", slog.String("status", data.OrderStatus))
	}

	return nil
}

func (s *WebhookService) handleExchange(ctx context.Context, event *kundto.WebhookEvent) error {
	dataBytes, _ := json.Marshal(event.Data)
	var data kundto.ExchangeData
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return err
	}

	order, err := s.exchangeOrderRepo.FindByKunRequestNo(ctx, data.RequestNo)
	if err != nil {
		return err
	}

	txRecord, err := s.transactionRecordRepo.FindByID(ctx, order.TransactionRecordID)
	if err != nil {
		return err
	}

	if txRecord.Status == "COMPLETED" || txRecord.Status == "FAILED" {
		return nil
	}

	toAmount, _ := decimal.NewFromString(data.ToAmount)
	exchangeRate, _ := decimal.NewFromString(data.ExchangeRate)
	kunFee, _ := decimal.NewFromString(data.TradeFee)

	switch data.OrderStatus {
	case "SUCCESS":
		if err := s.transactionRecordRepo.UpdateStatus(ctx, txRecord.ID, "COMPLETED"); err != nil {
			return err
		}
		return s.exchangeOrderRepo.UpdateFields(ctx, order.ID, map[string]interface{}{
			"to_amount":     toAmount,
			"exchange_rate": exchangeRate,
			"kun_fee":       kunFee,
		})
	case "FAIL":
		return s.transactionRecordRepo.UpdateStatus(ctx, txRecord.ID, "FAILED")
	default:
		slog.Info("exchange status update", slog.String("status", data.OrderStatus), slog.String("topic", event.EventTopic))
	}

	return nil
}

func (s *WebhookService) handleFundTransfer(ctx context.Context, event *kundto.WebhookEvent) error {
	requestNo, ok := event.Data["requestNo"].(string)
	if !ok {
		return nil
	}
	orderStatus, _ := event.Data["orderStatus"].(string)

	order, err := s.transferOrderRepo.FindByKunRequestNo(ctx, requestNo)
	if err != nil {
		return err
	}

	txRecord, err := s.transactionRecordRepo.FindByID(ctx, order.TransactionRecordID)
	if err != nil {
		return err
	}

	if txRecord.Status == "COMPLETED" || txRecord.Status == "FAILED" {
		return nil
	}

	switch orderStatus {
	case "SUCCESS":
		return s.transactionRecordRepo.UpdateStatus(ctx, txRecord.ID, "COMPLETED")
	case "FAIL":
		return s.transactionRecordRepo.UpdateStatus(ctx, txRecord.ID, "FAILED")
	default:
		slog.Info("fund transfer status update", slog.String("status", orderStatus))
	}

	return nil
}
