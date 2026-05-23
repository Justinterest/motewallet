package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"motewallet-withdrawal/backend/internal/config"
	dtoreq "motewallet-withdrawal/backend/internal/dto/request"
	bizerrors "motewallet-withdrawal/backend/internal/pkg/errors"
	"motewallet-withdrawal/backend/internal/pkg/jwt"
	"motewallet-withdrawal/backend/internal/pkg/response"
	"motewallet-withdrawal/backend/internal/service"
)

type AuthHandler struct {
	cfg         *config.Config
	authService *service.AuthService
}

func NewAuthHandler(cfg *config.Config, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		cfg:         cfg,
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dtoreq.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, err.Error())
		return
	}

	result, err := h.authService.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	secure := h.cfg.IsProduction()
	maxAge := int(h.cfg.JWT.Expiry.Seconds())
	jwt.SetTokenCookie(c, jwt.CookieMerchantToken, result.Token, maxAge, secure)
	response.Success(c, result.Merchant)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dtoreq.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, err.Error())
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	secure := h.cfg.IsProduction()
	maxAge := int(h.cfg.JWT.Expiry.Seconds())
	jwt.SetTokenCookie(c, jwt.CookieMerchantToken, result.Token, maxAge, secure)
	response.Success(c, result.Merchant)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	secure := h.cfg.IsProduction()
	jwt.ClearTokenCookie(c, jwt.CookieMerchantToken, secure)
	response.Success(c, nil)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	merchant, err := h.authService.GetMerchantByID(c.Request.Context(), userID.(uint64))
	if err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, merchant)
}
