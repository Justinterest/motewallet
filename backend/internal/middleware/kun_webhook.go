package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"motewallet-withdrawal/backend/internal/config"
	"motewallet-withdrawal/backend/internal/pkg/kun"
	kundto "motewallet-withdrawal/backend/internal/pkg/kun/dto"
	"motewallet-withdrawal/backend/internal/pkg/response"
)

func KUNWebhookAuth(kunCfg config.KUNConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			response.Error(c, http.StatusBadRequest, 40000, "failed to read request body")
			c.Abort()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		var event kundto.WebhookEvent
		if err := json.Unmarshal(bodyBytes, &event); err != nil {
			response.Error(c, http.StatusBadRequest, 40001, "invalid webhook payload")
			c.Abort()
			return
		}

		if event.Timestamp != "" {
			ts, err := strconv.ParseInt(event.Timestamp, 10, 64)
			if err == nil {
				diff := time.Since(time.UnixMilli(ts))
				if math.Abs(diff.Seconds()) > kunCfg.WebhookTimeDiff.Seconds() {
					response.Error(c, http.StatusUnauthorized, 40002, "webhook timestamp expired")
					c.Abort()
					return
				}
			}
		}

		if kunCfg.PublicKey != "" && !kunCfg.MockEnabled {
			if err := kun.VerifyWebhookSignature(kunCfg.PublicKey, event.Data, event.Timestamp, event.Sign); err != nil {
				response.Error(c, http.StatusUnauthorized, 40003, "webhook signature verification failed")
				c.Abort()
				return
			}
		}

		c.Set("webhook_event", &event)
		c.Next()
	}
}
