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

type AdminExchangeHandler struct {
	adminExchangeService *service.AdminExchangeService
}

func NewAdminExchangeHandler(adminExchangeService *service.AdminExchangeService) *AdminExchangeHandler {
	return &AdminExchangeHandler{adminExchangeService: adminExchangeService}
}

func (h *AdminExchangeHandler) List(c *gin.Context) {
	var req dtoreq.AdminListExchangesReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, err.Error())
		return
	}

	result, err := h.adminExchangeService.List(c.Request.Context(), &req)
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

func (h *AdminExchangeHandler) SyncStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, "invalid order id")
		return
	}

	result, err := h.adminExchangeService.SyncStatus(c.Request.Context(), id)
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
