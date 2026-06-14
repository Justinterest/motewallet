package utils

import "testing"

func TestKUNDepositChain(t *testing.T) {
	tests := []struct {
		currency string
		chain    string
		want     string
	}{
		{"USDT", "BSC_BEP20", "BNB_BEP20"},
		{"USDT", "BNB_BEP20", "BNB_BEP20"},
		{"USDT", "ETH_ERC20", "ETH_ERC20"},
		{"USDT", "TRX_TRC20", "TRX_TRC20"},
		{"BTC", "", "BTC_Bitcoin"},
		{"USDT", "", "TRX_TRC20"},
	}

	for _, tt := range tests {
		got := KUNDepositChain(tt.currency, tt.chain)
		if got != tt.want {
			t.Fatalf("currency=%s chain=%s got %s want %s", tt.currency, tt.chain, got, tt.want)
		}
	}
}
