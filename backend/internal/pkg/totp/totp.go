package totp

import (
	"fmt"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const issuer = "Motewallet"

// GenerateSecret creates a new TOTP secret for the given account email.
func GenerateSecret(email string) (secret string, uri string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: email,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", "", fmt.Errorf("generate totp secret: %w", err)
	}
	return key.Secret(), key.URL(), nil
}

// Validate checks a 6-digit TOTP code against the secret.
func Validate(code, secret string) bool {
	if secret == "" || code == "" {
		return false
	}
	return totp.Validate(code, secret)
}
