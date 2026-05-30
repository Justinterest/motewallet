package kun

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"
)

func TestParseRSAPrivateKey_rawBase64(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}

	raw := base64.StdEncoding.EncodeToString(der)
	if _, err := ParseRSAPrivateKey(raw); err != nil {
		t.Fatalf("parse raw base64 private key: %v", err)
	}
}

func TestParseRSAPrivateKey_pem(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	der := x509.MarshalPKCS1PrivateKey(key)
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: der})

	if _, err := ParseRSAPrivateKey(string(pemBytes)); err != nil {
		t.Fatalf("parse PEM private key: %v", err)
	}
}

func TestParseRSAPublicKey_rawBase64(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	raw := base64.StdEncoding.EncodeToString(der)
	if _, err := ParseRSAPublicKey(raw); err != nil {
		t.Fatalf("parse raw base64 public key: %v", err)
	}
}
