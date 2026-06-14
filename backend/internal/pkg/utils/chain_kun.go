package utils

import "strings"

// KUNDepositChain maps platform chain codes to KUN deposit API values.
// History query: https://opendocs.kun.global/docs/api/crypto-deposit-history-query
// Address query: https://opendocs.kun.global/docs/api/crypto-deposit-address-list-query
func KUNDepositChain(currency, chain string) string {
	switch strings.ToUpper(strings.TrimSpace(chain)) {
	case "TRC20", "TRX_TRC20":
		return "TRX_TRC20"
	case "ERC20", "ETH_ERC20":
		return "ETH_ERC20"
	case "BTC", "BTC_BITCOIN":
		return "BTC_Bitcoin"
	case "SOL", "SOL_SOLANA":
		return "SOL_Solana"
	case "BSC_BEP20", "BEP20", "BNB_BEP20":
		return "BNB_BEP20"
	case "TON":
		return "TON"
	default:
		if chain != "" {
			return chain
		}
		if strings.ToUpper(currency) == "BTC" {
			return "BTC_Bitcoin"
		}
		return "TRX_TRC20"
	}
}
