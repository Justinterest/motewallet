package kun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

var kunLogRequestHeaderKeys = []string{
	"Content-Type",
	"App-Key",
	"Api-Version",
	"Customer-No",
	"Timestamp",
	"Sign",
}

var kunLogResponseHeaderKeys = []string{
	"Content-Type",
	"Sign",
	"sign",
	"Timestamp",
	"timestamp",
}

type kunAPILog struct {
	Method          string            `json:"method"`
	URL             string            `json:"url"`
	Path            string            `json:"path"`
	SignString      string            `json:"sign_string,omitempty"`
	RequestHeaders  map[string]string `json:"request_headers"`
	RequestBody     any               `json:"request_body,omitempty"`
	HTTPStatus      int               `json:"http_status"`
	ResponseHeaders map[string]string `json:"response_headers"`
	ResponseBody    any               `json:"response_body,omitempty"`
}

func kunHeadersFrom(h http.Header, keys []string) map[string]string {
	out := make(map[string]string, len(keys))
	for _, key := range keys {
		if v := h.Get(key); v != "" {
			out[key] = v
		}
	}
	return out
}

func kunLogJSONValue(v any) any {
	switch t := v.(type) {
	case nil:
		return nil
	case json.RawMessage:
		if len(t) == 0 {
			return nil
		}
		var parsed any
		if json.Unmarshal(t, &parsed) == nil {
			return parsed
		}
		return string(t)
	case []byte:
		if len(t) == 0 {
			return nil
		}
		var parsed any
		if json.Unmarshal(t, &parsed) == nil {
			return parsed
		}
		return string(t)
	case string:
		if t == "" {
			return nil
		}
		var parsed any
		if json.Unmarshal([]byte(t), &parsed) == nil {
			return parsed
		}
		return t
	default:
		return v
	}
}

func logKUNRequest(
	method, url, path, signString string,
	reqHeaders, respHeaders http.Header,
	requestBody any,
	httpStatus int,
	responseBody []byte,
) {
	entry := kunAPILog{
		Method:         method,
		URL:            url,
		Path:           path,
		SignString:     signString,
		RequestHeaders: kunHeadersFrom(reqHeaders, kunLogRequestHeaderKeys),
		// RequestBody:     kunLogJSONValue(requestBody),
		HTTPStatus:      httpStatus,
		ResponseHeaders: kunHeadersFrom(respHeaders, kunLogResponseHeaderKeys),
		ResponseBody:    kunLogJSONValue(responseBody),
	}

	payload, err := marshalKUNLogJSON(entry)
	if err != nil {
		slog.Error("kun api call log marshal failed", slog.Any("error", err))
		return
	}

	// Write raw JSON to stdout so slog does not escape quotes/newlines (easy copy-paste).
	fmt.Fprintf(os.Stdout, "[kun api call]\n%s\n\n", payload)
}

func marshalKUNLogJSON(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

func kunMockResponseEnvelope(data json.RawMessage) []byte {
	envelope := map[string]any{
		"code": "000000",
		"msg":  "mock",
		"data": data,
	}
	b, err := json.Marshal(envelope)
	if err != nil {
		return []byte(`{"code":"000000","msg":"mock","data":{}}`)
	}
	return b
}

func logKUNMockCall(path, customerNo string, requestBody any, data json.RawMessage) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Customer-No", customerNo)
	headers.Set("X-Mock", "true")

	logKUNRequest(http.MethodPost, "mock://kun"+path, path, "", headers, nil, requestBody, http.StatusOK, kunMockResponseEnvelope(data))
}
