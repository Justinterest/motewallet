package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
	dtoresp "motewallet/internal/dto/response"
	"motewallet/internal/pkg/currency"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/kun"
	kundto "motewallet/internal/pkg/kun/dto"
)

const (
	kunRegionBalancePath = "/rest/v2.0/account/query/balance"
	kunTradeBalancePath  = "/rest/v2.0/trade/account/outAccount/query"
)

func (s *MerchantManagementService) SyncKUNBalances(ctx context.Context, adminID, merchantID uint64) (*dtoresp.SyncKUNBalancesResp, error) {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	if merchant.KunSubCustomerNo == nil || strings.TrimSpace(*merchant.KunSubCustomerNo) == "" {
		return nil, bizerrors.ErrMerchantNotRegisteredE
	}

	supportedCrypto, err := s.currencyConfigSvc.GetSupportedCrypto(ctx, merchant)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}
	supportedFiat, err := s.currencyConfigSvc.GetSupportedFiat(ctx, merchant)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}
	orderedCurrencies := append(supportedFiat, supportedCrypto...)

	subCustomerNo := *merchant.KunSubCustomerNo
	var balances []dtoresp.KUNWalletBalanceResp

	var fundingItems []kundto.AccountBalanceItem
	err = s.kunClient.PostAsCustomer(ctx, subCustomerNo, kunRegionBalancePath, &kundto.RegionBalanceQueryReq{
		RequestNo:  kun.GenerateRequestNo(),
		RegionCode: kun.AccountHK,
	}, &fundingItems)
	if err != nil {
		slog.Error("KUN query regional account balance failed", slog.Any("error", err), slog.Uint64("merchant_id", merchantID))
		return nil, bizerrors.ErrKUNAPIFailedE
	}

	for _, item := range fundingItems {
		if !currency.Contains(orderedCurrencies, item.Currency) {
			continue
		}
		balances = append(balances, mapKUNBalanceItem("FUNDING", item))
	}

	for _, code := range orderedCurrencies {
		var tradeItems []kundto.AccountBalanceItem
		err = s.kunClient.PostAsCustomer(ctx, subCustomerNo, kunTradeBalancePath, &kundto.OutAccountBalanceQueryReq{
			RequestNo:    kun.GenerateRequestNo(),
			Currency:     code,
			CurrencyType: kunCurrencyType(code),
		}, &tradeItems)
		if err != nil {
			slog.Warn("KUN query account balance failed",
				slog.Any("error", err),
				slog.Uint64("merchant_id", merchantID),
				slog.String("currency", code),
			)
			continue
		}
		for _, item := range tradeItems {
			if !currency.Contains(orderedCurrencies, item.Currency) {
				continue
			}
			balances = append(balances, mapKUNBalanceItem("TRADING", item))
		}
	}

	syncedAt := time.Now()
	s.logAudit(ctx, adminID, "SYNC_KUN_BALANCES", "Merchant", fmt.Sprintf("%d", merchantID), map[string]any{
		"balance_count": len(balances),
	})

	return &dtoresp.SyncKUNBalancesResp{
		KUNBalances: sortKUNBalances(balances, orderedCurrencies),
		SyncedAt:    syncedAt,
	}, nil
}

func mapKUNBalanceItem(accountType string, item kundto.AccountBalanceItem) dtoresp.KUNWalletBalanceResp {
	balance := defaultAmount(item.Balance, "0")
	return dtoresp.KUNWalletBalanceResp{
		AccountType: accountType,
		Currency:    strings.ToUpper(item.Currency),
		Balance:     balance,
	}
}

func kunCurrencyType(code string) string {
	if currency.IsCrypto(code) {
		return "DIGITAL"
	}
	return "LEGAL"
}

func defaultAmount(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return "0"
}

func sortKUNBalances(balances []dtoresp.KUNWalletBalanceResp, orderedCurrencies []string) []dtoresp.KUNWalletBalanceResp {
	index := make(map[string]int, len(orderedCurrencies))
	for i, code := range orderedCurrencies {
		index[code] = i
	}

	accountRank := map[string]int{"FUNDING": 0, "TRADING": 1}
	sorted := append([]dtoresp.KUNWalletBalanceResp(nil), balances...)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if kunBalanceLess(sorted[j], sorted[i], index, accountRank) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}

func kunBalanceLess(left, right dtoresp.KUNWalletBalanceResp, currencyIndex map[string]int, accountRank map[string]int) bool {
	leftAccount := accountRank[left.AccountType]
	rightAccount := accountRank[right.AccountType]
	if leftAccount != rightAccount {
		return leftAccount < rightAccount
	}
	return currencyIndex[left.Currency] < currencyIndex[right.Currency]
}
