package service

import (
	"context"
	"testing"

	"motewallet/internal/model"
	"motewallet/internal/pkg/currency"
)

type currencyConfigRepoStub struct {
	values map[string]string
}

func (r *currencyConfigRepoStub) GetByKey(_ context.Context, key string) (string, error) {
	return r.values[key], nil
}

func (r *currencyConfigRepoStub) Upsert(_ context.Context, key, value, _ string) error {
	r.values[key] = value
	return nil
}

func newCurrencyConfigServiceForTest() *CurrencyConfigService {
	return NewCurrencyConfigService(&currencyConfigRepoStub{values: map[string]string{
		currency.ConfigKeyCrypto:        "USDT",
		currency.ConfigKeyFiat:          "USD",
		currency.ConfigKeyCryptoChains:  `{"USDT":["TRX_TRC20"]}`,
		currency.ConfigKeyDefaultChains: `{"USDT":"TRX_TRC20"}`,
	}})
}

func TestMerchantCurrenciesUseSystemConfigOnlyAsDefault(t *testing.T) {
	svc := newCurrencyConfigServiceForTest()
	ctx := context.Background()

	defaultCrypto, err := svc.GetSupportedCrypto(ctx, &model.Merchant{})
	if err != nil {
		t.Fatal(err)
	}
	assertStringSlice(t, defaultCrypto, []string{"USDT"})

	explicitCrypto := "BTC,USDC"
	explicitFiat := "EUR,HKD"
	merchant := &model.Merchant{
		SupportedCryptoCurrencies: &explicitCrypto,
		SupportedFiatCurrencies:   &explicitFiat,
	}

	gotCrypto, err := svc.GetSupportedCrypto(ctx, merchant)
	if err != nil {
		t.Fatal(err)
	}
	assertStringSlice(t, gotCrypto, []string{"BTC", "USDC"})

	gotFiat, err := svc.GetSupportedFiat(ctx, merchant)
	if err != nil {
		t.Fatal(err)
	}
	assertStringSlice(t, gotFiat, []string{"EUR", "HKD"})
}

func TestNormalizeMerchantSelectionUsesFullCatalog(t *testing.T) {
	svc := newCurrencyConfigServiceForTest()

	cryptoSelection, fiatSelection, err := svc.NormalizeMerchantSelection(
		context.Background(),
		[]string{"BTC", "USDC"},
		[]string{"EUR"},
	)
	if err != nil {
		t.Fatal(err)
	}
	assertStringSlice(t, cryptoSelection, []string{"BTC", "USDC"})
	assertStringSlice(t, fiatSelection, []string{"EUR"})
}

func TestMerchantChainsUseSystemChainConfig(t *testing.T) {
	svc := newCurrencyConfigServiceForTest()
	ctx := context.Background()

	selected, defaults, err := svc.NormalizeMerchantChainSelection(
		ctx,
		[]string{"USDT", "BTC"},
		map[string][]string{
			"USDT": {"TRX_TRC20"},
			"BTC":  {"BTC"},
		},
		map[string]string{
			"USDT": "TRX_TRC20",
			"BTC":  "BTC",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	assertStringSlice(t, selected["USDT"], []string{"TRX_TRC20"})
	assertStringSlice(t, selected["BTC"], []string{"BTC"})
	if defaults["USDT"] != "TRX_TRC20" || defaults["BTC"] != "BTC" {
		t.Fatalf("unexpected defaults: %#v", defaults)
	}

	explicitCrypto := "USDT"
	explicitChains := `{"USDT":["TRX_TRC20"]}`
	merchant := &model.Merchant{
		SupportedCryptoCurrencies: &explicitCrypto,
		SupportedCryptoChains:     &explicitChains,
	}
	got, err := svc.GetSupportedChains(ctx, merchant)
	if err != nil {
		t.Fatal(err)
	}
	assertStringSlice(t, got["USDT"], []string{"TRX_TRC20"})
}

func TestSystemChainsAreConfiguredForCurrenciesThatAreDefaultDisabled(t *testing.T) {
	svc := newCurrencyConfigServiceForTest()

	chains, err := svc.GetAvailableChains(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	assertStringSlice(t, chains["USDT"], []string{"TRX_TRC20"})
	assertStringSlice(t, chains["USDC"], currency.CatalogChains("USDC"))
	assertStringSlice(t, chains["BTC"], []string{"BTC"})
}

func TestNormalizeGlobalSelectionRequiresChainsForEveryCryptoCurrency(t *testing.T) {
	svc := newCurrencyConfigServiceForTest()

	_, _, chains, defaults, err := svc.NormalizeGlobalSelection(
		context.Background(),
		[]string{"USDT"},
		[]string{"USD"},
		map[string][]string{
			"USDT": {"TRX_TRC20"},
			"USDC": {"ETH_ERC20"},
			"BTC":  {"BTC"},
		},
		map[string]string{
			"USDT": "TRX_TRC20",
			"USDC": "ETH_ERC20",
			"BTC":  "BTC",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	assertStringSlice(t, chains["USDT"], []string{"TRX_TRC20"})
	assertStringSlice(t, chains["USDC"], []string{"ETH_ERC20"})
	assertStringSlice(t, chains["BTC"], []string{"BTC"})
	if defaults["USDC"] != "ETH_ERC20" || defaults["BTC"] != "BTC" {
		t.Fatalf("unexpected defaults: %#v", defaults)
	}
}

func assertStringSlice(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	}
}
