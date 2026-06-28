package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dtoreq "motewallet/internal/dto/request"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/response"
	"motewallet/internal/service"
)

type OnboardingHandler struct {
	onboardingService *service.OnboardingService
	kycFileService    *service.KycFileService
}

func NewOnboardingHandler(
	onboardingService *service.OnboardingService,
	kycFileService *service.KycFileService,
) *OnboardingHandler {
	return &OnboardingHandler{
		onboardingService: onboardingService,
		kycFileService:    kycFileService,
	}
}

func (h *OnboardingHandler) GetAgreements(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	result, err := h.onboardingService.GetAgreements(c.Request.Context(), userID.(uint64))
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

func (h *OnboardingHandler) SignAgreements(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	err := h.onboardingService.SignAgreements(c.Request.Context(), userID.(uint64))
	if err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, nil)
}

func (h *OnboardingHandler) SubmitKyc(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	var req dtoreq.SubmitKycReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, err.Error())
		return
	}

	err := h.onboardingService.SubmitKyc(c.Request.Context(), userID.(uint64), &req)
	if err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			if bizErr.Data != nil {
				response.ErrorWithData(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message, bizErr.Data)
			} else {
				response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			}
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, nil)
}

func (h *OnboardingHandler) GetKycStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	result, err := h.onboardingService.GetKycStatus(c.Request.Context(), userID.(uint64))
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

func (h *OnboardingHandler) PresignKycFile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	var req dtoreq.PresignKycFileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, err.Error())
		return
	}

	result, err := h.kycFileService.PresignUpload(c.Request.Context(), userID.(uint64), &req)
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

func (h *OnboardingHandler) PresignKycFileAccess(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	var req dtoreq.PresignKycFileAccessReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, err.Error())
		return
	}

	result, err := h.kycFileService.PresignAccess(c.Request.Context(), userID.(uint64), &req)
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

func (h *OnboardingHandler) ListKycCountries(c *gin.Context) {
	if _, exists := c.Get("user_id"); !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	scene := c.DefaultQuery("scene", "REGISTER_ADDRESS")
	language := c.DefaultQuery("language", "ZH_CN")
	currency := c.Query("currency")

	result, err := h.onboardingService.ListKycCountries(c.Request.Context(), scene, language, currency)
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

func (h *OnboardingHandler) ListKycCountryAuthTypes(c *gin.Context) {
	if _, exists := c.Get("user_id"); !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	countryCode := c.Query("country_code")
	if countryCode == "" {
		countryCode = c.Param("country_code")
	}

	result, err := h.onboardingService.ListKycCountryAuthTypes(c.Request.Context(), countryCode)
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
