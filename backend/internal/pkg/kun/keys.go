package kun

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
)

// decodeKeyMaterial accepts PEM-wrapped keys or raw base64 DER (as provided by KUN dashboard).
func decodeKeyMaterial(keyStr string) ([]byte, error) {
	keyStr = strings.TrimSpace(keyStr)
	if keyStr == "" {
		return nil, fmt.Errorf("empty key")
	}

	// Allow literal \n in .env values.
	keyStr = strings.ReplaceAll(keyStr, `\n`, "\n")

	if strings.Contains(keyStr, "-----BEGIN") {
		block, _ := pem.Decode([]byte(keyStr))
		if block == nil {
			return nil, fmt.Errorf("failed to decode PEM block")
		}
		return block.Bytes, nil
	}

	der, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return nil, fmt.Errorf("decode base64 key: %w", err)
	}
	return der, nil
}

func ParseRSAPrivateKey(keyStr string) (*rsa.PrivateKey, error) {
	der, err := decodeKeyMaterial(keyStr)
	if err != nil {
		return nil, err
	}

	if parsed, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		key, ok := parsed.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA private key")
		}
		return key, nil
	}

	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}

	return nil, fmt.Errorf("parse private key: unsupported format")
}

func ParseRSAPublicKey(keyStr string) (*rsa.PublicKey, error) {
	der, err := decodeKeyMaterial(keyStr)
	if err != nil {
		return nil, err
	}

	if pub, err := x509.ParsePKIXPublicKey(der); err == nil {
		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA public key")
		}
		return rsaPub, nil
	}

	if pub, err := x509.ParsePKCS1PublicKey(der); err == nil {
		return pub, nil
	}

	return nil, fmt.Errorf("parse public key: unsupported format")
}
