package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
	"motewallet/internal/config"
	dtoreq "motewallet/internal/dto/request"
	dtoresp "motewallet/internal/dto/response"
	"motewallet/internal/model"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/kun"
	kundto "motewallet/internal/pkg/kun/dto"
	"motewallet/internal/repository"
)

type OnboardingService struct {
	cfg              *config.Config
	merchantRepo     repository.MerchantRepository
	kycSubmissionRepo repository.MerchantKycSubmissionRepository
	walletSvc        *WalletService
	kunClient        kun.KUNClient
}

func NewOnboardingService(
	cfg *config.Config,
	merchantRepo repository.MerchantRepository,
	kycSubmissionRepo repository.MerchantKycSubmissionRepository,
	walletSvc *WalletService,
	kunClient kun.KUNClient,
) *OnboardingService {
	return &OnboardingService{
		cfg:              cfg,
		merchantRepo:     merchantRepo,
		kycSubmissionRepo: kycSubmissionRepo,
		walletSvc:        walletSvc,
		kunClient:        kunClient,
	}
}

func (s *OnboardingService) GetAgreements(ctx context.Context, merchantID uint64) (*dtoresp.AgreementListResp, error) {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	if merchant.AgreementSignedAt != nil {
		return &dtoresp.AgreementListResp{Agreements: nil, Signed: true}, nil
	}

	if merchant.KunSubCustomerNo != nil {
		var agreeListResp kundto.AgreeListResp
		err = s.kunClient.Post(ctx, "/rest/v2.0/customer/agreeList", &kundto.AgreeListReq{
			SubCustomerNo: *merchant.KunSubCustomerNo,
			SignStatus:    "UNSIGN",
			BizCode:       "KUN_SPACE_REGIST",
		}, &agreeListResp)
		if err != nil {
			slog.Error("failed to get KUN agreements, returning defaults", slog.Any("error", err))
			return &dtoresp.AgreementListResp{
				Agreements: []dtoresp.AgreementResp{
					{ID: 1, Title: "Service Agreement", Content: "Motewallet platform service agreement.", Required: true},
					{ID: 2, Title: "Privacy Policy", Content: "Privacy and data protection policy.", Required: true},
					{ID: 3, Title: "AML/KYC Policy", Content: "Anti-money laundering and KYC compliance policy.", Required: true},
				},
				Signed: false,
			}, nil
		}

		if len(agreeListResp.List) == 0 {
			return &dtoresp.AgreementListResp{Agreements: nil, Signed: true}, nil
		}

		agreements := make([]dtoresp.AgreementResp, len(agreeListResp.List))
		for i, a := range agreeListResp.List {
			agreements[i] = dtoresp.AgreementResp{
				ID:       i + 1,
				Title:    a.Title,
				Content:  a.URL,
				Required: true,
			}
		}

		return &dtoresp.AgreementListResp{
			Agreements: agreements,
			Signed:     false,
		}, nil
	}

	return &dtoresp.AgreementListResp{
		Agreements: []dtoresp.AgreementResp{
			{ID: 1, Title: "Service Agreement", Content: "Motewallet platform service agreement.", Required: true},
			{ID: 2, Title: "Privacy Policy", Content: "Privacy and data protection policy.", Required: true},
			{ID: 3, Title: "AML/KYC Policy", Content: "Anti-money laundering and KYC compliance policy.", Required: true},
		},
		Signed: false,
	}, nil
}

func (s *OnboardingService) SignAgreements(ctx context.Context, merchantID uint64) error {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrNotFoundError
		}
		return bizerrors.ErrInternalError
	}

	if merchant.Status != "PENDING_AGREEMENT" {
		return bizerrors.ErrMerchantNotPendingAgreementE
	}

	if merchant.KunSubCustomerNo != nil {
		var agreeListResp kundto.AgreeListResp
		err = s.kunClient.Post(ctx, "/rest/v2.0/customer/agreeList", &kundto.AgreeListReq{
			SubCustomerNo: *merchant.KunSubCustomerNo,
			SignStatus:    "UNSIGN",
			BizCode:       "KUN_SPACE_REGIST",
		}, &agreeListResp)
		if err == nil && len(agreeListResp.List) > 0 {
			ids := make([]string, len(agreeListResp.List))
			for i, a := range agreeListResp.List {
				ids[i] = a.ProtocolId
			}
			_ = s.kunClient.Post(ctx, "/rest/v2.0/customer/agree/auth", &kundto.AgreeAuthReq{
				SubCustomerNo: *merchant.KunSubCustomerNo,
				ProtocolIds:   strings.Join(ids, ","),
			}, nil)
		}
	}

	now := time.Now()
	merchant.AgreementSignedAt = &now
	merchant.Status = "PENDING_KYC"

	if err := s.merchantRepo.Update(ctx, merchant); err != nil {
		return bizerrors.ErrInternalError
	}

	return nil
}

func (s *OnboardingService) SubmitKyc(ctx context.Context, merchantID uint64, req *dtoreq.SubmitKycReq) error {
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

	if merchant.KycStatus == "AUTHING" {
		return bizerrors.ErrKycAlreadyProcessingE
	}

	if merchant.KunSubCustomerNo == nil {
		return bizerrors.ErrMerchantNotRegisteredE
	}

	if req.RequestNo == "" {
		req.RequestNo = kun.GenerateRequestNo()
	}

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return bizerrors.ErrInternalError
	}

	submission := &model.MerchantKycSubmission{
		MerchantID:   merchant.ID,
		KunRequestNo: req.RequestNo,
		Payload:      payloadBytes,
		Status:       "SUBMITTED",
	}
	if err := s.kycSubmissionRepo.Create(ctx, submission); err != nil {
		slog.Error("failed to save local KYC submission", slog.Any("error", err))
		return bizerrors.ErrInternalError
	}

	var registerResp kundto.SubMerchantRegisterResp
	if err := s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, "/rest/v2.0/customer/subMerchant/register", req, &registerResp); err != nil {
		slog.Error("KUN sub-merchant onboarding authentication failed", slog.Any("error", err))
		return bizerrors.ErrKUNAPIFailedE
	}

	now := time.Now()
	updates := map[string]interface{}{
		"kyc_status":       "AUTHING",
		"kyc_submitted_at": now,
		"status":           "KYC_REVIEWING",
	}
	if registerResp.AuthID != "" {
		updates["kyc_auth_id"] = registerResp.AuthID
	}

	if err := s.merchantRepo.UpdateFields(ctx, merchant.ID, updates); err != nil {
		return bizerrors.ErrInternalError
	}

	submissionUpdates := map[string]interface{}{"status": "AUTHING"}
	if registerResp.AuthID != "" {
		submissionUpdates["kun_auth_id"] = registerResp.AuthID
	}
	if err := s.kycSubmissionRepo.UpdateFields(ctx, submission.ID, submissionUpdates); err != nil {
		slog.Error("failed to update local KYC submission after KUN success", slog.Any("error", err))
	}

	return nil
}

func (s *OnboardingService) GetKycStatus(ctx context.Context, merchantID uint64) (*dtoresp.KycStatusResp, error) {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	return &dtoresp.KycStatusResp{
		Status:            merchant.Status,
		KycStatus:         merchant.KycStatus,
		KycFailReason:     merchant.KycFailReason,
		KycSubmittedAt:    merchant.KycSubmittedAt,
		KycCompletedAt:    merchant.KycCompletedAt,
		AgreementSignedAt: merchant.AgreementSignedAt,
	}, nil
}

func (s *OnboardingService) PollKycStatus(ctx context.Context, merchantID uint64) error {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		return err
	}

	if merchant.KunSubCustomerNo == nil {
		return bizerrors.ErrMerchantNotRegisteredE
	}

	if merchant.KycStatus != "AUTHING" {
		return nil
	}

	var queryResp kundto.MerchantRegisterQueryResp
	if err := s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, "/rest/v2.0/customer/merchant/register/query", &kundto.MerchantRegisterQueryReq{
		SubCustomerNo: *merchant.KunSubCustomerNo,
	}, &queryResp); err != nil {
		return err
	}

	now := time.Now()
	switch queryResp.AuthStatus {
	case "AUTH_SUC":
		err = s.merchantRepo.UpdateFields(ctx, merchant.ID, map[string]interface{}{
			"status":           "ACTIVE",
			"kyc_status":       "AUTH_SUC",
			"kyc_completed_at": now,
		})
		if err != nil {
			return err
		}
		return s.walletSvc.InitializeWallets(ctx, merchant.ID)
	case "AUTH_FAIL":
		failReason := queryResp.FailReason
		fields := map[string]interface{}{
			"status":           "KYC_FAILED",
			"kyc_status":       "AUTH_FAIL",
			"kyc_completed_at": now,
		}
		if failReason != "" {
			fields["kyc_fail_reason"] = failReason
		}
		return s.merchantRepo.UpdateFields(ctx, merchant.ID, fields)
	}

	return nil
}
