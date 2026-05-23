package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	dtoreq "motewallet-withdrawal/backend/internal/dto/request"
	dtoresp "motewallet-withdrawal/backend/internal/dto/response"
	"motewallet-withdrawal/backend/internal/model"
	bizerrors "motewallet-withdrawal/backend/internal/pkg/errors"
	"motewallet-withdrawal/backend/internal/repository"
)

type MerchantManagementService struct {
	merchantRepo       repository.MerchantRepository
	merchantWalletRepo repository.MerchantWalletRepository
	feeTemplateRepo    repository.FeeTemplateRepository
	auditLogRepo       repository.AuditLogRepository
}

func NewMerchantManagementService(
	merchantRepo repository.MerchantRepository,
	merchantWalletRepo repository.MerchantWalletRepository,
	feeTemplateRepo repository.FeeTemplateRepository,
	auditLogRepo repository.AuditLogRepository,
) *MerchantManagementService {
	return &MerchantManagementService{
		merchantRepo:       merchantRepo,
		merchantWalletRepo: merchantWalletRepo,
		feeTemplateRepo:    feeTemplateRepo,
		auditLogRepo:       auditLogRepo,
	}
}

func (s *MerchantManagementService) List(ctx context.Context, req *dtoreq.AdminListMerchantsReq) ([]*dtoresp.AdminMerchantListItem, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	merchants, total, err := s.merchantRepo.ListWithPagination(ctx, req.Page, req.PageSize, req.Status, req.KycStatus, req.Search)
	if err != nil {
		return nil, 0, bizerrors.ErrInternalError
	}

	var items []*dtoresp.AdminMerchantListItem
	for _, m := range merchants {
		items = append(items, &dtoresp.AdminMerchantListItem{
			ID:                m.ID,
			Email:             m.Email,
			Status:            m.Status,
			KycStatus:         m.KycStatus,
			FeeTemplateID:     m.FeeTemplateID,
			KunSubCustomerNo:  m.KunSubCustomerNo,
			AgreementSignedAt: m.AgreementSignedAt,
			KycSubmittedAt:    m.KycSubmittedAt,
			KycCompletedAt:    m.KycCompletedAt,
			FrozenAt:          m.FrozenAt,
			CreatedAt:         m.CreatedAt,
		})
	}

	return items, total, nil
}

func (s *MerchantManagementService) GetDetail(ctx context.Context, id uint64) (*dtoresp.AdminMerchantDetailResp, error) {
	merchant, err := s.merchantRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	wallets, err := s.merchantWalletRepo.FindByMerchantID(ctx, id)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	var walletResps []dtoresp.WalletBalanceResp
	for _, w := range wallets {
		available := w.Balance.Sub(w.FrozenBalance)
		walletResps = append(walletResps, dtoresp.WalletBalanceResp{
			AccountType:      w.AccountType,
			Currency:         w.Currency,
			Balance:          w.Balance.String(),
			FrozenBalance:    w.FrozenBalance.String(),
			AvailableBalance: available.String(),
		})
	}

	resp := &dtoresp.AdminMerchantDetailResp{
		ID:                merchant.ID,
		Email:             merchant.Email,
		Status:            merchant.Status,
		KycStatus:         merchant.KycStatus,
		KycFailReason:     merchant.KycFailReason,
		FeeTemplateID:     merchant.FeeTemplateID,
		KunSubCustomerNo:  merchant.KunSubCustomerNo,
		AgreementSignedAt: merchant.AgreementSignedAt,
		KycSubmittedAt:    merchant.KycSubmittedAt,
		KycCompletedAt:    merchant.KycCompletedAt,
		FrozenAt:          merchant.FrozenAt,
		CreatedAt:         merchant.CreatedAt,
		Wallets:           walletResps,
	}

	if merchant.FeeTemplateID != nil {
		template, err := s.feeTemplateRepo.FindByID(ctx, *merchant.FeeTemplateID)
		if err == nil {
			resp.FeeTemplateName = &template.Name
		}
	}

	return resp, nil
}

func (s *MerchantManagementService) UpdateStatus(ctx context.Context, adminID, merchantID uint64, req *dtoreq.UpdateMerchantStatusReq) error {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrNotFoundError
		}
		return bizerrors.ErrInternalError
	}

	if req.Status == "FROZEN" && merchant.Status != "ACTIVE" {
		return bizerrors.ErrInvalidStatusTransitionE
	}
	if req.Status == "ACTIVE" && merchant.Status != "FROZEN" {
		return bizerrors.ErrInvalidStatusTransitionE
	}

	fields := map[string]interface{}{"status": req.Status}
	if req.Status == "FROZEN" {
		now := time.Now()
		fields["frozen_at"] = &now
	} else {
		fields["frozen_at"] = nil
	}

	if err := s.merchantRepo.UpdateFields(ctx, merchantID, fields); err != nil {
		return bizerrors.ErrInternalError
	}

	s.logAudit(ctx, adminID, "UPDATE_MERCHANT_STATUS", "Merchant", fmt.Sprintf("%d", merchantID), map[string]string{
		"old_status": merchant.Status,
		"new_status": req.Status,
	})

	return nil
}

func (s *MerchantManagementService) UpdateFeeTemplate(ctx context.Context, adminID, merchantID uint64, req *dtoreq.UpdateMerchantFeeTemplateReq) error {
	_, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrNotFoundError
		}
		return bizerrors.ErrInternalError
	}

	_, err = s.feeTemplateRepo.FindByID(ctx, req.FeeTemplateID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrNotFoundError
		}
		return bizerrors.ErrInternalError
	}

	if err := s.merchantRepo.UpdateFields(ctx, merchantID, map[string]interface{}{"fee_template_id": req.FeeTemplateID}); err != nil {
		return bizerrors.ErrInternalError
	}

	s.logAudit(ctx, adminID, "UPDATE_MERCHANT_FEE_TEMPLATE", "Merchant", fmt.Sprintf("%d", merchantID), map[string]uint64{
		"fee_template_id": req.FeeTemplateID,
	})

	return nil
}

func (s *MerchantManagementService) ApproveKyc(ctx context.Context, adminID, merchantID uint64) error {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrNotFoundError
		}
		return bizerrors.ErrInternalError
	}

	if merchant.Status != "PENDING_KYC" {
		return bizerrors.ErrMerchantNotPendingKycE
	}

	now := time.Now()
	fields := map[string]interface{}{
		"status":           "ACTIVE",
		"kyc_status":       "AUTH_SUC",
		"kyc_completed_at": &now,
	}

	if err := s.merchantRepo.UpdateFields(ctx, merchantID, fields); err != nil {
		return bizerrors.ErrInternalError
	}

	s.logAudit(ctx, adminID, "APPROVE_KYC", "Merchant", fmt.Sprintf("%d", merchantID), nil)

	return nil
}

func (s *MerchantManagementService) RejectKyc(ctx context.Context, adminID, merchantID uint64, req *dtoreq.KycRejectReq) error {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrNotFoundError
		}
		return bizerrors.ErrInternalError
	}

	if merchant.Status != "PENDING_KYC" {
		return bizerrors.ErrMerchantNotPendingKycE
	}

	now := time.Now()
	fields := map[string]interface{}{
		"kyc_status":       "AUTH_FAIL",
		"kyc_fail_reason":  req.Reason,
		"kyc_completed_at": &now,
	}

	if err := s.merchantRepo.UpdateFields(ctx, merchantID, fields); err != nil {
		return bizerrors.ErrInternalError
	}

	s.logAudit(ctx, adminID, "REJECT_KYC", "Merchant", fmt.Sprintf("%d", merchantID), map[string]string{
		"reason": req.Reason,
	})

	return nil
}

func (s *MerchantManagementService) logAudit(ctx context.Context, operatorID uint64, action, targetType, targetID string, detail interface{}) {
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
