package kun

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// CanonicalizeAndSignRequest builds the KUN signature string and signs it using SHA256withRSA.
//
// Doc rules (high level):
// - Business params: top-level fields of query/form/json body participate in signature.
// - System params: Customer-No and Timestamp headers MUST participate in signature as:
//   Customer-No -> customerNo, Timestamp -> timestamp
// - Do NOT include Sign itself in the signature set.
// - Sort by ASCII (String.compareTo) and concatenate key=value with '&'
// - Values are percent-encoded (URL encode) consistently
// - Signature output is Base64-encoded SHA256withRSA.
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
// - keys are trimmed and lower-cased
// - "sign" key is ignored
// - nil values become empty string
// - values are percent-encoded
// - keys are ASCII sorted and joined as key=value&key2=value2
//
// extraPlain is used for injecting system params (already in the correct key casing),
// such as customerNo/timestamp.
func buildCanonicalKVString(
	biz map[string]interface{},
	extraPlain map[string]string,
) string {
	params := make(map[string]string, len(biz)+len(extraPlain))

	for k, v := range biz {
		key := strings.ToLower(strings.TrimSpace(k))
		if key == "" || key == "sign" {
			continue
		}
		params[key] = percentEncode(valueToString(v))
	}

	for k, v := range extraPlain {
		key := strings.ToLower(strings.TrimSpace(k))
		if key == "" || key == "sign" {
			continue
		}
		params[key] = percentEncode(v)
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

// percentEncode applies an RFC3986-like encoding suitable for signature canonicalization.
// It starts from QueryEscape but normalizes space to %20 (not '+').
func percentEncode(s string) string {
	enc := url.QueryEscape(s)
	enc = strings.ReplaceAll(enc, "+", "%20")
	// Keep "~" unescaped for RFC3986 compatibility.
	enc = strings.ReplaceAll(enc, "%7E", "~")
	return enc
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
