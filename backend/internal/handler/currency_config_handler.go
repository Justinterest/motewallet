package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"errors"
	"gorm.io/gorm"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/response"
	"motewallet/internal/repository"
	"motewallet/internal/service"
)

type CurrencyConfigHandler struct {
	merchantRepo      repository.MerchantRepository
	currencyConfigSvc *service.CurrencyConfigService
}

func NewCurrencyConfigHandler(
	merchantRepo repository.MerchantRepository,
	currencyConfigSvc *service.CurrencyConfigService,
) *CurrencyConfigHandler {
	return &CurrencyConfigHandler{
		merchantRepo:      merchantRepo,
		currencyConfigSvc: currencyConfigSvc,
	}
}

func (h *CurrencyConfigHandler) GetSupportedCurrencies(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	merchant, err := h.merchantRepo.FindByID(c.Request.Context(), userID.(uint64))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, bizerrors.ErrNotFound, "merchant not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	crypto, err := h.currencyConfigSvc.GetSupportedCrypto(c.Request.Context(), merchant)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}
	fiat, err := h.currencyConfigSvc.GetSupportedFiat(c.Request.Context(), merchant)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, map[string][]string{
		"crypto_currencies": crypto,
		"fiat_currencies":   fiat,
	})
}
