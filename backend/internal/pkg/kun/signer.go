package kun

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// SignRequest produces HMAC-SHA256 signature for outgoing API requests.
// It merges params with the timestamp, sorts keys by ASCII, concatenates as
// key=value pairs joined by &, then HMAC-SHA256 with the secret.
func SignRequest(params map[string]interface{}, timestamp string, secret string) string {
	merged := make(map[string]interface{}, len(params)+1)
	for k, v := range params {
		merged[k] = v
	}
	merged["timestamp"] = timestamp

	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%v", k, merged[k]))
	}
	payload := strings.Join(pairs, "&")

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
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
