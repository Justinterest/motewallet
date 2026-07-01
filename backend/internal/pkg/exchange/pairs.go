package exchange

import bizerrors "motewallet/internal/pkg/errors"

var supportedPairs = map[string]map[string]struct{}{
	"USDT": {"USD": {}},
	"USD":  {"USDT": {}, "USDC": {}},
	"USDC": {"USD": {}},
}

func IsSupportedPair(fromCurrency, toCurrency string) bool {
	targets, ok := supportedPairs[fromCurrency]
	if !ok {
		return false
	}
	_, ok = targets[toCurrency]
	return ok
}

func EnsureSupportedPair(fromCurrency, toCurrency string) error {
	if fromCurrency == toCurrency {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "from and to currencies must be different")
	}
	if !IsSupportedPair(fromCurrency, toCurrency) {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "unsupported exchange pair")
	}
	return nil
}

func SupportedCurrencies() []string {
	return []string{"USD", "USDT", "USDC"}
}
