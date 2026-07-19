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

type MerchantManagementHandler struct {
	merchantMgmtService *service.MerchantManagementService
}

func NewMerchantManagementHandler(merchantMgmtService *service.MerchantManagementService) *MerchantManagementHandler {
	return &MerchantManagementHandler{
		merchantMgmtService: merchantMgmtService,
	}
}

func (h *MerchantManagementHandler) List(c *gin.Context) {
	var req dtoreq.AdminListMerchantsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, err.Error())
		return
	}

	items, total, err := h.merchantMgmtService.List(c.Request.Context(), &req)
	if err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Paginated(c, items, total, req.Page, req.PageSize)
}

func (h *MerchantManagementHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, "invalid id")
		return
	}

	result, err := h.merchantMgmtService.GetDetail(c.Request.Context(), id)
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

func (h *MerchantManagementHandler) UpdateStatus(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, "invalid id")
		return
	}

	var req dtoreq.UpdateMerchantStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, err.Error())
		return
	}

	if err := h.merchantMgmtService.UpdateStatus(c.Request.Context(), adminID.(uint64), id, &req); err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, nil)
}

func (h *MerchantManagementHandler) UpdateFeeTemplate(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, "invalid id")
		return
	}

	var req dtoreq.UpdateMerchantFeeTemplateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, err.Error())
		return
	}

	if err := h.merchantMgmtService.UpdateFeeTemplate(c.Request.Context(), adminID.(uint64), id, &req); err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, nil)
}

func (h *MerchantManagementHandler) UpdateSupportedCurrencies(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, "invalid id")
		return
	}

	var req dtoreq.UpdateMerchantSupportedCurrenciesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, err.Error())
		return
	}

	if err := h.merchantMgmtService.UpdateSupportedCurrencies(c.Request.Context(), adminID.(uint64), id, &req); err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, nil)
}

func (h *MerchantManagementHandler) SyncKUNBalances(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, "invalid id")
		return
	}

	result, err := h.merchantMgmtService.SyncKUNBalances(c.Request.Context(), adminID.(uint64), id)
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

func (h *MerchantManagementHandler) SyncDeposits(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, "invalid id")
		return
	}

	result, err := h.merchantMgmtService.SyncDeposits(c.Request.Context(), adminID.(uint64), id)
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

func (h *MerchantManagementHandler) ApproveKyc(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, "invalid id")
		return
	}

	if err := h.merchantMgmtService.ApproveKyc(c.Request.Context(), adminID.(uint64), id); err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, nil)
}

func (h *MerchantManagementHandler) RejectKyc(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, "invalid id")
		return
	}

	var req dtoreq.KycRejectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, err.Error())
		return
	}

	if err := h.merchantMgmtService.RejectKyc(c.Request.Context(), adminID.(uint64), id, &req); err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, nil)
}

func (h *MerchantManagementHandler) Reset2FA(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, "invalid id")
		return
	}

	if err := h.merchantMgmtService.Reset2FA(c.Request.Context(), adminID.(uint64), id); err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, nil)
}

func (h *MerchantManagementHandler) ListLedger(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, "invalid id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.merchantMgmtService.ListLedger(
		c.Request.Context(),
		id,
		c.Query("account_type"),
		c.Query("currency"),
		c.Query("biz_type"),
		c.Query("entry_type"),
		page,
		pageSize,
	)
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
