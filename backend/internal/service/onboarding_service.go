package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
	"motewallet/internal/config"
	dtoreq "motewallet/internal/dto/request"
	dtoresp "motewallet/internal/dto/response"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/kun"
	kundto "motewallet/internal/pkg/kun/dto"
	"motewallet/internal/repository"
)

type OnboardingService struct {
	cfg          *config.Config
	merchantRepo repository.MerchantRepository
	walletSvc    *WalletService
	kunClient    kun.KUNClient
}

func NewOnboardingService(cfg *config.Config, merchantRepo repository.MerchantRepository, walletSvc *WalletService, kunClient kun.KUNClient) *OnboardingService {
	return &OnboardingService{
		cfg:          cfg,
		merchantRepo: merchantRepo,
		walletSvc:    walletSvc,
		kunClient:    kunClient,
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

	if merchant.KunSubCustomerNo == nil {
		return &dtoresp.AgreementListResp{
			Agreements: []dtoresp.AgreementResp{
				{ID: 1, Title: "Service Agreement", Content: "Motewallet platform service agreement.", Required: true},
				{ID: 2, Title: "Privacy Policy", Content: "Privacy and data protection policy.", Required: true},
				{ID: 3, Title: "AML/KYC Policy", Content: "Anti-money laundering and KYC compliance policy.", Required: true},
			},
			Signed: false,
		}, nil
	}

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

	if merchant.KycStatus == "PROCESSING" {
		return bizerrors.ErrKycAlreadyProcessingE
	}

	now := time.Now()

	if merchant.KunSubCustomerNo != nil {
		kunReq := &kundto.SubMerchantRegisterReq{
			SubCustomerNo:      *merchant.KunSubCustomerNo,
			CompanyName:        req.CompanyName,
			Country:            req.Country,
			RegistrationNumber: req.RegistrationNumber,
			ContactName:        req.ContactName,
			ContactPhone:       req.ContactPhone,
			RequestNo:          kun.GenerateRequestNo(),
		}
		if err := s.kunClient.Post(ctx, "/rest/v2.0/customer/subMerchant/register", kunReq, nil); err != nil {
			slog.Error("KUN KYC submission failed", slog.Any("error", err))
			return bizerrors.ErrKUNAPIFailedE
		}
	}

	return s.merchantRepo.UpdateFields(ctx, merchant.ID, map[string]interface{}{
		"kyc_status":     "PROCESSING",
		"kyc_submitted_at": now,
	})
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

	if merchant.KycStatus != "PROCESSING" {
		return nil
	}

	var queryResp kundto.MerchantRegisterQueryResp
	if err := s.kunClient.Post(ctx, "/rest/v2.0/customer/merchant/register/query", &kundto.MerchantRegisterQueryReq{
		SubCustomerNo: *merchant.KunSubCustomerNo,
	}, &queryResp); err != nil {
		return err
	}

	now := time.Now()
	switch queryResp.AuthStatus {
	case "AUTH_SUC":
		err = s.merchantRepo.UpdateFields(ctx, merchant.ID, map[string]interface{}{
			"status":          "ACTIVE",
			"kyc_status":      "AUTH_SUC",
			"kyc_completed_at": now,
		})
		if err != nil {
			return err
		}
		return s.walletSvc.InitializeWallets(ctx, merchant.ID)
	case "AUTH_FAIL":
		return s.merchantRepo.UpdateFields(ctx, merchant.ID, map[string]interface{}{
			"kyc_status":      "AUTH_FAIL",
			"kyc_completed_at": now,
		})
	}

	return nil
}
