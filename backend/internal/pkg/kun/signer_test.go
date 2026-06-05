package kun

import (
	"strings"
	"testing"
)

func TestBuildCanonicalKVString_docExample(t *testing.T) {
	got := buildCanonicalKVString(
		map[string]interface{}{
			"orderId": "123456789",
			"amount":  "100.50",
		},
		map[string]string{
			"customerNo": "10037227",
			"timestamp":  "1758012303512",
		},
	)

	want := "amount=100.50&customerNo=10037227&orderid=123456789&timestamp=1758012303512"
	if got != want {
		t.Fatalf("canonical string mismatch\nwant: %s\ngot:  %s", want, got)
	}
}

func TestBuildCanonicalKVString_customerNoNotLowercased(t *testing.T) {
	got := buildCanonicalKVString(nil, map[string]string{
		"customerNo": "SUB001",
		"timestamp":  "1758012303512",
	})

	if !strings.Contains(got, "customerNo=SUB001") {
		t.Fatalf("expected customerNo with camelCase, got: %s", got)
	}
	if strings.Contains(got, "customerno=") {
		t.Fatalf("customerNo must not be lowercased to customerno, got: %s", got)
	}
}
