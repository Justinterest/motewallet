package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	dtoreq "motewallet/internal/dto/request"
	dtoresp "motewallet/internal/dto/response"
	"motewallet/internal/pkg/currency"
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
	chains, err := h.currencyConfigSvc.GetSupportedChains(c.Request.Context(), merchant)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}
	defaults, err := h.currencyConfigSvc.GetDefaultChains(c.Request.Context(), merchant)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, dtoresp.SupportedCurrenciesResp{
		CryptoCurrencies: crypto,
		FiatCurrencies:   fiat,
		CryptoChains:     chains,
		DefaultChains:    defaults,
	})
}

func (h *CurrencyConfigHandler) GetSystemCurrencyConfig(c *gin.Context) {
	crypto, err := h.currencyConfigSvc.GetAvailableCrypto(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}
	fiat, err := h.currencyConfigSvc.GetAvailableFiat(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}
	chains, err := h.currencyConfigSvc.GetAvailableChains(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}
	defaults, err := h.currencyConfigSvc.GetAvailableDefaultChains(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, dtoresp.SystemCurrencyConfigResp{
		CryptoCurrencies: crypto,
		FiatCurrencies:   fiat,
		CryptoChains:     chains,
		DefaultChains:    defaults,
		CatalogChains:    h.currencyConfigSvc.GetCatalogChains(),
		AllCrypto:        append([]string(nil), currency.AllCrypto...),
		AllFiat:          append([]string(nil), currency.AllFiat...),
	})
}

func (h *CurrencyConfigHandler) UpdateSystemCurrencyConfig(c *gin.Context) {
	if _, exists := c.Get("admin_id"); !exists {
		response.Error(c, http.StatusUnauthorized, bizerrors.ErrUnauthorized, "unauthorized")
		return
	}

	var req dtoreq.UpdateSystemCurrencyConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, bizerrors.ErrValidation, err.Error())
		return
	}

	if err := h.currencyConfigSvc.UpdateGlobalConfig(
		c.Request.Context(),
		req.CryptoCurrencies,
		req.FiatCurrencies,
		req.CryptoChains,
		req.DefaultChains,
	); err != nil {
		if bizErr, ok := err.(*bizerrors.BusinessError); ok {
			response.Error(c, bizErr.HTTPStatus, bizErr.Code, bizErr.Message)
			return
		}
		response.Error(c, http.StatusInternalServerError, bizerrors.ErrInternal, "internal server error")
		return
	}

	response.Success(c, nil)
}
