package service

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"
	"motewallet/internal/config"
	dtoresp "motewallet/internal/dto/response"
	"motewallet/internal/model"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/jwt"
	"motewallet/internal/pkg/totp"
	"motewallet/internal/pkg/utils"
	"motewallet/internal/repository"
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

type AdminAuthResult struct {
	IssueSession bool
	Token        string
	Challenge    *dtoresp.AdminAuthChallengeResp
}

func (s *AdminAuthService) AdminLogin(ctx context.Context, username, password string) (*AdminAuthResult, error) {
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

	if user.TotpEnabled && user.TotpSecret != nil && *user.TotpSecret != "" {
		tempToken, err := jwt.GenerateTokenWithPurpose(user.ID, "ADMIN", user.Email, jwt.Purpose2FAVerify, s.cfg.JWT.Secret, totpChallengeTTL)
		if err != nil {
			return nil, bizerrors.ErrInternalError
		}
		return &AdminAuthResult{
			IssueSession: false,
			Challenge: &dtoresp.AdminAuthChallengeResp{
				Status:    dtoresp.AuthStatusRequires2FA,
				TempToken: tempToken,
			},
		}, nil
	}

	secret, uri, err := totp.GenerateSecret(user.Email)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}
	if err := s.adminUserRepo.UpdateFields(ctx, user.ID, map[string]interface{}{
		"totp_secret":         secret,
		"totp_enabled":        false,
		"totp_pending_secret": nil,
	}); err != nil {
		return nil, bizerrors.ErrInternalError
	}

	tempToken, err := jwt.GenerateTokenWithPurpose(user.ID, "ADMIN", user.Email, jwt.Purpose2FASetup, s.cfg.JWT.Secret, totpChallengeTTL)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	return &AdminAuthResult{
		IssueSession: false,
		Challenge: &dtoresp.AdminAuthChallengeResp{
			Status:     dtoresp.AuthStatusRequires2FASetup,
			TempToken:  tempToken,
			TotpSecret: secret,
			TotpURI:    uri,
		},
	}, nil
}

func (s *AdminAuthService) Verify2FA(ctx context.Context, tempToken, code string) (*AdminAuthResult, error) {
	claims, err := jwt.ParseToken(tempToken, s.cfg.JWT.Secret)
	if err != nil || claims.UserType != "ADMIN" || claims.Purpose != jwt.Purpose2FAVerify {
		return nil, bizerrors.ErrInvalidTokenError
	}

	user, err := s.adminUserRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	if !user.TotpEnabled || user.TotpSecret == nil || !totp.Validate(code, *user.TotpSecret) {
		return nil, bizerrors.ErrInvalidTOTPCodeE
	}

	return s.issueSession(ctx, user)
}

func (s *AdminAuthService) Confirm2FASetup(ctx context.Context, tempToken, code string) (*AdminAuthResult, error) {
	claims, err := jwt.ParseToken(tempToken, s.cfg.JWT.Secret)
	if err != nil || claims.UserType != "ADMIN" || claims.Purpose != jwt.Purpose2FASetup {
		return nil, bizerrors.ErrInvalidTokenError
	}

	user, err := s.adminUserRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	if user.TotpSecret == nil || *user.TotpSecret == "" {
		return nil, bizerrors.ErrTOTPNotEnabledE
	}
	if !totp.Validate(code, *user.TotpSecret) {
		return nil, bizerrors.ErrInvalidTOTPCodeE
	}

	if err := s.adminUserRepo.UpdateFields(ctx, user.ID, map[string]interface{}{
		"totp_enabled":        true,
		"totp_pending_secret": nil,
	}); err != nil {
		return nil, bizerrors.ErrInternalError
	}

	user.TotpEnabled = true
	return s.issueSession(ctx, user)
}

func (s *AdminAuthService) ChangePassword(ctx context.Context, tempToken, newPassword string) (*AdminAuthResult, error) {
	claims, err := jwt.ParseToken(tempToken, s.cfg.JWT.Secret)
	if err != nil || claims.UserType != "ADMIN" || claims.Purpose != jwt.PurposePasswordChange {
		return nil, bizerrors.ErrInvalidTokenError
	}

	user, err := s.adminUserRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	if !user.MustChangePassword {
		return nil, bizerrors.ErrValidationError
	}

	if utils.CheckPassword(newPassword, user.PasswordHash) {
		return nil, bizerrors.NewBusinessError(http.StatusBadRequest, bizerrors.ErrValidation, "新密码不能与当前密码相同")
	}

	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	if err := s.adminUserRepo.UpdateFields(ctx, user.ID, map[string]interface{}{
		"password_hash":        hash,
		"must_change_password": false,
	}); err != nil {
		return nil, bizerrors.ErrInternalError
	}

	user.MustChangePassword = false
	user.PasswordHash = hash
	return s.issueSession(ctx, user)
}

func (s *AdminAuthService) GetAdminByID(ctx context.Context, id uint64) (*dtoresp.AdminInfoResp, error) {
	user, err := s.adminUserRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	return toAdminInfo(user), nil
}

func (s *AdminAuthService) issueSession(ctx context.Context, user *model.AdminUser) (*AdminAuthResult, error) {
	if user.MustChangePassword {
		tempToken, err := jwt.GenerateTokenWithPurpose(
			user.ID, "ADMIN", user.Email, jwt.PurposePasswordChange, s.cfg.JWT.Secret, totpChallengeTTL,
		)
		if err != nil {
			return nil, bizerrors.ErrInternalError
		}
		return &AdminAuthResult{
			IssueSession: false,
			Challenge: &dtoresp.AdminAuthChallengeResp{
				Status:    dtoresp.AuthStatusRequiresPasswordChange,
				TempToken: tempToken,
			},
		}, nil
	}

	token, err := jwt.GenerateToken(user.ID, "ADMIN", user.Email, s.cfg.JWT.Secret, s.cfg.JWT.Expiry)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	_ = s.adminUserRepo.UpdateLastLogin(ctx, user.ID)

	return &AdminAuthResult{
		IssueSession: true,
		Token:        token,
		Challenge: &dtoresp.AdminAuthChallengeResp{
			Status: dtoresp.AuthStatusSuccess,
			Admin:  toAdminInfo(user),
		},
	}, nil
}

func toAdminInfo(user *model.AdminUser) *dtoresp.AdminInfoResp {
	return &dtoresp.AdminInfoResp{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}
}
