package dto

import (
	"encoding/json"
	"testing"
)

func TestExchangeOrderQueryRespUnmarshalsFailureReasons(t *testing.T) {
	var resp ExchangeOrderQueryResp
	err := json.Unmarshal([]byte(`{
		"orderId": "order-1",
		"orderStatus": "FAIL",
		"rejectReason": "pair is unavailable",
		"failReason": "risk control rejected"
	}`), &resp)
	if err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.RejectReason != "pair is unavailable" {
		t.Fatalf("unexpected reject reason: %q", resp.RejectReason)
	}
	if resp.FailReason != "risk control rejected" {
		t.Fatalf("unexpected fail reason: %q", resp.FailReason)
	}
}

func TestInnerMatchQueryRespUnmarshalsFailureReasons(t *testing.T) {
	var resp InnerMatchQueryResp
	err := json.Unmarshal([]byte(`{
		"orderId": "order-2",
		"orderStatus": "FAIL",
		"rejectReason": "insufficient balance",
		"failReason": "order rejected"
	}`), &resp)
	if err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.RejectReason != "insufficient balance" {
		t.Fatalf("unexpected reject reason: %q", resp.RejectReason)
	}
	if resp.FailReason != "order rejected" {
		t.Fatalf("unexpected fail reason: %q", resp.FailReason)
	}
}
