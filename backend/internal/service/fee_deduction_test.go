package service

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestFeeDeductionAmounts(t *testing.T) {
	amount := decimal.NewFromInt(100)
	fee := decimal.NewFromInt(3)

	tests := []struct {
		name          string
		method        string
		walletAmount  string
		receivedValue string
	}{
		{name: "wallet", method: FeeDeductionWallet, walletAmount: "103", receivedValue: "100"},
		{name: "received amount", method: FeeDeductionReceivedAmount, walletAmount: "100", receivedValue: "97"},
		{name: "legacy empty defaults to wallet", method: "", walletAmount: "103", receivedValue: "100"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := walletDeductionAmount(amount, fee, tt.method).String(); got != tt.walletAmount {
				t.Fatalf("wallet deduction = %s, want %s", got, tt.walletAmount)
			}
			if got := receivedAmount(amount, fee, tt.method).String(); got != tt.receivedValue {
				t.Fatalf("received amount = %s, want %s", got, tt.receivedValue)
			}
		})
	}
}

func TestValidateFeeDeductionMethod(t *testing.T) {
	for _, method := range []string{"", "WALLET", "wallet", "RECEIVED_AMOUNT", " received_amount "} {
		if err := validateFeeDeductionMethod(method); err != nil {
			t.Fatalf("method %q should be valid: %v", method, err)
		}
	}
	if err := validateFeeDeductionMethod("UNKNOWN"); err == nil {
		t.Fatal("unknown method should be rejected")
	}
}
