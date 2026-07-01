package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dtoreq "motewallet/internal/dto/request"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/response"
	"motewallet/internal/service"
)

type AdminWithdrawalHandler struct {
	adminWithdrawalService *service.AdminWithdrawalService
}

func NewAdminWithdrawalHandler(adminWithdrawalService *service.AdminWithdrawalService) *AdminWithdrawalHandler {
	return &AdminWithdrawalHandler{adminWithdrawalService: adminWithdrawalService}
}

func (h *AdminWithdrawalHandler) List(c *gin.Context) {
	var req dtoreq.AdminListWithdrawalsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, err.Error())
		return
	}

	result, err := h.adminWithdrawalService.List(c.Request.Context(), &req)
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
