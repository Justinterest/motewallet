package dto

import (
	"encoding/json"
	"testing"
)

func TestAgreementListUnmarshalArray(t *testing.T) {
	raw := `[{"protocolId":"p1","title":"Service Agreement","url":"https://example.com/a.pdf","signStatus":"UNSIGN","version":"1.0"}]`
	var list AgreementList
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		t.Fatalf("unmarshal array: %v", err)
	}
	if len(list) != 1 || list[0].ProtocolId != "p1" {
		t.Fatalf("unexpected list: %+v", list)
	}
}

func TestAgreementListUnmarshalWrapped(t *testing.T) {
	raw := `{"list":[{"protocolId":"p2","title":"Privacy Policy","url":"https://example.com/b.pdf","signStatus":"UNSIGN"}]}`
	var list AgreementList
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		t.Fatalf("unmarshal wrapped: %v", err)
	}
	if len(list) != 1 || list[0].ProtocolId != "p2" {
		t.Fatalf("unexpected list: %+v", list)
	}
}
