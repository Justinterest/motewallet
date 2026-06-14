package utils

import "strings"

// ChainDisplayName maps internal chain codes to user-facing network labels.
func ChainDisplayName(chain string) string {
	switch strings.ToUpper(strings.TrimSpace(chain)) {
	case "TRX_TRC20", "TRC20":
		return "TRC20（波场）"
	case "ETH_ERC20", "ERC20":
		return "ERC20（以太坊）"
	case "BTC_Bitcoin", "BTC":
		return "Bitcoin"
	case "SOL_Solana", "SOL":
		return "Solana"
	case "BSC_BEP20", "BNB_BEP20", "BEP20":
		return "BEP20（BNB Chain）"
	case "TON":
		return "TON"
	default:
		if chain == "" {
			return ""
		}
		return chain
	}
}
