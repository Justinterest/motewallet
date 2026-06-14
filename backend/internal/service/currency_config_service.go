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

func (s *CurrencyConfigService) GetSupportedCrypto(ctx context.Context, merchant *model.Merchant) ([]string, error) {
	available, err := s.GetAvailableCrypto(ctx)
	if err != nil {
		return nil, err
	}
	if merchant == nil || merchant.SupportedCryptoCurrencies == nil || strings.TrimSpace(*merchant.SupportedCryptoCurrencies) == "" {
		return available, nil
	}
	return currency.FilterAllowed(currency.ParseList(*merchant.SupportedCryptoCurrencies), available), nil
}

func (s *CurrencyConfigService) GetSupportedFiat(ctx context.Context, merchant *model.Merchant) ([]string, error) {
	available, err := s.GetAvailableFiat(ctx)
	if err != nil {
		return nil, err
	}
	if merchant == nil || merchant.SupportedFiatCurrencies == nil || strings.TrimSpace(*merchant.SupportedFiatCurrencies) == "" {
		return available, nil
	}
	return currency.FilterAllowed(currency.ParseList(*merchant.SupportedFiatCurrencies), available), nil
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

func (s *CurrencyConfigService) NormalizeMerchantSelection(ctx context.Context, crypto, fiat []string) ([]string, []string, error) {
	availableCrypto, err := s.GetAvailableCrypto(ctx)
	if err != nil {
		return nil, nil, err
	}
	availableFiat, err := s.GetAvailableFiat(ctx)
	if err != nil {
		return nil, nil, err
	}

	selectedCrypto := currency.ValidateSelection(crypto, availableCrypto)
	selectedFiat := currency.ValidateSelection(fiat, availableFiat)
	if len(selectedCrypto) == 0 {
		return nil, nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "at least one crypto currency must be selected")
	}
	if len(selectedFiat) == 0 {
		return nil, nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "at least one fiat currency must be selected")
	}

	return selectedCrypto, selectedFiat, nil
}

func (s *CurrencyConfigService) SerializeSelection(crypto, fiat []string) (*string, *string) {
	cryptoCSV := currency.JoinList(crypto)
	fiatCSV := currency.JoinList(fiat)
	return &cryptoCSV, &fiatCSV
}
