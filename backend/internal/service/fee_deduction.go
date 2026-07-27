package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"motewallet/internal/model"
)

const (
	FeeDeductionWallet         = "WALLET"
	FeeDeductionReceivedAmount = "RECEIVED_AMOUNT"
)

func normalizeFeeDeductionMethod(method string) string {
	if strings.ToUpper(strings.TrimSpace(method)) == FeeDeductionReceivedAmount {
		return FeeDeductionReceivedAmount
	}
	return FeeDeductionWallet
}

func validateFeeDeductionMethod(method string) error {
	normalized := strings.ToUpper(strings.TrimSpace(method))
	if normalized == "" || normalized == FeeDeductionWallet || normalized == FeeDeductionReceivedAmount {
		return nil
	}
	return fmt.Errorf("invalid fee deduction method: %s", method)
}

func walletDeductionAmount(amount, fee decimal.Decimal, method string) decimal.Decimal {
	if normalizeFeeDeductionMethod(method) == FeeDeductionReceivedAmount {
		return amount
	}
	return amount.Add(fee)
}

func receivedAmount(amount, fee decimal.Decimal, method string) decimal.Decimal {
	if normalizeFeeDeductionMethod(method) == FeeDeductionReceivedAmount {
		return amount.Sub(fee)
	}
	return amount
}

func feeDeductionMethods(ctx context.Context, db *gorm.DB, templateID *uint64) (exchange, crypto, fiat string) {
	exchange, crypto, fiat = FeeDeductionWallet, FeeDeductionWallet, FeeDeductionWallet
	if templateID == nil {
		return
	}
	var template model.FeeTemplate
	if err := db.WithContext(ctx).
		Select("exchange_fee_deduction_method", "crypto_withdrawal_fee_deduction_method", "fiat_withdrawal_fee_deduction_method").
		First(&template, *templateID).Error; err != nil {
		return
	}
	return normalizeFeeDeductionMethod(template.ExchangeFeeDeductionMethod),
		normalizeFeeDeductionMethod(template.CryptoWithdrawalFeeDeductionMethod),
		normalizeFeeDeductionMethod(template.FiatWithdrawalFeeDeductionMethod)
}
