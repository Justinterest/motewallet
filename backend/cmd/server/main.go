package main

import (
	"log"
	"log/slog"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"motewallet/internal/config"
	"motewallet/internal/handler"
	"motewallet/internal/pkg/email"
	"motewallet/internal/pkg/kun"
	"motewallet/internal/pkg/storage"
	"motewallet/internal/repository"
	"motewallet/internal/router"
	"motewallet/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logLevel := slog.LevelInfo
	if cfg.App.Env == "development" {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	gormLogLevel := gormlogger.Warn
	if cfg.App.Env == "development" {
		gormLogLevel = gormlogger.Info
	}
	db, err := gorm.Open(mysql.Open(cfg.DB.DSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Repositories
	merchantRepo := repository.NewMerchantRepository(db)
	merchantKycSubmissionRepo := repository.NewMerchantKycSubmissionRepository(db)
	adminUserRepo := repository.NewAdminUserRepository(db)
	merchantWalletRepo := repository.NewMerchantWalletRepository(db)
	feeTemplateRepo := repository.NewFeeTemplateRepository(db)
	exchangeItemRepo := repository.NewFeeTemplateExchangeItemRepository(db)
	cryptoWithdrawalItemRepo := repository.NewFeeTemplateCryptoWithdrawalItemRepository(db)
	fiatWithdrawalItemRepo := repository.NewFeeTemplateFiatWithdrawalItemRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)
	webhookLogRepo := repository.NewWebhookLogRepository(db)
	transactionRecordRepo := repository.NewTransactionRecordRepository(db)
	cryptoAddressRepo := repository.NewCryptoAddressRepository(db)
	bankAccountRepo := repository.NewBankAccountRepository(db)
	depositOrderRepo := repository.NewDepositOrderRepository(db)
	withdrawalOrderRepo := repository.NewWithdrawalOrderRepository(db)
	exchangeOrderRepo := repository.NewExchangeOrderRepository(db)
	transferOrderRepo := repository.NewTransferOrderRepository(db)

	// KUN Client
	var kunClient kun.KUNClient
	if cfg.KUN.MockEnabled {
		slog.Info("using mock KUN client")
		kunClient = kun.NewMockClient(cfg.KUN)
	} else {
		kunClient = kun.NewClient(cfg.KUN)
	}

	s3Storage, err := storage.NewS3Storage(cfg.S3)
	if err != nil {
		log.Fatalf("failed to init S3 storage: %v", err)
	}
	if !s3Storage.Enabled() {
		slog.Warn("S3 storage is not configured; KYC file upload will be unavailable")
	}

	emailSender, err := email.NewSender(email.Config{
		AccessKeyID:     cfg.DirectMail.AccessKeyID,
		AccessKeySecret: cfg.DirectMail.AccessKeySecret,
		Region:          cfg.DirectMail.Region,
		Endpoint:        cfg.DirectMail.Endpoint,
		AccountName:     cfg.DirectMail.AccountName,
		FromAlias:       cfg.DirectMail.FromAlias,
		TagName:         cfg.DirectMail.TagName,
	})
	if err != nil {
		log.Fatalf("failed to init DirectMail sender: %v", err)
	}
	if !emailSender.Enabled() {
		slog.Warn("Aliyun DirectMail is not configured; verification codes will be logged to console")
	}

	// Services
	authService := service.NewAuthService(cfg, merchantRepo, feeTemplateRepo, kunClient, emailSender)
	adminAuthService := service.NewAdminAuthService(cfg, adminUserRepo)
	walletService := service.NewWalletService(merchantWalletRepo)
	kycFileService := service.NewKycFileService(cfg, merchantRepo, s3Storage, kunClient)
	onboardingService := service.NewOnboardingService(cfg, merchantRepo, merchantKycSubmissionRepo, walletService, kycFileService, kunClient)
	feeTemplateService := service.NewFeeTemplateService(db, feeTemplateRepo, exchangeItemRepo, cryptoWithdrawalItemRepo, fiatWithdrawalItemRepo, auditLogRepo)
	merchantMgmtService := service.NewMerchantManagementService(merchantRepo, merchantWalletRepo, feeTemplateRepo, auditLogRepo)
	addressService := service.NewAddressService(kunClient, merchantRepo, cryptoAddressRepo, bankAccountRepo)
	depositService := service.NewDepositService(kunClient, merchantRepo)
	adminDepositService := service.NewAdminDepositService(depositOrderRepo)
	withdrawalService := service.NewWithdrawalService(db, merchantRepo, walletService, withdrawalOrderRepo, transactionRecordRepo, cryptoWithdrawalItemRepo, fiatWithdrawalItemRepo, bankAccountRepo, kunClient)
	exchangeService := service.NewExchangeService(db, merchantRepo, walletService, exchangeOrderRepo, transactionRecordRepo, exchangeItemRepo, kunClient)
	transferService := service.NewTransferService(db, merchantRepo, walletService, transferOrderRepo, transactionRecordRepo, kunClient)
	webhookService := service.NewWebhookService(db, webhookLogRepo, merchantRepo, walletService, transactionRecordRepo, depositOrderRepo, withdrawalOrderRepo, exchangeOrderRepo, transferOrderRepo)

	// Handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(cfg, authService)
	adminAuthHandler := handler.NewAdminAuthHandler(cfg, adminAuthService)
	onboardingHandler := handler.NewOnboardingHandler(onboardingService, kycFileService)
	walletHandler := handler.NewWalletHandler(walletService)
	feeTemplateHandler := handler.NewFeeTemplateHandler(feeTemplateService)
	merchantMgmtHandler := handler.NewMerchantManagementHandler(merchantMgmtService)
	webhookHandler := handler.NewWebhookHandler(webhookService)
	addressHandler := handler.NewAddressHandler(addressService)
	depositHandler := handler.NewDepositHandler(depositService)
	adminDepositHandler := handler.NewAdminDepositHandler(adminDepositService)
	withdrawalHandler := handler.NewWithdrawalHandler(withdrawalService)
	exchangeHandler := handler.NewExchangeHandler(exchangeService)
	transferHandler := handler.NewTransferHandler(transferService)

	// Router
	r := router.Setup(
		cfg,
		healthHandler,
		authHandler,
		adminAuthHandler,
		onboardingHandler,
		walletHandler,
		feeTemplateHandler,
		merchantMgmtHandler,
		webhookHandler,
		addressHandler,
		depositHandler,
		adminDepositHandler,
		withdrawalHandler,
		exchangeHandler,
		transferHandler,
	)

	addr := ":" + cfg.App.Port
	slog.Info("starting server", slog.String("addr", addr), slog.String("env", cfg.App.Env))
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
