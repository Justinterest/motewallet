package service

import (
	"context"
	"strings"

	"motewallet/internal/model"
	"motewallet/internal/pkg/currency"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/repository"
)

type CurrencyConfigService struct {
	systemConfigRepo repository.SystemConfigRepository
}

func NewCurrencyConfigService(systemConfigRepo repository.SystemConfigRepository) *CurrencyConfigService {
	return &CurrencyConfigService{
		systemConfigRepo: systemConfigRepo,
	}
}

func (s *CurrencyConfigService) GetAvailableCrypto(ctx context.Context) ([]string, error) {
	value, err := s.systemConfigRepo.GetByKey(ctx, currency.ConfigKeyCrypto)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(value) == "" {
		return currency.ParseList(currency.DefaultCryptoCSV), nil
	}
	return currency.ParseList(value), nil
}

func (s *CurrencyConfigService) GetAvailableFiat(ctx context.Context) ([]string, error) {
	value, err := s.systemConfigRepo.GetByKey(ctx, currency.ConfigKeyFiat)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(value) == "" {
		return currency.ParseList(currency.DefaultFiatCSV), nil
	}
	return currency.ParseList(value), nil
}

func (s *CurrencyConfigService) GetAvailableChains(ctx context.Context) (map[string][]string, error) {
	value, err := s.systemConfigRepo.GetByKey(ctx, currency.ConfigKeyCryptoChains)
	if err != nil {
		return nil, err
	}
	raw := currency.ParseChainMap(value)
	if len(raw) == 0 {
		raw = currency.ParseChainMap(currency.DefaultSupportedChainsJSON)
	}
	return s.filterChainsByCurrencies(raw, currency.AllCrypto), nil
}

func (s *CurrencyConfigService) GetAvailableDefaultChains(ctx context.Context) (map[string]string, error) {
	value, err := s.systemConfigRepo.GetByKey(ctx, currency.ConfigKeyDefaultChains)
	if err != nil {
		return nil, err
	}
	availableChains, err := s.GetAvailableChains(ctx)
	if err != nil {
		return nil, err
	}
	raw := currency.ParseDefaultChainMap(value)
	if len(raw) == 0 {
		raw = currency.ParseDefaultChainMap(currency.DefaultDefaultChainsJSON)
	}
	return s.resolveDefaults(availableChains, raw), nil
}

func (s *CurrencyConfigService) GetSupportedCrypto(ctx context.Context, merchant *model.Merchant) ([]string, error) {
	if merchant == nil || merchant.SupportedCryptoCurrencies == nil || strings.TrimSpace(*merchant.SupportedCryptoCurrencies) == "" {
		return s.GetAvailableCrypto(ctx)
	}
	return currency.FilterAllowed(currency.ParseList(*merchant.SupportedCryptoCurrencies), currency.AllCrypto), nil
}

func (s *CurrencyConfigService) GetSupportedFiat(ctx context.Context, merchant *model.Merchant) ([]string, error) {
	if merchant == nil || merchant.SupportedFiatCurrencies == nil || strings.TrimSpace(*merchant.SupportedFiatCurrencies) == "" {
		return s.GetAvailableFiat(ctx)
	}
	return currency.FilterAllowed(currency.ParseList(*merchant.SupportedFiatCurrencies), currency.AllFiat), nil
}

func (s *CurrencyConfigService) GetSupportedChains(ctx context.Context, merchant *model.Merchant) (map[string][]string, error) {
	supportedCrypto, err := s.GetSupportedCrypto(ctx, merchant)
	if err != nil {
		return nil, err
	}

	if merchant == nil || merchant.SupportedCryptoChains == nil || strings.TrimSpace(*merchant.SupportedCryptoChains) == "" {
		available, err := s.GetAvailableChains(ctx)
		if err != nil {
			return nil, err
		}
		return s.filterChainsByCurrencies(available, supportedCrypto), nil
	}

	merchantChains := currency.ParseChainMap(*merchant.SupportedCryptoChains)
	available, err := s.GetAvailableChains(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string][]string)
	for _, code := range supportedCrypto {
		allowed := available[code]
		if len(allowed) == 0 {
			continue
		}
		selected := currency.FilterChains(merchantChains[code], allowed)
		if len(selected) == 0 {
			selected = append([]string(nil), allowed...)
		}
		result[code] = selected
	}
	return result, nil
}

func (s *CurrencyConfigService) GetDefaultChains(ctx context.Context, merchant *model.Merchant) (map[string]string, error) {
	supportedChains, err := s.GetSupportedChains(ctx, merchant)
	if err != nil {
		return nil, err
	}
	globalDefaults, err := s.GetAvailableDefaultChains(ctx)
	if err != nil {
		return nil, err
	}

	configured := currency.CloneDefaultChainMap(globalDefaults)
	if merchant != nil && merchant.DefaultCryptoChains != nil && strings.TrimSpace(*merchant.DefaultCryptoChains) != "" {
		merchantDefaults := currency.ParseDefaultChainMap(*merchant.DefaultCryptoChains)
		for code, chain := range merchantDefaults {
			configured[code] = chain
		}
	}
	return s.resolveDefaults(supportedChains, configured), nil
}

func (s *CurrencyConfigService) EnsureCurrencySupported(ctx context.Context, merchant *model.Merchant, code string) error {
	normalized := strings.ToUpper(strings.TrimSpace(code))
	if normalized == "" {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "currency is required")
	}

	switch {
	case currency.IsCrypto(normalized):
		supported, err := s.GetSupportedCrypto(ctx, merchant)
		if err != nil {
			return bizerrors.ErrInternalError
		}
		if !currency.Contains(supported, normalized) {
			return bizerrors.ErrUnsupportedCurrencyE
		}
	case currency.IsFiat(normalized):
		supported, err := s.GetSupportedFiat(ctx, merchant)
		if err != nil {
			return bizerrors.ErrInternalError
		}
		if !currency.Contains(supported, normalized) {
			return bizerrors.ErrUnsupportedCurrencyE
		}
	default:
		return bizerrors.ErrUnsupportedCurrencyE
	}

	return nil
}

func (s *CurrencyConfigService) EnsureChainSupported(ctx context.Context, merchant *model.Merchant, currencyCode, chain string) error {
	normalizedCurrency := strings.ToUpper(strings.TrimSpace(currencyCode))
	normalizedChain := strings.TrimSpace(chain)
	if normalizedCurrency == "" {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "currency is required")
	}
	if normalizedChain == "" {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "chain is required")
	}
	if err := s.EnsureCurrencySupported(ctx, merchant, normalizedCurrency); err != nil {
		return err
	}
	if !currency.IsCrypto(normalizedCurrency) {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "chain is only applicable to crypto currencies")
	}

	supportedChains, err := s.GetSupportedChains(ctx, merchant)
	if err != nil {
		return bizerrors.ErrInternalError
	}
	chains := supportedChains[normalizedCurrency]
	if !currency.Contains(chains, normalizedChain) {
		// Accept BTC_Bitcoin when platform chain is BTC (deposit vs withdraw code divergence).
		if normalizedCurrency == "BTC" && strings.EqualFold(normalizedChain, "BTC_Bitcoin") && currency.Contains(chains, "BTC") {
			return nil
		}
		return bizerrors.ErrUnsupportedChainE
	}
	return nil
}

func (s *CurrencyConfigService) NormalizeMerchantSelection(_ context.Context, crypto, fiat []string) ([]string, []string, error) {
	selectedCrypto := currency.ValidateSelection(crypto, currency.AllCrypto)
	selectedFiat := currency.ValidateSelection(fiat, currency.AllFiat)
	if len(selectedCrypto) == 0 {
		return nil, nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "at least one crypto currency must be selected")
	}
	if len(selectedFiat) == 0 {
		return nil, nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "at least one fiat currency must be selected")
	}

	return selectedCrypto, selectedFiat, nil
}

func (s *CurrencyConfigService) NormalizeMerchantChainSelection(
	ctx context.Context,
	selectedCrypto []string,
	chains map[string][]string,
	defaults map[string]string,
) (map[string][]string, map[string]string, error) {
	availableChains, err := s.GetAvailableChains(ctx)
	if err != nil {
		return nil, nil, err
	}

	normalizedChains := make(map[string][]string)
	for _, code := range selectedCrypto {
		available := availableChains[code]
		if len(available) == 0 {
			continue
		}
		requested := chains[code]
		var selected []string
		if len(requested) == 0 {
			selected = append([]string(nil), available...)
		} else {
			selected = currency.ValidateChainSelection(code, requested, available)
		}
		if len(selected) == 0 {
			return nil, nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "at least one chain must be selected for "+code)
		}
		normalizedChains[code] = selected
	}

	normalizedDefaults := make(map[string]string)
	for _, code := range selectedCrypto {
		supported := normalizedChains[code]
		if len(supported) == 0 {
			continue
		}
		normalizedDefaults[code] = currency.ResolveDefaultChain(code, supported, defaults[code])
	}

	return normalizedChains, normalizedDefaults, nil
}

func (s *CurrencyConfigService) NormalizeGlobalSelection(
	_ context.Context,
	crypto, fiat []string,
	chains map[string][]string,
	defaults map[string]string,
) ([]string, []string, map[string][]string, map[string]string, error) {
	selectedCrypto := currency.ValidateSelection(crypto, currency.AllCrypto)
	selectedFiat := currency.ValidateSelection(fiat, currency.AllFiat)
	if len(selectedCrypto) == 0 {
		return nil, nil, nil, nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "at least one crypto currency must be selected")
	}
	if len(selectedFiat) == 0 {
		return nil, nil, nil, nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "at least one fiat currency must be selected")
	}

	normalizedChains := make(map[string][]string)
	for _, code := range currency.AllCrypto {
		catalog := currency.CatalogChains(code)
		selected := currency.ValidateChainSelection(code, chains[code], catalog)
		if len(selected) == 0 {
			return nil, nil, nil, nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "at least one chain must be selected for "+code)
		}
		normalizedChains[code] = selected
	}

	normalizedDefaults := make(map[string]string)
	for _, code := range currency.AllCrypto {
		normalizedDefaults[code] = currency.ResolveDefaultChain(code, normalizedChains[code], defaults[code])
	}

	return selectedCrypto, selectedFiat, normalizedChains, normalizedDefaults, nil
}

func (s *CurrencyConfigService) SerializeSelection(crypto, fiat []string) (*string, *string) {
	cryptoCSV := currency.JoinList(crypto)
	fiatCSV := currency.JoinList(fiat)
	return &cryptoCSV, &fiatCSV
}

func (s *CurrencyConfigService) SerializeChainSelection(chains map[string][]string, defaults map[string]string) (*string, *string) {
	chainsJSON := currency.JoinChainMap(chains)
	defaultsJSON := currency.JoinDefaultChainMap(defaults)
	return &chainsJSON, &defaultsJSON
}

func (s *CurrencyConfigService) UpdateGlobalConfig(
	ctx context.Context,
	crypto, fiat []string,
	chains map[string][]string,
	defaults map[string]string,
) error {
	selectedCrypto, selectedFiat, normalizedChains, normalizedDefaults, err := s.NormalizeGlobalSelection(ctx, crypto, fiat, chains, defaults)
	if err != nil {
		return err
	}

	cryptoCSV := currency.JoinList(selectedCrypto)
	fiatCSV := currency.JoinList(selectedFiat)
	chainsJSON := currency.JoinChainMap(normalizedChains)
	defaultsJSON := currency.JoinDefaultChainMap(normalizedDefaults)

	if err := s.systemConfigRepo.Upsert(ctx, currency.ConfigKeyCrypto, cryptoCSV, "Supported crypto currencies"); err != nil {
		return bizerrors.ErrInternalError
	}
	if err := s.systemConfigRepo.Upsert(ctx, currency.ConfigKeyFiat, fiatCSV, "Supported fiat currencies"); err != nil {
		return bizerrors.ErrInternalError
	}
	if err := s.systemConfigRepo.Upsert(ctx, currency.ConfigKeyCryptoChains, chainsJSON, "Default supported chains per crypto currency (JSON)"); err != nil {
		return bizerrors.ErrInternalError
	}
	if err := s.systemConfigRepo.Upsert(ctx, currency.ConfigKeyDefaultChains, defaultsJSON, "Default selected chain per crypto currency (JSON)"); err != nil {
		return bizerrors.ErrInternalError
	}
	return nil
}

func (s *CurrencyConfigService) GetCatalogChains() map[string][]string {
	return currency.CloneChainMap(currency.CatalogChainsByCurrency)
}

func (s *CurrencyConfigService) filterChainsByCurrencies(chains map[string][]string, currencies []string) map[string][]string {
	result := make(map[string][]string)
	for _, code := range currencies {
		items := chains[code]
		if len(items) == 0 {
			items = currency.CatalogChains(code)
		}
		filtered := currency.ValidateChainSelection(code, items, currency.CatalogChains(code))
		if len(filtered) > 0 {
			result[code] = filtered
		}
	}
	return result
}

func (s *CurrencyConfigService) resolveDefaults(supported map[string][]string, configured map[string]string) map[string]string {
	result := make(map[string]string)
	for code, chains := range supported {
		result[code] = currency.ResolveDefaultChain(code, chains, configured[code])
	}
	return result
}
