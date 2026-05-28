package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"motewallet/internal/config"
	dtoreq "motewallet/internal/dto/request"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/jwt"
	"motewallet/internal/pkg/response"
	"motewallet/internal/service"
)

type AdminAuthHandler struct {
	cfg              *config.Config
	adminAuthService *service.AdminAuthService
}

func NewAdminAuthHandler(cfg *config.Config, adminAuthService *service.AdminAuthService) *AdminAuthHandler {
	return &AdminAuthHandler{
		cfg:              cfg,
		adminAuthService: adminAuthService,
	}
}

func (h *AdminAuthHandler) Login(c *gin.Context) {
	var req dtoreq.AdminLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, err.Error())
		return
	}

	result, err := h.adminAuthService.AdminLogin(c.Request.Context(), req.Username, req.Password)
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
	jwt.SetTokenCookie(c, jwt.CookieAdminToken, result.Token, maxAge, secure)
	response.Success(c, result.Admin)
}

func (h *AdminAuthHandler) Logout(c *gin.Context) {
	secure := h.cfg.IsProduction()
	jwt.ClearTokenCookie(c, jwt.CookieAdminToken, secure)
	response.Success(c, nil)
}

func (h *AdminAuthHandler) Me(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	admin, err := h.adminAuthService.GetAdminByID(c.Request.Context(), adminID.(uint64))
	if err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, admin)
}
