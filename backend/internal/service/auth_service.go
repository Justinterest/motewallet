package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"motewallet/internal/config"
	dtoresp "motewallet/internal/dto/response"
	"motewallet/internal/model"
	"motewallet/internal/pkg/email"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/jwt"
	"motewallet/internal/pkg/kun"
	kundto "motewallet/internal/pkg/kun/dto"
	"motewallet/internal/pkg/totp"
	"motewallet/internal/pkg/utils"
	"motewallet/internal/pkg/verification"
	"motewallet/internal/repository"

	"gorm.io/gorm"
)

const totpChallengeTTL = 10 * time.Minute

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

type AuthResult struct {
	IssueSession bool
	Token        string
	Challenge    *dtoresp.AuthChallengeResp
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

func (s *AuthService) Register(ctx context.Context, email, password, verificationCode string) (*AuthResult, error) {
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

	secret, uri, err := totp.GenerateSecret(email)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	merchant := &model.Merchant{
		Email:            email,
		PasswordHash:     hash,
		TotpSecret:       &secret,
		TotpEnabled:      false,
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

	tempToken, err := jwt.GenerateTokenWithPurpose(merchant.ID, "MERCHANT", merchant.Email, jwt.Purpose2FASetup, s.cfg.JWT.Secret, totpChallengeTTL)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	return &AuthResult{
		IssueSession: false,
		Challenge: &dtoresp.AuthChallengeResp{
			Status:     dtoresp.AuthStatusRequires2FASetup,
			TempToken:  tempToken,
			TotpSecret: secret,
			TotpURI:    uri,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
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

	if merchant.TotpEnabled && merchant.TotpSecret != nil && *merchant.TotpSecret != "" {
		tempToken, err := jwt.GenerateTokenWithPurpose(merchant.ID, "MERCHANT", merchant.Email, jwt.Purpose2FAVerify, s.cfg.JWT.Secret, totpChallengeTTL)
		if err != nil {
			return nil, bizerrors.ErrInternalError
		}
		return &AuthResult{
			IssueSession: false,
			Challenge: &dtoresp.AuthChallengeResp{
				Status:    dtoresp.AuthStatusRequires2FA,
				TempToken: tempToken,
			},
		}, nil
	}

	// 2FA not enabled (new account incomplete setup, or admin reset) → force setup
	secret, uri, err := totp.GenerateSecret(merchant.Email)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}
	if err := s.merchantRepo.UpdateFields(ctx, merchant.ID, map[string]interface{}{
		"totp_secret":         secret,
		"totp_enabled":        false,
		"totp_pending_secret": nil,
	}); err != nil {
		return nil, bizerrors.ErrInternalError
	}

	tempToken, err := jwt.GenerateTokenWithPurpose(merchant.ID, "MERCHANT", merchant.Email, jwt.Purpose2FASetup, s.cfg.JWT.Secret, totpChallengeTTL)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	return &AuthResult{
		IssueSession: false,
		Challenge: &dtoresp.AuthChallengeResp{
			Status:     dtoresp.AuthStatusRequires2FASetup,
			TempToken:  tempToken,
			TotpSecret: secret,
			TotpURI:    uri,
		},
	}, nil
}

func (s *AuthService) Verify2FA(ctx context.Context, tempToken, code string) (*AuthResult, error) {
	claims, err := jwt.ParseToken(tempToken, s.cfg.JWT.Secret)
	if err != nil || claims.UserType != "MERCHANT" || claims.Purpose != jwt.Purpose2FAVerify {
		return nil, bizerrors.ErrInvalidTokenError
	}

	merchant, err := s.merchantRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	if !merchant.TotpEnabled || merchant.TotpSecret == nil || !totp.Validate(code, *merchant.TotpSecret) {
		return nil, bizerrors.ErrInvalidTOTPCodeE
	}

	return s.issueSession(merchant)
}

func (s *AuthService) Confirm2FASetup(ctx context.Context, tempToken, code string) (*AuthResult, error) {
	claims, err := jwt.ParseToken(tempToken, s.cfg.JWT.Secret)
	if err != nil || claims.UserType != "MERCHANT" || claims.Purpose != jwt.Purpose2FASetup {
		return nil, bizerrors.ErrInvalidTokenError
	}

	merchant, err := s.merchantRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	if merchant.TotpSecret == nil || *merchant.TotpSecret == "" {
		return nil, bizerrors.ErrTOTPNotEnabledE
	}
	if !totp.Validate(code, *merchant.TotpSecret) {
		return nil, bizerrors.ErrInvalidTOTPCodeE
	}

	if err := s.merchantRepo.UpdateFields(ctx, merchant.ID, map[string]interface{}{
		"totp_enabled":        true,
		"totp_pending_secret": nil,
	}); err != nil {
		return nil, bizerrors.ErrInternalError
	}

	merchant.TotpEnabled = true
	return s.issueSession(merchant)
}

func (s *AuthService) GetTotpStatus(ctx context.Context, merchantID uint64) (*dtoresp.TotpStatusResp, error) {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}
	return &dtoresp.TotpStatusResp{Enabled: merchant.TotpEnabled}, nil
}

func (s *AuthService) PrepareTotpRebind(ctx context.Context, merchantID uint64, currentCode string) (*dtoresp.TotpSetupResp, error) {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	if !merchant.TotpEnabled || merchant.TotpSecret == nil {
		return nil, bizerrors.ErrTOTPNotEnabledE
	}
	if !totp.Validate(currentCode, *merchant.TotpSecret) {
		return nil, bizerrors.ErrInvalidTOTPCodeE
	}

	secret, uri, err := totp.GenerateSecret(merchant.Email)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	if err := s.merchantRepo.UpdateFields(ctx, merchant.ID, map[string]interface{}{
		"totp_pending_secret": secret,
	}); err != nil {
		return nil, bizerrors.ErrInternalError
	}

	return &dtoresp.TotpSetupResp{
		TotpSecret: secret,
		TotpURI:    uri,
	}, nil
}

func (s *AuthService) ConfirmTotpRebind(ctx context.Context, merchantID uint64, code string) error {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrNotFoundError
		}
		return bizerrors.ErrInternalError
	}

	if merchant.TotpPendingSecret == nil || *merchant.TotpPendingSecret == "" {
		return bizerrors.ErrValidationError
	}
	if !totp.Validate(code, *merchant.TotpPendingSecret) {
		return bizerrors.ErrInvalidTOTPCodeE
	}

	return s.merchantRepo.UpdateFields(ctx, merchant.ID, map[string]interface{}{
		"totp_secret":         *merchant.TotpPendingSecret,
		"totp_enabled":        true,
		"totp_pending_secret": nil,
	})
}

func (s *AuthService) GetMerchantByID(ctx context.Context, id uint64) (*dtoresp.MerchantInfoResp, error) {
	merchant, err := s.merchantRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	return toMerchantInfo(merchant), nil
}

func (s *AuthService) issueSession(merchant *model.Merchant) (*AuthResult, error) {
	token, err := jwt.GenerateToken(merchant.ID, "MERCHANT", merchant.Email, s.cfg.JWT.Secret, s.cfg.JWT.Expiry)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	return &AuthResult{
		IssueSession: true,
		Token:        token,
		Challenge: &dtoresp.AuthChallengeResp{
			Status:   dtoresp.AuthStatusSuccess,
			Merchant: toMerchantInfo(merchant),
		},
	}, nil
}

func toMerchantInfo(merchant *model.Merchant) *dtoresp.MerchantInfoResp {
	return &dtoresp.MerchantInfoResp{
		ID:            merchant.ID,
		Email:         merchant.Email,
		Status:        merchant.Status,
		KycStatus:     merchant.KycStatus,
		TotpEnabled:   merchant.TotpEnabled,
		FeeTemplateID: merchant.FeeTemplateID,
		CreatedAt:     merchant.CreatedAt,
	}
}
