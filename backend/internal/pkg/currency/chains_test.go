package currency

import "testing"

func TestParseChainMapAndDefaults(t *testing.T) {
	chains := ParseChainMap(`{"USDT":["ETH_ERC20","TRX_TRC20"],"BTC":["BTC"]}`)
	if len(chains["USDT"]) != 2 {
		t.Fatalf("expected 2 USDT chains, got %v", chains["USDT"])
	}
	defaults := ParseDefaultChainMap(`{"USDT":"TRX_TRC20","BTC":"BTC"}`)
	if defaults["USDT"] != "TRX_TRC20" {
		t.Fatalf("unexpected default: %s", defaults["USDT"])
	}
	resolved := ResolveDefaultChain("USDT", chains["USDT"], "ETH_ERC20")
	if resolved != "ETH_ERC20" {
		t.Fatalf("expected ETH_ERC20, got %s", resolved)
	}
	resolved = ResolveDefaultChain("USDT", chains["USDT"], "TON")
	if resolved != "ETH_ERC20" {
		t.Fatalf("expected first supported chain fallback, got %s", resolved)
	}
}

func TestValidateChainSelection(t *testing.T) {
	selected := ValidateChainSelection("USDT", []string{"ETH_ERC20", "UNKNOWN", "TRX_TRC20"}, CatalogChains("USDT"))
	if len(selected) != 2 {
		t.Fatalf("expected 2 valid chains, got %v", selected)
	}
}
