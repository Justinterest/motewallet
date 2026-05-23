package service

import (
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"
	"motewallet-withdrawal/backend/internal/config"
	dtoresp "motewallet-withdrawal/backend/internal/dto/response"
	"motewallet-withdrawal/backend/internal/model"
	bizerrors "motewallet-withdrawal/backend/internal/pkg/errors"
	"motewallet-withdrawal/backend/internal/pkg/jwt"
	"motewallet-withdrawal/backend/internal/pkg/kun"
	kundto "motewallet-withdrawal/backend/internal/pkg/kun/dto"
	"motewallet-withdrawal/backend/internal/pkg/utils"
	"motewallet-withdrawal/backend/internal/repository"
)

type AuthService struct {
	cfg          *config.Config
	merchantRepo repository.MerchantRepository
	kunClient    kun.KUNClient
}

func NewAuthService(cfg *config.Config, merchantRepo repository.MerchantRepository, kunClient kun.KUNClient) *AuthService {
	return &AuthService{
		cfg:          cfg,
		merchantRepo: merchantRepo,
		kunClient:    kunClient,
	}
}

type RegisterResult struct {
	Token    string
	Merchant *dtoresp.MerchantInfoResp
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*RegisterResult, error) {
	// Check if email already exists
	existing, err := s.merchantRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, bizerrors.ErrInternalError
	}
	if existing != nil {
		return nil, bizerrors.ErrEmailAlreadyExistsError
	}

	// Hash password
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	// Create merchant
	merchant := &model.Merchant{
		Email:        email,
		PasswordHash: hash,
		Status:       "PENDING_AGREEMENT",
		KycStatus:    "NONE",
	}
	if err := s.merchantRepo.Create(ctx, merchant); err != nil {
		return nil, bizerrors.ErrInternalError
	}

	var registerResp kundto.RegisterResp
	kunErr := s.kunClient.Post(ctx, "/rest/v2.0/customer/register", &kundto.RegisterReq{
		Email:     email,
		RequestNo: kun.GenerateRequestNo(),
	}, &registerResp)
	if kunErr != nil {
		slog.Error("KUN registration failed, merchant created locally", slog.String("email", email), slog.Any("error", kunErr))
	} else {
		_ = s.merchantRepo.UpdateFields(ctx, merchant.ID, map[string]interface{}{
			"kun_sub_customer_no": registerResp.SubCustomerNo,
		})
		merchant.KunSubCustomerNo = &registerResp.SubCustomerNo
	}

	// Generate JWT
	token, err := jwt.GenerateToken(merchant.ID, "MERCHANT", merchant.Email, s.cfg.JWT.Secret, s.cfg.JWT.Expiry)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	return &RegisterResult{
		Token: token,
		Merchant: &dtoresp.MerchantInfoResp{
			ID:        merchant.ID,
			Email:     merchant.Email,
			Status:    merchant.Status,
			KycStatus: merchant.KycStatus,
			CreatedAt: merchant.CreatedAt,
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
			ID:        merchant.ID,
			Email:     merchant.Email,
			Status:    merchant.Status,
			KycStatus: merchant.KycStatus,
			CreatedAt: merchant.CreatedAt,
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
		ID:        merchant.ID,
		Email:     merchant.Email,
		Status:    merchant.Status,
		KycStatus: merchant.KycStatus,
		CreatedAt: merchant.CreatedAt,
	}, nil
}
