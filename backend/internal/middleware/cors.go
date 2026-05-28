package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"motewallet/internal/config"
)

func CORS(cfg *config.Config) gin.HandlerFunc {
	allowedOrigins := []string{cfg.Frontend.URL, cfg.Frontend.AdminURL}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		allowed := false
		for _, o := range allowedOrigins {
			if strings.TrimRight(o, "/") == strings.TrimRight(origin, "/") {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
