package currency

import (
	"strings"
)

const (
	ConfigKeyCrypto = "supported_crypto_currencies"
	ConfigKeyFiat   = "supported_fiat_currencies"

	DefaultCryptoCSV = "USDT,USDC,BTC"
	DefaultFiatCSV   = "USD,HKD,EUR"
)

var (
	AllCrypto = []string{"USDT", "USDC", "BTC"}
	AllFiat   = []string{"USD", "HKD", "EUR"}
)

func ParseList(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	var result []string
	seen := make(map[string]struct{})
	for _, part := range parts {
		item := strings.ToUpper(strings.TrimSpace(part))
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

func JoinList(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return strings.Join(items, ",")
}

func Contains(items []string, item string) bool {
	target := strings.ToUpper(strings.TrimSpace(item))
	for _, current := range items {
		if strings.ToUpper(current) == target {
			return true
		}
	}
	return false
}

func IsCrypto(item string) bool {
	return Contains(AllCrypto, item)
}

func IsFiat(item string) bool {
	return Contains(AllFiat, item)
}

func FilterAllowed(requested, allowed []string) []string {
	if len(requested) == 0 {
		return allowed
	}
	var result []string
	for _, item := range requested {
		if Contains(allowed, item) {
			result = append(result, strings.ToUpper(item))
		}
	}
	return result
}

func ValidateSelection(requested, available []string) []string {
	var result []string
	for _, item := range requested {
		normalized := strings.ToUpper(strings.TrimSpace(item))
		if normalized == "" || !Contains(available, normalized) {
			continue
		}
		if !Contains(result, normalized) {
			result = append(result, normalized)
		}
	}
	return result
}
