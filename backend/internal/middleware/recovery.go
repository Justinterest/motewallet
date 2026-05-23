package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	bizerrors "motewallet-withdrawal/backend/internal/pkg/errors"
	"motewallet-withdrawal/backend/internal/pkg/response"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic recovered",
					slog.Any("error", r),
					slog.String("stack", string(debug.Stack())),
				)
				response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
