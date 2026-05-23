package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	bizerrors "motewallet-withdrawal/backend/internal/pkg/errors"
	"motewallet-withdrawal/backend/internal/pkg/jwt"
	"motewallet-withdrawal/backend/internal/pkg/response"
)

// AdminAuth extracts the "admin_token" cookie, validates the JWT,
// checks that UserType is ADMIN, and sets admin_id, admin_email, admin_role in the context.
func AdminAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(jwt.CookieAdminToken)
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

		if claims.UserType != "ADMIN" {
			response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		c.Set("admin_id", claims.UserID)
		c.Set("admin_email", claims.Email)
		// Role is not in the JWT claims — we'd need to look it up if needed.
		// For now, handlers that need the role can query the service.
		c.Next()
	}
}
