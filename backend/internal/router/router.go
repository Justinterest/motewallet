package router

import (
	"github.com/gin-gonic/gin"
	"motewallet/internal/config"
	"motewallet/internal/handler"
	"motewallet/internal/middleware"
)

func Setup(
	cfg *config.Config,
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	adminAuthHandler *handler.AdminAuthHandler,
	onboardingHandler *handler.OnboardingHandler,
	walletHandler *handler.WalletHandler,
	feeTemplateHandler *handler.FeeTemplateHandler,
	merchantMgmtHandler *handler.MerchantManagementHandler,
	webhookHandler *handler.WebhookHandler,
	addressHandler *handler.AddressHandler,
	depositHandler *handler.DepositHandler,
	adminDepositHandler *handler.AdminDepositHandler,
	withdrawalHandler *handler.WithdrawalHandler,
	exchangeHandler *handler.ExchangeHandler,
	transferHandler *handler.TransferHandler,
	currencyConfigHandler *handler.CurrencyConfigHandler,
) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(cfg))

	// Health check
	r.GET("/health", healthHandler.Health)

	// API v1
	v1 := r.Group("/api/v1")

	// --- Public routes ---
	authPublic := v1.Group("/auth")
	{
		authPublic.POST("/send-verification-code", authHandler.SendVerificationCode)
		authPublic.POST("/register", authHandler.Register)
		authPublic.POST("/login", authHandler.Login)
	}

	adminAuthPublic := v1.Group("/admin/auth")
	{
		adminAuthPublic.POST("/login", adminAuthHandler.Login)
	}

	// --- Merchant protected routes ---
	authProtected := v1.Group("/auth")
	authProtected.Use(middleware.JWTAuth(cfg.JWT.Secret))
	{
		authProtected.POST("/logout", authHandler.Logout)
		authProtected.GET("/me", authHandler.Me)
	}

	// --- Admin protected routes ---
	adminAuthProtected := v1.Group("/admin/auth")
	adminAuthProtected.Use(middleware.AdminAuth(cfg.JWT.Secret))
	{
		adminAuthProtected.POST("/logout", adminAuthHandler.Logout)
		adminAuthProtected.GET("/me", adminAuthHandler.Me)
	}

	// --- Merchant protected: onboarding ---
	onboarding := v1.Group("/onboarding")
	onboarding.Use(middleware.JWTAuth(cfg.JWT.Secret))
	{
		onboarding.GET("/agreements", onboardingHandler.GetAgreements)
		onboarding.POST("/agreements/sign", onboardingHandler.SignAgreements)
		onboarding.POST("/kyc", onboardingHandler.SubmitKyc)
		onboarding.GET("/kyc/status", onboardingHandler.GetKycStatus)
		onboarding.POST("/files/presign", onboardingHandler.PresignKycFile)
		onboarding.POST("/files/access", onboardingHandler.PresignKycFileAccess)
		onboarding.GET("/countries", onboardingHandler.ListKycCountries)
		onboarding.GET("/countries/:country_code/auth-types", onboardingHandler.ListKycCountryAuthTypes)
	}

	// --- Merchant protected: wallet/account ---
	account := v1.Group("/account")
	account.Use(middleware.JWTAuth(cfg.JWT.Secret))
	{
		account.GET("/balances", walletHandler.GetBalances)
		account.GET("/supported-currencies", currencyConfigHandler.GetSupportedCurrencies)
		account.POST("/transfer", transferHandler.Transfer)
		account.GET("/transfers", transferHandler.ListTransfers)
	}

	// --- Merchant protected: addresses ---
	addresses := v1.Group("/addresses")
	addresses.Use(middleware.JWTAuth(cfg.JWT.Secret))
	{
		addresses.POST("/crypto", addressHandler.AddCryptoAddress)
		addresses.GET("/crypto", addressHandler.ListCryptoAddresses)
		addresses.DELETE("/crypto/:id", addressHandler.DeleteCryptoAddress)
		addresses.POST("/bank", addressHandler.AddBankAccount)
		addresses.GET("/bank", addressHandler.ListBankAccounts)
		addresses.DELETE("/bank/:id", addressHandler.DeleteBankAccount)
	}

	// --- Merchant protected: deposit ---
	deposit := v1.Group("/deposit")
	deposit.Use(middleware.JWTAuth(cfg.JWT.Secret))
	{
		deposit.GET("/addresses", depositHandler.GetDepositAddresses)
		deposit.GET("/orders", depositHandler.ListDepositOrders)
	}

	// --- Merchant protected: withdrawal ---
	withdraw := v1.Group("/withdraw")
	withdraw.Use(middleware.JWTAuth(cfg.JWT.Secret))
	{
		withdraw.POST("/crypto", withdrawalHandler.SubmitCryptoWithdrawal)
		withdraw.POST("/fiat", withdrawalHandler.SubmitFiatWithdrawal)
		withdraw.POST("/fee-preview", withdrawalHandler.PreviewWithdrawalFee)
		withdraw.GET("/orders", withdrawalHandler.ListWithdrawals)
		withdraw.GET("/orders/:id", withdrawalHandler.GetWithdrawalDetail)
	}

	// --- Merchant protected: exchange ---
	exchange := v1.Group("/exchange")
	exchange.Use(middleware.JWTAuth(cfg.JWT.Secret))
	{
		exchange.POST("/quote", exchangeHandler.GetQuote)
		exchange.POST("/order", exchangeHandler.CreateExchangeOrder)
		exchange.POST("/1to1", exchangeHandler.Create1to1Order)
		exchange.GET("/orders", exchangeHandler.ListExchangeOrders)
	}

	// --- Admin protected: fee templates ---
	adminFeeTemplates := v1.Group("/admin/fee-templates")
	adminFeeTemplates.Use(middleware.AdminAuth(cfg.JWT.Secret))
	{
		adminFeeTemplates.POST("", feeTemplateHandler.Create)
		adminFeeTemplates.GET("", feeTemplateHandler.List)
		adminFeeTemplates.GET("/:id", feeTemplateHandler.GetByID)
		adminFeeTemplates.PUT("/:id", feeTemplateHandler.Update)
		adminFeeTemplates.DELETE("/:id", feeTemplateHandler.Delete)
	}

	// --- Admin protected: merchant management ---
	adminMerchants := v1.Group("/admin/merchants")
	adminMerchants.Use(middleware.AdminAuth(cfg.JWT.Secret))
	{
		adminMerchants.GET("", merchantMgmtHandler.List)
		adminMerchants.GET("/:id", merchantMgmtHandler.GetDetail)
		adminMerchants.PUT("/:id/status", merchantMgmtHandler.UpdateStatus)
		adminMerchants.PUT("/:id/fee-template", merchantMgmtHandler.UpdateFeeTemplate)
		adminMerchants.PUT("/:id/supported-currencies", merchantMgmtHandler.UpdateSupportedCurrencies)
		adminMerchants.POST("/:id/sync-kun-balances", merchantMgmtHandler.SyncKUNBalances)
		adminMerchants.POST("/:id/sync-deposits", merchantMgmtHandler.SyncDeposits)
		adminMerchants.POST("/:id/kyc/approve", merchantMgmtHandler.ApproveKyc)
		adminMerchants.POST("/:id/kyc/reject", merchantMgmtHandler.RejectKyc)
	}

	// --- Admin protected: crypto deposits ---
	adminDeposits := v1.Group("/admin/deposits")
	adminDeposits.Use(middleware.AdminAuth(cfg.JWT.Secret))
	{
		adminDeposits.GET("", adminDepositHandler.List)
	}

	// --- Admin protected: withdrawal review ---
	adminWithdrawals := v1.Group("/admin/withdrawals")
	adminWithdrawals.Use(middleware.AdminAuth(cfg.JWT.Secret))
	{
		adminWithdrawals.GET("/pending", withdrawalHandler.AdminListPendingReviews)
		adminWithdrawals.POST("/:id/approve", withdrawalHandler.AdminApproveWithdrawal)
		adminWithdrawals.POST("/:id/reject", withdrawalHandler.AdminRejectWithdrawal)
	}

	// --- Webhook (KUN callback, no JWT, signature verified) ---
	webhook := v1.Group("/webhook")
	webhook.Use(middleware.KUNWebhookAuth(cfg.KUN))
	{
		webhook.POST("/kun", webhookHandler.HandleKUN)
	}

	return r
}
