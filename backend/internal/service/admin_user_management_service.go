package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	dtoreq "motewallet/internal/dto/request"
	dtoresp "motewallet/internal/dto/response"
	"motewallet/internal/model"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/utils"
	"motewallet/internal/repository"
)

type AdminUserManagementService struct {
	adminUserRepo repository.AdminUserRepository
	auditLogRepo  repository.AuditLogRepository
}

func NewAdminUserManagementService(
	adminUserRepo repository.AdminUserRepository,
	auditLogRepo repository.AuditLogRepository,
) *AdminUserManagementService {
	return &AdminUserManagementService{
		adminUserRepo: adminUserRepo,
		auditLogRepo:  auditLogRepo,
	}
}

func (s *AdminUserManagementService) List(ctx context.Context, operatorID uint64) ([]dtoresp.AdminUserItemResp, error) {
	if err := s.requireSuperAdmin(ctx, operatorID); err != nil {
		return nil, err
	}

	users, err := s.adminUserRepo.List(ctx)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	items := make([]dtoresp.AdminUserItemResp, 0, len(users))
	for i := range users {
		items = append(items, toAdminUserItem(&users[i]))
	}
	return items, nil
}

func (s *AdminUserManagementService) Create(ctx context.Context, operatorID uint64, req *dtoreq.CreateAdminUserReq) (*dtoresp.CreateAdminUserResp, error) {
	if err := s.requireSuperAdmin(ctx, operatorID); err != nil {
		return nil, err
	}

	if _, err := s.adminUserRepo.FindByUsername(ctx, req.Username); err == nil {
		return nil, bizerrors.ErrUsernameAlreadyExistsError
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, bizerrors.ErrInternalError
	}

	if _, err := s.adminUserRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, bizerrors.ErrEmailAlreadyExistsError
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, bizerrors.ErrInternalError
	}

	password := req.Password
	if password == "" {
		var err error
		password, err = utils.GenerateRandomPassword(12)
		if err != nil {
			return nil, bizerrors.ErrInternalError
		}
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	user := &model.AdminUser{
		Username:           req.Username,
		Email:              req.Email,
		PasswordHash:       hash,
		MustChangePassword: true,
		Role:               req.Role,
		Status:             "ACTIVE",
	}

	if err := s.adminUserRepo.Create(ctx, user); err != nil {
		return nil, bizerrors.ErrInternalError
	}

	s.logAudit(ctx, operatorID, "CREATE_ADMIN_USER", "AdminUser", fmt.Sprintf("%d", user.ID), map[string]string{
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})

	resp := &dtoresp.CreateAdminUserResp{
		User:            toAdminUserItem(user),
		InitialPassword: password,
	}
	return resp, nil
}

func (s *AdminUserManagementService) ResetPassword(ctx context.Context, operatorID, targetID uint64) (*dtoresp.ResetAdminPasswordResp, error) {
	if err := s.requireSuperAdmin(ctx, operatorID); err != nil {
		return nil, err
	}

	user, err := s.adminUserRepo.FindByID(ctx, targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	password, err := utils.GenerateRandomPassword(12)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	if err := s.adminUserRepo.UpdateFields(ctx, targetID, map[string]interface{}{
		"password_hash":        hash,
		"must_change_password": true,
	}); err != nil {
		return nil, bizerrors.ErrInternalError
	}

	s.logAudit(ctx, operatorID, "RESET_ADMIN_PASSWORD", "AdminUser", fmt.Sprintf("%d", targetID), map[string]string{
		"username": user.Username,
	})

	return &dtoresp.ResetAdminPasswordResp{NewPassword: password}, nil
}

func (s *AdminUserManagementService) Reset2FA(ctx context.Context, operatorID, targetID uint64) error {
	if err := s.requireSuperAdmin(ctx, operatorID); err != nil {
		return err
	}

	user, err := s.adminUserRepo.FindByID(ctx, targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrNotFoundError
		}
		return bizerrors.ErrInternalError
	}

	if err := s.adminUserRepo.UpdateFields(ctx, targetID, map[string]interface{}{
		"totp_secret":         nil,
		"totp_enabled":        false,
		"totp_pending_secret": nil,
	}); err != nil {
		return bizerrors.ErrInternalError
	}

	s.logAudit(ctx, operatorID, "RESET_ADMIN_2FA", "AdminUser", fmt.Sprintf("%d", targetID), map[string]string{
		"username":         user.Username,
		"previous_enabled": fmt.Sprintf("%v", user.TotpEnabled),
	})

	return nil
}

func (s *AdminUserManagementService) requireSuperAdmin(ctx context.Context, operatorID uint64) error {
	operator, err := s.adminUserRepo.FindByID(ctx, operatorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrForbiddenError
		}
		return bizerrors.ErrInternalError
	}
	if operator.Role != "SUPER_ADMIN" {
		return bizerrors.ErrForbiddenError
	}
	return nil
}

func (s *AdminUserManagementService) logAudit(ctx context.Context, operatorID uint64, action, targetType, targetID string, detail interface{}) {
	var detailJSON json.RawMessage
	if detail != nil {
		data, err := json.Marshal(detail)
		if err == nil {
			detailJSON = data
		}
	}
	tt := targetType
	tid := targetID
	_ = s.auditLogRepo.Create(ctx, &model.AuditLog{
		OperatorID:   operatorID,
		OperatorType: "ADMIN",
		Action:       action,
		TargetType:   &tt,
		TargetID:     &tid,
		Detail:       detailJSON,
	})
}

func toAdminUserItem(user *model.AdminUser) dtoresp.AdminUserItemResp {
	return dtoresp.AdminUserItemResp{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
		Status:      user.Status,
		TotpEnabled: user.TotpEnabled,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
	}
}
