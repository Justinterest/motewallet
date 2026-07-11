package currency

import (
	"encoding/json"
	"strings"
)

const (
	ConfigKeyCryptoChains  = "supported_crypto_chains"
	ConfigKeyDefaultChains = "default_crypto_chains"
)

// CatalogChainsByCurrency is the full universe of chains the platform can enable.
var CatalogChainsByCurrency = map[string][]string{
	"USDT": {"ETH_ERC20", "TRX_TRC20", "TON", "SOL_Solana", "BSC_BEP20"},
	"USDC": {"ETH_ERC20", "TRX_TRC20", "SOL_Solana", "BSC_BEP20"},
	"BTC":  {"BTC"},
}

var DefaultSupportedChainsJSON = `{"USDT":["ETH_ERC20","TRX_TRC20","SOL_Solana","BSC_BEP20"],"USDC":["ETH_ERC20","TRX_TRC20","SOL_Solana","BSC_BEP20"],"BTC":["BTC"]}`

var DefaultDefaultChainsJSON = `{"USDT":"TRX_TRC20","USDC":"ETH_ERC20","BTC":"BTC"}`

func ParseChainMap(value string) map[string][]string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	var raw map[string][]string
	if err := json.Unmarshal([]byte(trimmed), &raw); err != nil {
		return nil
	}
	result := make(map[string][]string, len(raw))
	for currency, chains := range raw {
		code := strings.ToUpper(strings.TrimSpace(currency))
		if code == "" || !IsCrypto(code) {
			continue
		}
		var normalized []string
		seen := make(map[string]struct{})
		for _, chain := range chains {
			item := normalizeChain(chain)
			if item == "" {
				continue
			}
			if _, ok := seen[item]; ok {
				continue
			}
			seen[item] = struct{}{}
			normalized = append(normalized, item)
		}
		if len(normalized) > 0 {
			result[code] = normalized
		}
	}
	return result
}

func ParseDefaultChainMap(value string) map[string]string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	var raw map[string]string
	if err := json.Unmarshal([]byte(trimmed), &raw); err != nil {
		return nil
	}
	result := make(map[string]string, len(raw))
	for currency, chain := range raw {
		code := strings.ToUpper(strings.TrimSpace(currency))
		item := normalizeChain(chain)
		if code == "" || item == "" || !IsCrypto(code) {
			continue
		}
		result[code] = item
	}
	return result
}

func JoinChainMap(chains map[string][]string) string {
	if len(chains) == 0 {
		return "{}"
	}
	data, err := json.Marshal(chains)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func JoinDefaultChainMap(defaults map[string]string) string {
	if len(defaults) == 0 {
		return "{}"
	}
	data, err := json.Marshal(defaults)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func CatalogChains(currencyCode string) []string {
	return append([]string(nil), CatalogChainsByCurrency[strings.ToUpper(strings.TrimSpace(currencyCode))]...)
}

func IsKnownChain(currencyCode, chain string) bool {
	return Contains(CatalogChains(currencyCode), normalizeChain(chain))
}

func FilterChains(requested, allowed []string) []string {
	if len(requested) == 0 {
		return append([]string(nil), allowed...)
	}
	var result []string
	for _, item := range requested {
		normalized := normalizeChain(item)
		if normalized != "" && Contains(allowed, normalized) && !Contains(result, normalized) {
			result = append(result, normalized)
		}
	}
	return result
}

func ValidateChainSelection(currencyCode string, requested, available []string) []string {
	catalog := CatalogChains(currencyCode)
	var result []string
	for _, item := range requested {
		normalized := normalizeChain(item)
		if normalized == "" || !Contains(catalog, normalized) || !Contains(available, normalized) {
			continue
		}
		if !Contains(result, normalized) {
			result = append(result, normalized)
		}
	}
	return result
}

func ResolveDefaultChain(currencyCode string, supported []string, configured string) string {
	if len(supported) == 0 {
		return ""
	}
	preferred := normalizeChain(configured)
	if preferred != "" && Contains(supported, preferred) {
		return preferred
	}
	return supported[0]
}

func normalizeChain(chain string) string {
	return strings.TrimSpace(chain)
}

func CloneChainMap(src map[string][]string) map[string][]string {
	if src == nil {
		return nil
	}
	dst := make(map[string][]string, len(src))
	for k, v := range src {
		dst[k] = append([]string(nil), v...)
	}
	return dst
}

func CloneDefaultChainMap(src map[string]string) map[string]string {
	if src == nil {
		return nil
	}
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
