package jwt

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

const (
	CookieMerchantToken = "token"
	CookieAdminToken    = "admin_token"

	PurposeSession        = ""
	Purpose2FAVerify      = "2fa_verify"
	Purpose2FASetup       = "2fa_setup"
	Purpose2FARebind      = "2fa_rebind"
	PurposePasswordChange = "password_change"
)

type Claims struct {
	UserID   uint64 `json:"user_id"`
	UserType string `json:"user_type"` // MERCHANT or ADMIN
	Email    string `json:"email"`
	Purpose  string `json:"purpose,omitempty"`
	jwtlib.RegisteredClaims
}

// GenerateToken creates a signed JWT session token.
func GenerateToken(userID uint64, userType, email, secret string, expiry time.Duration) (string, error) {
	return GenerateTokenWithPurpose(userID, userType, email, PurposeSession, secret, expiry)
}

// GenerateTokenWithPurpose creates a signed JWT with an optional purpose claim.
func GenerateTokenWithPurpose(userID uint64, userType, email, purpose, secret string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		UserType: userType,
		Email:    email,
		Purpose:  purpose,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwtlib.NewNumericDate(now),
			NotBefore: jwtlib.NewNumericDate(now),
			Issuer:    "motewallet",
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken validates a JWT token string and returns the claims.
func ParseToken(tokenString, secret string) (*Claims, error) {
	token, err := jwtlib.ParseWithClaims(tokenString, &Claims{}, func(token *jwtlib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// SetTokenCookie sets an httpOnly cookie on the response.
func SetTokenCookie(c *gin.Context, name, token string, maxAge int, secure bool) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(name, token, maxAge, "/", "", secure, true)
}

// ClearTokenCookie removes a cookie by setting maxAge to -1.
func ClearTokenCookie(c *gin.Context, name string, secure bool) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(name, "", -1, "/", "", secure, true)
}
