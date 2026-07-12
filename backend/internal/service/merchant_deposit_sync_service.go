package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	dtoresp "motewallet/internal/dto/response"
	"motewallet/internal/model"
	bizerrors "motewallet/internal/pkg/errors"
	kundto "motewallet/internal/pkg/kun/dto"
	"motewallet/internal/pkg/utils"
)

const kunDepositSyncPageSize = 50

type syncDepositStats struct {
	Created int
	Updated int
	Skipped int
}

func (s *MerchantManagementService) SyncDeposits(ctx context.Context, adminID, merchantID uint64) (*dtoresp.SyncDepositsResp, error) {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	if merchant.KunSubCustomerNo == nil || *merchant.KunSubCustomerNo == "" {
		return nil, bizerrors.ErrMerchantNotRegisteredE
	}

	records, err := s.fetchAllKUNDepositRecords(ctx, *merchant.KunSubCustomerNo)
	if err != nil {
		slog.Error("KUN fetch deposit history failed", slog.Any("error", err), slog.Uint64("merchant_id", merchantID))
		return nil, bizerrors.ErrKUNAPIFailedE
	}

	stats := &syncDepositStats{}
	for _, record := range records {
		if err := s.syncSingleDepositRecord(ctx, merchant.ID, record, stats); err != nil {
			slog.Error("sync deposit record failed",
				slog.Any("error", err),
				slog.Uint64("merchant_id", merchantID),
				slog.String("kun_order_id", record.OrderID),
			)
			return nil, bizerrors.ErrInternalError
		}
	}

	syncedAt := time.Now()
	s.logAudit(ctx, adminID, "SYNC_DEPOSITS", "Merchant", fmt.Sprintf("%d", merchantID), map[string]any{
		"total_fetched": len(records),
		"created":       stats.Created,
		"updated":       stats.Updated,
		"skipped":       stats.Skipped,
	})

	return &dtoresp.SyncDepositsResp{
		SyncedCount:  stats.Created,
		UpdatedCount: stats.Updated,
		SkippedCount: stats.Skipped,
		TotalFetched: len(records),
		SyncedAt:     syncedAt,
	}, nil
}

func (s *MerchantManagementService) fetchAllKUNDepositRecords(ctx context.Context, subCustomerNo string) ([]kundto.DepositHistoryRecord, error) {
	loc := kunDepositLocation()
	now := time.Now().In(loc)
	startTime := now.AddDate(0, 0, -kunDepositHistoryMaxDays).Format("2006-01-02 15:04:05")
	endTime := now.Format("2006-01-02 15:04:05")

	var all []kundto.DepositHistoryRecord
	page := 1

	for {
		var resp kundto.DepositHistoryQueryResp
		err := s.kunClient.PostAsCustomer(ctx, subCustomerNo, kunDepositHistoryPath, kundto.DepositHistoryQueryReq{
			StartTime: startTime,
			EndTime:   endTime,
			PageNo:    page,
			PageSize:  kunDepositSyncPageSize,
		}, &resp)
		if err != nil {
			return nil, err
		}

		items := resp.Items()
		all = append(all, items...)

		totalPage := resp.TotalPage.Int()
		if len(items) == 0 {
			break
		}
		if totalPage > 0 {
			if page >= totalPage {
				break
			}
		} else if len(items) < kunDepositSyncPageSize {
			break
		}
		page++
	}

	return all, nil
}

func (s *MerchantManagementService) syncSingleDepositRecord(
	ctx context.Context,
	merchantID uint64,
	record kundto.DepositHistoryRecord,
	stats *syncDepositStats,
) error {
	if record.OrderID == "" {
		stats.Skipped++
		return nil
	}

	status := mapDepositOrderStatus(record.OrderStatus)
	existing, err := s.depositOrderRepo.FindByKunOrderID(ctx, record.OrderID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil {
		updated, err := s.updateExistingDepositRecord(ctx, merchantID, existing, record, status)
		if err != nil {
			return err
		}
		if updated {
			stats.Updated++
		} else {
			stats.Skipped++
		}
		return nil
	}

	if err := s.createDepositRecord(ctx, merchantID, record, status); err != nil {
		return err
	}
	stats.Created++
	return nil
}

func (s *MerchantManagementService) createDepositRecord(
	ctx context.Context,
	merchantID uint64,
	record kundto.DepositHistoryRecord,
	status string,
) error {
	amount, err := decimal.NewFromString(record.OrderAmount)
	if err != nil {
		return fmt.Errorf("parse amount: %w", err)
	}

	orderTime := parseKunDateTime(record.OrderTime)
	completedAt := orderCompletedAt(record, status, orderTime)
	kunOrderID := record.OrderID

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		subType := "CRYPTO"
		txRecord := &model.TransactionRecord{
			PlatformOrderID: utils.GeneratePlatformOrderID("DEP"),
			MerchantID:      merchantID,
			Type:            "DEPOSIT",
			SubType:         &subType,
			Amount:          amount,
			Currency:        record.OrderCurrency,
			PlatformFee:     decimal.Zero,
			Status:          status,
			CompletedAt:     completedAt,
		}
		if err := tx.WithContext(ctx).Create(txRecord).Error; err != nil {
			return err
		}

		depositOrder := &model.DepositOrder{
			TransactionRecordID: txRecord.ID,
			MerchantID:          merchantID,
			Currency:            record.OrderCurrency,
			Chain:               record.Chain,
			ToAddress:           defaultDepositAddress(record.ReceiveWalletAddress),
			FromAddress:         nullableString(record.SendWalletAddress),
			TxID:                nullableString(record.TxID),
			KunOrderID:          &kunOrderID,
			Confirmations:       parseDepositConfirmations(record.NetworkConfirmNumber),
			CompletedAt:         completedAt,
		}
		if err := tx.WithContext(ctx).Create(depositOrder).Error; err != nil {
			return err
		}

		if status == "COMPLETED" {
			ref := WalletChangeRef{TransactionRecordID: txRecord.ID, BizType: "DEPOSIT"}
			return s.walletSvc.CreditBalance(ctx, tx, merchantID, "FUNDING", record.OrderCurrency, amount, ref)
		}
		return nil
	})
}

func (s *MerchantManagementService) updateExistingDepositRecord(
	ctx context.Context,
	merchantID uint64,
	existing *model.DepositOrder,
	record kundto.DepositHistoryRecord,
	status string,
) (bool, error) {
	txRecord, err := s.transactionRecordRepo.FindByID(ctx, existing.TransactionRecordID)
	if err != nil {
		return false, err
	}

	if txRecord.Status == status && status != "COMPLETED" {
		return false, nil
	}
	if txRecord.Status == "COMPLETED" {
		return false, nil
	}
	if status != "COMPLETED" {
		if txRecord.Status == status {
			return false, nil
		}
		if err := s.transactionRecordRepo.UpdateStatus(ctx, txRecord.ID, status); err != nil {
			return false, err
		}
		return true, nil
	}

	amount, err := decimal.NewFromString(record.OrderAmount)
	if err != nil {
		return false, fmt.Errorf("parse amount: %w", err)
	}
	completedAt := orderCompletedAt(record, status, parseKunDateTime(record.OrderTime))

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Model(&model.TransactionRecord{}).Where("id = ?", txRecord.ID).Updates(map[string]any{
			"status":       "COMPLETED",
			"completed_at": completedAt,
		}).Error; err != nil {
			return err
		}

		depositUpdates := map[string]any{
			"confirmations": parseDepositConfirmations(record.NetworkConfirmNumber),
			"completed_at":  completedAt,
		}
		if record.TxID != "" {
			depositUpdates["tx_id"] = record.TxID
		}
		if record.SendWalletAddress != "" {
			depositUpdates["from_address"] = record.SendWalletAddress
		}
		if err := tx.WithContext(ctx).Model(&model.DepositOrder{}).Where("id = ?", existing.ID).Updates(depositUpdates).Error; err != nil {
			return err
		}

		ref := WalletChangeRef{TransactionRecordID: txRecord.ID, BizType: "DEPOSIT"}
		return s.walletSvc.CreditBalance(ctx, tx, merchantID, "FUNDING", record.OrderCurrency, amount, ref)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func orderCompletedAt(record kundto.DepositHistoryRecord, status string, fallback time.Time) *time.Time {
	if status != "COMPLETED" {
		return nil
	}
	if record.UpdateTime != "" {
		t := parseKunDateTime(record.UpdateTime)
		return &t
	}
	return &fallback
}

func defaultDepositAddress(address string) string {
	if address == "" {
		return "-"
	}
	return address
}

func parseDepositConfirmations(value string) *uint {
	if value == "" {
		return nil
	}
	n, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return nil
	}
	v := uint(n)
	return &v
}
