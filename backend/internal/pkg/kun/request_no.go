package kun

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// GenerateRequestNo produces a unique request number for KUN API idempotency.
// Format: "MW" + YYYYMMDDHHmmss + 8-char random hex
func GenerateRequestNo() string {
	ts := time.Now().Format("20060102150405")
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return "MW" + ts + hex.EncodeToString(b)
}
