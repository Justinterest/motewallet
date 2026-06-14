package dto

import (
	"encoding/json"
	"testing"
)

func TestFlexInt_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{"number", `1`, 1, false},
		{"string", `"2"`, 2, false},
		{"empty string", `""`, 0, false},
		{"null", `null`, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got FlexInt
			err := got.UnmarshalJSON([]byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Int() != tt.want {
				t.Fatalf("got %d want %d", got.Int(), tt.want)
			}
		})
	}
}

func TestDepositHistoryQueryResp_UnmarshalJSON_stringPagination(t *testing.T) {
	payload := []byte(`{
		"pageNo": "1",
		"pageSize": "10",
		"totalSize": "3",
		"totalPage": "1",
		"rows": [{
			"orderId": "1302347500618572718081",
			"orderCurrency": "USDT",
			"orderAmount": "1.5",
			"chain": "ETH_ERC20",
			"orderStatus": "SUCCESS",
			"orderTime": "2025-11-18 10:07:41"
		}]
	}`)

	var resp DepositHistoryQueryResp
	if err := json.Unmarshal(payload, &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if resp.PageNo.Int() != 1 {
		t.Fatalf("pageNo=%d", resp.PageNo.Int())
	}
	if resp.Total() != 3 {
		t.Fatalf("total=%d", resp.Total())
	}
	if len(resp.Items()) != 1 {
		t.Fatalf("items=%d", len(resp.Items()))
	}
}
