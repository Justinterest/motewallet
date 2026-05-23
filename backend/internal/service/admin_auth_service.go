package service

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"motewallet-withdrawal/backend/internal/config"
	dtoresp "motewallet-withdrawal/backend/internal/dto/response"
	bizerrors "motewallet-withdrawal/backend/internal/pkg/errors"
	"motewallet-withdrawal/backend/internal/pkg/jwt"
	"motewallet-withdrawal/backend/internal/pkg/utils"
	"motewallet-withdrawal/backend/internal/repository"
)

type AdminAuthService struct {
	cfg           *config.Config
	adminUserRepo repository.AdminUserRepository
}

func NewAdminAuthService(cfg *config.Config, adminUserRepo repository.AdminUserRepository) *AdminAuthService {
	return &AdminAuthService{
		cfg:           cfg,
		adminUserRepo: adminUserRepo,
	}
}

type AdminLoginResult struct {
	Token string
	Admin *dtoresp.AdminInfoResp
}

func (s *AdminAuthService) AdminLogin(ctx context.Context, username, password string) (*AdminLoginResult, error) {
	user, err := s.adminUserRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrInvalidCredentialsError
		}
		return nil, bizerrors.ErrInternalError
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, bizerrors.ErrInvalidCredentialsError
	}

	if user.Status != "ACTIVE" {
		return nil, bizerrors.ErrForbiddenError
	}

	token, err := jwt.GenerateToken(user.ID, "ADMIN", user.Email, s.cfg.JWT.Secret, s.cfg.JWT.Expiry)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	// Update last login time
	_ = s.adminUserRepo.UpdateLastLogin(ctx, user.ID)

	return &AdminLoginResult{
		Token: token,
		Admin: &dtoresp.AdminInfoResp{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}

func (s *AdminAuthService) GetAdminByID(ctx context.Context, id uint64) (*dtoresp.AdminInfoResp, error) {
	user, err := s.adminUserRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	return &dtoresp.AdminInfoResp{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}
