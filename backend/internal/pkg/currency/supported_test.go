package currency

import "testing"

func TestParseList(t *testing.T) {
	items := ParseList("USDT, USDC, BTC,USDT")
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
	if items[0] != "USDT" || items[1] != "USDC" || items[2] != "BTC" {
		t.Fatalf("unexpected items: %#v", items)
	}
}

func TestFilterAllowed(t *testing.T) {
	result := FilterAllowed([]string{"USDT", "ETH"}, []string{"USDT", "USDC", "BTC"})
	if len(result) != 1 || result[0] != "USDT" {
		t.Fatalf("unexpected result: %#v", result)
	}
}
