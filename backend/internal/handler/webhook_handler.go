package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	kundto "motewallet/internal/pkg/kun/dto"
	"motewallet/internal/service"
)

type WebhookHandler struct {
	webhookService *service.WebhookService
}

func NewWebhookHandler(webhookService *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookService: webhookService,
	}
}

func (h *WebhookHandler) HandleKUN(c *gin.Context) {
	eventVal, exists := c.Get("webhook_event")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400"})
		return
	}

	event := eventVal.(*kundto.WebhookEvent)

	_ = h.webhookService.ProcessEvent(c.Request.Context(), event)

	c.JSON(http.StatusOK, gin.H{"code": "200"})
}
