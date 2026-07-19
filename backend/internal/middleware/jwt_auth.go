package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/jwt"
	"motewallet/internal/pkg/response"
)

// JWTAuth extracts the merchant "token" cookie, validates the JWT,
// checks that UserType is MERCHANT, and sets user_id + user_email in the context.
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(jwt.CookieMerchantToken)
		if err != nil || tokenString == "" {
			response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(tokenString, secret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, bizerrors.ErrInvalidToken, "invalid or expired token")
			c.Abort()
			return
		}

		if claims.UserType != "MERCHANT" {
			response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		// Challenge tokens (2FA verify/setup) must not access protected routes.
		if claims.Purpose != jwt.PurposeSession {
			response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Next()
	}
}
