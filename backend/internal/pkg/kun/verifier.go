package kun

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// VerifyWebhookSignature verifies incoming webhook using SHA256withRSA.
func VerifyWebhookSignature(pubKeyPEM string, data map[string]interface{}, timestamp string, signature string) error {
	pubKey, err := ParseRSAPublicKey(pubKeyPEM)
	if err != nil {
		return fmt.Errorf("parse public key: %w", err)
	}

	payload := buildCanonicalKVString(data, map[string]string{
		"timestamp": timestamp,
	})

	hash := sha256.Sum256([]byte(payload))

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("base64 decode signature: %w", err)
	}

	if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], sigBytes); err != nil {
		return fmt.Errorf("RSA signature verification failed: %w", err)
	}

	return nil
}

// VerifyResponseSignature verifies the KUN response signature using SHA256withRSA.
// Per doc, signature is computed from response "data" only (plus timestamp).
func VerifyResponseSignature(pubKeyPEM string, data map[string]interface{}, timestamp string, signature string) error {
	return VerifyWebhookSignature(pubKeyPEM, data, timestamp, signature)
}

