package service

import (
	"context"
	"errors"
	"log/slog"

	"motewallet/internal/config"
	dtoresp "motewallet/internal/dto/response"
	"motewallet/internal/model"
	"motewallet/internal/pkg/email"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/jwt"
	"motewallet/internal/pkg/kun"
	kundto "motewallet/internal/pkg/kun/dto"
	"motewallet/internal/pkg/utils"
	"motewallet/internal/pkg/verification"
	"motewallet/internal/repository"

	"gorm.io/gorm"
)

type AuthService struct {
	cfg               *config.Config
	merchantRepo      repository.MerchantRepository
	feeTemplateRepo   repository.FeeTemplateRepository
	kunClient         kun.KUNClient
	verificationStore *verification.Store
	emailSender       *email.Sender
}

func NewAuthService(cfg *config.Config, merchantRepo repository.MerchantRepository, feeTemplateRepo repository.FeeTemplateRepository, kunClient kun.KUNClient, emailSender *email.Sender) *AuthService {
	return &AuthService{
		cfg:               cfg,
		merchantRepo:      merchantRepo,
		feeTemplateRepo:   feeTemplateRepo,
		kunClient:         kunClient,
		verificationStore: verification.NewStore(),
		emailSender:       emailSender,
	}
}

type RegisterResult struct {
	Token    string
	Merchant *dtoresp.MerchantInfoResp
}

func (s *AuthService) SendVerificationCode(ctx context.Context, emailAddr string) error {
	existing, err := s.merchantRepo.FindByEmail(ctx, emailAddr)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return bizerrors.ErrInternalError
	}
	if existing != nil {
		return bizerrors.ErrEmailAlreadyExistsError
	}

	code, retryAfter, err := s.verificationStore.Issue(emailAddr)
	if err != nil {
		if errors.Is(err, verification.ErrSendTooFrequent) {
			slog.Debug("verification code send throttled",
				slog.String("email", emailAddr),
				slog.Duration("retry_after", retryAfter),
			)
			return bizerrors.ErrVerificationSendTooFrequentE
		}
		return bizerrors.ErrInternalError
	}

	if err := s.emailSender.SendVerificationCode(emailAddr, code); err != nil {
		slog.Error("failed to send verification email", slog.String("email", emailAddr), slog.Any("error", err))
		return bizerrors.ErrInternalError
	}

	return nil
}

func (s *AuthService) Register(ctx context.Context, email, password, verificationCode string) (*RegisterResult, error) {
	// if err := s.verificationStore.Verify(email, verificationCode); err != nil {
	// 	switch {
	// 	case errors.Is(err, verification.ErrExpiredCode):
	// 		return nil, bizerrors.ErrVerificationCodeExpiredE
	// 	default:
	// 		return nil, bizerrors.ErrInvalidVerificationCodeE
	// 	}
	// }

	existing, err := s.merchantRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, bizerrors.ErrInternalError
	}
	if existing != nil {
		return nil, bizerrors.ErrEmailAlreadyExistsError
	}

	var registerResp kundto.RegisterResp
	if err := s.kunClient.Post(ctx, "/rest/v2.0/customer/register", &kundto.RegisterReq{
		Email:     email,
		RequestNo: kun.GenerateRequestNo(),
	}, &registerResp); err != nil {
		slog.Error("KUN sub-merchant registration failed", slog.String("email", email), slog.Any("error", err))
		return nil, bizerrors.ErrKUNAPIFailedE
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	merchant := &model.Merchant{
		Email:            email,
		PasswordHash:     hash,
		KunSubCustomerNo: &registerResp.SubCustomerNo,
		Status:           "PENDING_AGREEMENT",
		KycStatus:        "NONE",
	}

	defaultTemplate, err := s.feeTemplateRepo.FindDefault(ctx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("failed to find default fee template", slog.Any("error", err))
	}
	if defaultTemplate != nil {
		merchant.FeeTemplateID = &defaultTemplate.ID
	}

	if err := s.merchantRepo.Create(ctx, merchant); err != nil {
		return nil, bizerrors.ErrInternalError
	}

	token, err := jwt.GenerateToken(merchant.ID, "MERCHANT", merchant.Email, s.cfg.JWT.Secret, s.cfg.JWT.Expiry)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	return &RegisterResult{
		Token: token,
		Merchant: &dtoresp.MerchantInfoResp{
			ID:            merchant.ID,
			Email:         merchant.Email,
			Status:        merchant.Status,
			KycStatus:     merchant.KycStatus,
			FeeTemplateID: merchant.FeeTemplateID,
			CreatedAt:     merchant.CreatedAt,
		},
	}, nil
}

type LoginResult struct {
	Token    string
	Merchant *dtoresp.MerchantInfoResp
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	merchant, err := s.merchantRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrInvalidCredentialsError
		}
		return nil, bizerrors.ErrInternalError
	}

	if !utils.CheckPassword(password, merchant.PasswordHash) {
		return nil, bizerrors.ErrInvalidCredentialsError
	}

	token, err := jwt.GenerateToken(merchant.ID, "MERCHANT", merchant.Email, s.cfg.JWT.Secret, s.cfg.JWT.Expiry)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	return &LoginResult{
		Token: token,
		Merchant: &dtoresp.MerchantInfoResp{
			ID:            merchant.ID,
			Email:         merchant.Email,
			Status:        merchant.Status,
			KycStatus:     merchant.KycStatus,
			FeeTemplateID: merchant.FeeTemplateID,
			CreatedAt:     merchant.CreatedAt,
		},
	}, nil
}

func (s *AuthService) GetMerchantByID(ctx context.Context, id uint64) (*dtoresp.MerchantInfoResp, error) {
	merchant, err := s.merchantRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	return &dtoresp.MerchantInfoResp{
		ID:            merchant.ID,
		Email:         merchant.Email,
		Status:        merchant.Status,
		KycStatus:     merchant.KycStatus,
		FeeTemplateID: merchant.FeeTemplateID,
		CreatedAt:     merchant.CreatedAt,
	}, nil
}
