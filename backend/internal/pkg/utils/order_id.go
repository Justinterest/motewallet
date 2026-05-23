package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func GeneratePlatformOrderID(prefix string) string {
	ts := time.Now().Format("20060102150405")
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return "MW" + prefix + ts + hex.EncodeToString(b)
}
