package kun

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
// - business keys are trimmed and lower-cased (e.g. orderId -> orderid)
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
		// requestNo 不需要转成小写
		key := strings.TrimSpace(k)
		if strings.EqualFold(key, "requestNo") {
			key = key
		} else {
			key = strings.ToLower(key)
		}
		if key == "" || key == "sign" {
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
		// For objects/arrays/numbers we fall back to JSON representation to avoid Go's "%v" map formatting.
		b, err := json.Marshal(t)
		if err != nil {
			return fmt.Sprintf("%v", t)
		}
		return string(b)
	}
}

// StructToMap flattens a struct to map[string]interface{} using JSON marshal/unmarshal.
func StructToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal struct: %w", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("unmarshal to map: %w", err)
	}
	return m, nil
}
