package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		attrs := []slog.Attr{
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Duration("latency", latency),
			slog.String("client_ip", c.ClientIP()),
		}
		if query != "" {
			attrs = append(attrs, slog.String("query", query))
		}
		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()))
		}

		args := make([]any, len(attrs))
		for i, a := range attrs {
			args[i] = a
		}

		if status >= 500 {
			slog.Error("request completed", args...)
		} else if status >= 400 {
			slog.Warn("request completed", args...)
		} else {
			slog.Info("request completed", args...)
		}
	}
}
