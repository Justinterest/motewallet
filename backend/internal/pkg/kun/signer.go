package kun

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// CanonicalizeAndSignRequest builds the KUN signature string and signs it using SHA256withRSA.
//
// Doc rules (high level):
//   - Business params: top-level fields of query/form/json body participate in signature.
//   - System params: Customer-No and Timestamp headers MUST participate in signature as:
//     Customer-No -> customerNo, Timestamp -> timestamp
//   - Do NOT include Sign itself in the signature set.
//   - Sort by ASCII (String.compareTo) and concatenate key=value with '&'
//   - Values are used as-is when building key=value (same as KUN Java getSignRowData sample)
//   - Signature output is Base64-encoded SHA256withRSA (identical to Java SHA256withRSA).
func CanonicalizeAndSignRequest(
	bizParams map[string]interface{},
	customerNo string,
	timestamp string,
	privateKeyPEM string,
) (signString string, signatureBase64 string, err error) {
	signString = buildCanonicalKVString(bizParams, map[string]string{
		"customerNo": customerNo,
		"timestamp":  timestamp,
	})

	priv, err := ParseRSAPrivateKey(privateKeyPEM)
	if err != nil {
		return "", "", err
	}

	sum := sha256.Sum256([]byte(signString))
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, sum[:])
	if err != nil {
		return "", "", fmt.Errorf("RSA sign failed: %w", err)
	}

	return signString, base64.StdEncoding.EncodeToString(sig), nil
}

// buildCanonicalKVString applies KUN canonicalization rules:
// - business keys are trimmed and keep JSON body casing (e.g. requestNo, directorInfo)
// - system keys in extraPlain keep doc-mapped casing (customerNo, timestamp)
// - "sign" key is ignored
// - null values are skipped (KUN Java sample skips entry.getValue() == null)
// - keys are ASCII sorted and joined as key=value&key2=value2
func buildCanonicalKVString(
	biz map[string]interface{},
	extraPlain map[string]string,
) string {
	params := make(map[string]string, len(biz)+len(extraPlain))

	for k, v := range biz {
		if v == nil {
			continue
		}
		key := strings.TrimSpace(k)
		if key == "" || strings.EqualFold(key, "sign") {
			continue
		}
		params[key] = valueToString(v)
	}

	for k, v := range extraPlain {
		key := strings.TrimSpace(k)
		if key == "" || strings.EqualFold(key, "sign") {
			continue
		}
		params[key] = v
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, k+"="+params[k])
	}
	return strings.Join(pairs, "&")
}

func valueToString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case json.Number:
		return t.String()
	case float64:
		// StructToMap uses encoding/json which decodes numbers into float64 by default.
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		s, err := marshalNestedJSON(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return s
	}
}

// marshalNestedJSON serializes nested objects/arrays without sorting map keys.
// It recurses through every map/slice level. encoding/json.Marshal is not used for maps
// because it sorts keys alphabetically; KUN nested values follow hash-map iteration order.
func marshalNestedJSON(v interface{}) (string, error) {
	if v == nil {
		return "null", nil
	}

	switch t := v.(type) {
	case string:
		return marshalJSONString(t)
	case json.Number:
		return t.String(), nil
	case json.RawMessage:
		return compactJSONBytes(t), nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	case bool:
		if t {
			return "true", nil
		}
		return "false", nil
	case map[string]interface{}:
		return marshalJSONObject(t)
	case []interface{}:
		return marshalJSONArray(t)
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		return marshalReflectMap(rv)
	case reflect.Slice, reflect.Array:
		return marshalReflectSlice(rv)
	default:
		b, err := MarshalRequestBody(v)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
}

func marshalJSONObject(m map[string]interface{}) (string, error) {
	if len(m) == 0 {
		return "{}", nil
	}
	var buf strings.Builder
	buf.WriteByte('{')
	first := true
	for k, val := range m {
		if !first {
			buf.WriteByte(',')
		}
		first = false
		kb, err := marshalJSONString(k)
		if err != nil {
			return "", err
		}
		buf.WriteString(kb)
		buf.WriteByte(':')
		part, err := marshalNestedJSON(val)
		if err != nil {
			return "", err
		}
		buf.WriteString(part)
	}
	buf.WriteByte('}')
	return buf.String(), nil
}

func marshalJSONArray(items []interface{}) (string, error) {
	if len(items) == 0 {
		return "[]", nil
	}
	var buf strings.Builder
	buf.WriteByte('[')
	for i, val := range items {
		if i > 0 {
			buf.WriteByte(',')
		}
		part, err := marshalNestedJSON(val)
		if err != nil {
			return "", err
		}
		buf.WriteString(part)
	}
	buf.WriteByte(']')
	return buf.String(), nil
}

func marshalReflectMap(rv reflect.Value) (string, error) {
	if rv.IsNil() {
		return "null", nil
	}
	var buf strings.Builder
	buf.WriteByte('{')
	first := true
	iter := rv.MapRange()
	for iter.Next() {
		if !first {
			buf.WriteByte(',')
		}
		first = false
		kb, err := marshalJSONString(fmt.Sprint(iter.Key().Interface()))
		if err != nil {
			return "", err
		}
		buf.WriteString(kb)
		buf.WriteByte(':')
		part, err := marshalNestedJSON(iter.Value().Interface())
		if err != nil {
			return "", err
		}
		buf.WriteString(part)
	}
	buf.WriteByte('}')
	return buf.String(), nil
}

func marshalReflectSlice(rv reflect.Value) (string, error) {
	if rv.Kind() == reflect.Slice && rv.IsNil() {
		return "null", nil
	}
	length := rv.Len()
	if length == 0 {
		return "[]", nil
	}
	items := make([]interface{}, length)
	for i := 0; i < length; i++ {
		items[i] = rv.Index(i).Interface()
	}
	return marshalJSONArray(items)
}

func marshalJSONString(s string) (string, error) {
	return strconv.Quote(s), nil
}

// MarshalRequestBody serializes KUN API request bodies without HTML escaping.
// KUN expects "<" to remain literal in JSON (not "\u003c") for both HTTP body and signing.
func MarshalRequestBody(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimSpace(buf.Bytes()), nil
}

// StructToMap flattens a struct to map[string]interface{} for signing.
// Uses MarshalRequestBody so nested JSON substrings match the HTTP request body exactly.
func StructToMap(v interface{}) (map[string]interface{}, error) {
	data, err := MarshalRequestBody(v)
	if err != nil {
		return nil, fmt.Errorf("marshal struct: %w", err)
	}
	return mapFromJSON(data)
}

func mapFromJSON(data []byte) (map[string]interface{}, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("unmarshal to map: %w", err)
	}
	m := make(map[string]interface{}, len(raw))
	for k, v := range raw {
		parsed, err := parseJSONValue(v)
		if err != nil {
			return nil, fmt.Errorf("parse field %q: %w", k, err)
		}
		m[k] = parsed
	}
	return m, nil
}

func parseJSONValue(raw json.RawMessage) (interface{}, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	switch raw[0] {
	case '"':
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return nil, err
		}
		return s, nil
	case '{', '[':
		return compactJSONBytes(raw), nil
	case 't', 'f':
		var b bool
		if err := json.Unmarshal(raw, &b); err != nil {
			return nil, err
		}
		return b, nil
	default:
		var n json.Number
		if err := json.Unmarshal(raw, &n); err != nil {
			return nil, err
		}
		return n, nil
	}
}

func compactJSONBytes(raw json.RawMessage) string {
	var buf bytes.Buffer
	if err := json.Compact(&buf, raw); err != nil {
		return string(raw)
	}
	return buf.String()
}
