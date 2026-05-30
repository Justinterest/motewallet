package kun

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
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

// ParseRSAPublicKey parses a PEM-encoded RSA public key.
func ParseRSAPublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaPub, nil
}
