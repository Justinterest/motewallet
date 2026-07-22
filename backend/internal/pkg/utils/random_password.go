package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const passwordCharset = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// GenerateRandomPassword returns a cryptographically random password of the given length.
func GenerateRandomPassword(length int) (string, error) {
	if length < 8 {
		length = 12
	}
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(passwordCharset)))
	for i := range result {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", fmt.Errorf("generate random password: %w", err)
		}
		result[i] = passwordCharset[n.Int64()]
	}
	return string(result), nil
}
