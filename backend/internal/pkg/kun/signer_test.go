package kun

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMarshalRequestBody_noHTMLEscape(t *testing.T) {
	type payload struct {
		EmployeeNum string `json:"employeeNum"`
	}
	got, err := MarshalRequestBody(payload{EmployeeNum: "<10"})
	if err != nil {
		t.Fatal(err)
	}
	body := string(got)
	if strings.Contains(body, `\u003c`) {
		t.Fatalf("KUN JSON must not HTML-escape '<': %s", body)
	}
	if !strings.Contains(body, `"employeeNum":"<10"`) {
		t.Fatalf("unexpected body: %s", body)
	}
}

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

	want := "amount=100.50&customerNo=10037227&orderId=123456789&timestamp=1758012303512"
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

func TestMarshalNestedJSON_preservesMapKeyOrder(t *testing.T) {
	nested := map[string]interface{}{
		"zLast":  "3",
		"aFirst": "1",
		"mMid":   "2",
	}
	got, err := marshalNestedJSON(nested)
	if err != nil {
		t.Fatalf("marshalNestedJSON: %v", err)
	}

	sorted, _ := json.Marshal(nested)
	if got == string(sorted) {
		t.Fatalf("nested JSON must not use encoding/json alphabetical key order\ngot:  %s\nsort: %s", got, sorted)
	}

	for _, key := range []string{"zLast", "aFirst", "mMid"} {
		if !strings.Contains(got, `"`+key+`"`) {
			t.Fatalf("expected key %q in %s", key, got)
		}
	}
}

func TestStructToMap_nestedPreservesBodyKeyOrder(t *testing.T) {
	type inner struct {
		ZField string `json:"zField"`
		AField string `json:"aField"`
	}
	type outer struct {
		Nested inner `json:"nested"`
	}
	body, err := json.Marshal(outer{Nested: inner{ZField: "z", AField: "a"}})
	if err != nil {
		t.Fatal(err)
	}

	m, err := StructToMap(outer{Nested: inner{ZField: "z", AField: "a"}})
	if err != nil {
		t.Fatal(err)
	}

	nestedStr, ok := m["nested"].(string)
	if !ok {
		t.Fatalf("expected nested value as JSON string, got %T %#v", m["nested"], m["nested"])
	}
	if !strings.Contains(nestedStr, `"zField"`) || !strings.HasPrefix(nestedStr, `{"zField"`) {
		t.Fatalf("nested JSON should preserve struct field order from body\nbody:   %s\nnested: %s", body, nestedStr)
	}
}

func TestMarshalNestedJSON_deepNestingPreservesMapKeyOrder(t *testing.T) {
	nested := map[string]interface{}{
		"zOuter": map[string]interface{}{
			"zInner": "last",
			"aInner": "first",
		},
		"aOuter": "top",
	}
	got, err := marshalNestedJSON(nested)
	if err != nil {
		t.Fatalf("marshalNestedJSON: %v", err)
	}

	sorted, _ := json.Marshal(nested)
	if got == string(sorted) {
		t.Fatalf("deep nested JSON must not use encoding/json alphabetical key order\ngot:  %s\nsort: %s", got, sorted)
	}
}

func TestMarshalNestedJSON_mapStringString(t *testing.T) {
	got, err := marshalNestedJSON(map[string]string{"zKey": "1", "aKey": "2"})
	if err != nil {
		t.Fatalf("marshalNestedJSON: %v", err)
	}
	sorted, _ := json.Marshal(map[string]string{"zKey": "1", "aKey": "2"})
	if got == string(sorted) {
		t.Fatalf("map[string]string must not use encoding/json alphabetical key order\ngot:  %s\nsort: %s", got, sorted)
	}
}
