package handler

import (
	"github.com/gin-gonic/gin"
	"motewallet-withdrawal/backend/internal/pkg/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *gin.Context) {
	response.Success(c, gin.H{
		"status": "ok",
	})
}
