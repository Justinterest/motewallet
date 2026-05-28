package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	dtoreq "motewallet/internal/dto/request"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/response"
	"motewallet/internal/service"
)

type DepositHandler struct {
	depositService *service.DepositService
}

func NewDepositHandler(depositService *service.DepositService) *DepositHandler {
	return &DepositHandler{
		depositService: depositService,
	}
}

func (h *DepositHandler) GetDepositAddresses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	var req dtoreq.GetDepositAddressReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, err.Error())
		return
	}

	result, err := h.depositService.GetDepositAddresses(c.Request.Context(), userID.(uint64), req.Currency, req.Chain)
	if err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, result)
}

func (h *DepositHandler) ListDepositOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.depositService.ListDepositOrders(c.Request.Context(), userID.(uint64), page, pageSize)
	if err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, result)
}
