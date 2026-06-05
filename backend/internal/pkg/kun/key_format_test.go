package kun

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseBothKeyFormats(t *testing.T) {
	dir := filepath.Join("..", "..", "..", ".kun-keys")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skip("no .kun-keys directory")
	}

	pkcs8PEM, err := os.ReadFile(filepath.Join(dir, "app_private_key_pkcs8.pem"))
	if err != nil {
		t.Fatal(err)
	}
	pkcs8B64, err := os.ReadFile(filepath.Join(dir, "app_private_key_pkcs8.b64.txt"))
	if err != nil {
		t.Fatal(err)
	}
	pkcs1PEM, err := os.ReadFile(filepath.Join(dir, "app_private_key.pem"))
	if err != nil {
		t.Fatal(err)
	}

	for name, key := range map[string]string{
		"pkcs8_pem": string(pkcs8PEM),
		"pkcs8_b64": string(pkcs8B64),
		"pkcs1_pem": string(pkcs1PEM),
	} {
		if _, err := ParseRSAPrivateKey(key); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
}
