package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"motewallet/internal/config"
	dtoreq "motewallet/internal/dto/request"
	dtoresp "motewallet/internal/dto/response"
	"motewallet/internal/model"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/kun"
	kundto "motewallet/internal/pkg/kun/dto"
	"motewallet/internal/repository"

	"gorm.io/gorm"
)

type OnboardingService struct {
	cfg               *config.Config
	merchantRepo      repository.MerchantRepository
	kycSubmissionRepo repository.MerchantKycSubmissionRepository
	walletSvc         *WalletService
	kycFileSvc        *KycFileService
	kunClient         kun.KUNClient
}

func NewOnboardingService(
	cfg *config.Config,
	merchantRepo repository.MerchantRepository,
	kycSubmissionRepo repository.MerchantKycSubmissionRepository,
	walletSvc *WalletService,
	kycFileSvc *KycFileService,
	kunClient kun.KUNClient,
) *OnboardingService {
	return &OnboardingService{
		cfg:               cfg,
		merchantRepo:      merchantRepo,
		kycSubmissionRepo: kycSubmissionRepo,
		walletSvc:         walletSvc,
		kycFileSvc:        kycFileSvc,
		kunClient:         kunClient,
	}
}

const (
	kunAgreeListPath               = "/rest/v2.0/customer/agreeList"
	kunAgreeAuthPath               = "/rest/v2.0/customer/agree/auth"
	kunAgreementBizCode            = "KUN_SPACE_REGIST"
	kunAgreementSignStatusUnsigned = "UNSIGN"
)

func (s *OnboardingService) queryPendingAgreements(ctx context.Context, subCustomerNo string) (kundto.AgreementList, error) {
	var agreements kundto.AgreementList
	err := s.kunClient.PostAsCustomer(ctx, subCustomerNo, kunAgreeListPath, &kundto.AgreeListReq{
		RequestNo: kun.GenerateRequestNo(),
		// SubCustomerNo: subCustomerNo,
		SignStatus: kunAgreementSignStatusUnsigned,
		BizCode:    kunAgreementBizCode,
	}, &agreements)
	if err != nil {
		return nil, err
	}
	return agreements, nil
}

func mapKunAgreements(items kundto.AgreementList) []dtoresp.AgreementResp {
	agreements := make([]dtoresp.AgreementResp, 0, len(items))
	for _, a := range items {
		if a.ProtocolId == "" {
			continue
		}
		agreements = append(agreements, dtoresp.AgreementResp{
			ID:         a.ProtocolId,
			ProtocolID: a.ProtocolId,
			Title:      a.Title,
			Version:    a.Version,
			URL:        a.URL,
			Required:   true,
		})
	}
	return agreements
}

func (s *OnboardingService) completeAgreementSigning(ctx context.Context, merchant *model.Merchant) error {
	if merchant.AgreementSignedAt != nil && merchant.Status != "PENDING_AGREEMENT" {
		return nil
	}

	now := time.Now()
	merchant.AgreementSignedAt = &now
	if merchant.Status == "PENDING_AGREEMENT" {
		merchant.Status = "PENDING_KYC"
	}

	if err := s.merchantRepo.Update(ctx, merchant); err != nil {
		return bizerrors.ErrInternalError
	}
	return nil
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

	if merchant.KunSubCustomerNo == nil || *merchant.KunSubCustomerNo == "" {
		return nil, bizerrors.ErrMerchantNotRegisteredE
	}

	agreements, err := s.queryPendingAgreements(ctx, *merchant.KunSubCustomerNo)
	if err != nil {
		slog.Error("KUN pending agreement query failed",
			slog.Uint64("merchant_id", merchantID),
			slog.String("sub_customer_no", *merchant.KunSubCustomerNo),
			slog.Any("error", err),
		)
		return nil, bizerrors.ErrKUNAPIFailedE
	}

	if len(agreements) == 0 {
		if err := s.completeAgreementSigning(ctx, merchant); err != nil {
			return nil, err
		}
		return &dtoresp.AgreementListResp{Agreements: nil, Signed: true}, nil
	}

	return &dtoresp.AgreementListResp{
		Agreements: mapKunAgreements(agreements),
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

	if merchant.KunSubCustomerNo == nil || *merchant.KunSubCustomerNo == "" {
		return bizerrors.ErrMerchantNotRegisteredE
	}

	agreements, err := s.queryPendingAgreements(ctx, *merchant.KunSubCustomerNo)
	if err != nil {
		slog.Error("KUN pending agreement query failed before sign",
			slog.Uint64("merchant_id", merchantID),
			slog.Any("error", err),
		)
		return bizerrors.ErrKUNAPIFailedE
	}
	if len(agreements) == 0 {
		return s.completeAgreementSigning(ctx, merchant)
	}

	ids := make([]string, 0, len(agreements))
	for _, a := range agreements {
		if a.ProtocolId != "" {
			ids = append(ids, a.ProtocolId)
		}
	}
	if len(ids) == 0 {
		return bizerrors.NewBusinessError(http.StatusBadRequest, bizerrors.ErrValidation, "no valid agreement protocol ids")
	}

	if err := s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, kunAgreeAuthPath, &kundto.AgreeAuthReq{
		RequestNo: kun.GenerateRequestNo(),
		// SubCustomerNo: *merchant.KunSubCustomerNo,
		ProtocolIds: strings.Join(ids, ","),
	}, nil); err != nil {
		slog.Error("KUN agreement sign failed",
			slog.Uint64("merchant_id", merchantID),
			slog.Any("error", err),
		)
		return bizerrors.ErrKUNAPIFailedE
	}

	return s.completeAgreementSigning(ctx, merchant)
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

	// fmt.Println("ResolveSubmitKycFiles", req)

	if err := s.kycFileSvc.ResolveSubmitKycFiles(ctx, merchant.ID, req); err != nil {
		return err
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
		if kunErr, ok := kun.IsKUNError(err); ok && len(kunErr.Errors) > 0 {
			return bizerrors.NewKycValidationError(kunErr.Errors)
		}
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

const (
	kunCountriesPath            = "/rest/v2.0/customer/fiat/withdrawal/countries"
	kunAuthTypesPath            = "/rest/v2.0/customer/country/auth/types"
	kunSubMerchantAuthQueryPath = "/rest/v2.0/customer/merchant/register/query"
	kycCountrySceneAddress      = "REGISTER_ADDRESS"
	kycCountrySceneRegister     = "REGISTER"
	kycCountrySceneWithdrawal   = "WITHDRAWAL"
	kycCountryLanguage          = "ZH_CN"
)

var allowedKycCountryScenes = map[string]bool{
	kycCountrySceneAddress:    true,
	kycCountrySceneRegister:   true,
	kycCountrySceneWithdrawal: true,
}

// ListKycCountries returns countries/regions for KYC fields.
// REGISTER_ADDRESS: address-related fields; REGISTER: nationality and related fields; WITHDRAWAL: fiat withdrawal.
func (s *OnboardingService) ListKycCountries(ctx context.Context, scene, language string) (*dtoresp.CountryListResp, error) {
	if scene == "" {
		scene = kycCountrySceneAddress
	}
	if !allowedKycCountryScenes[scene] {
		return nil, bizerrors.NewBusinessError(
			http.StatusBadRequest,
			bizerrors.ErrValidation,
			"invalid scene; allowed: REGISTER_ADDRESS, REGISTER, WITHDRAWAL",
		)
	}
	if language == "" {
		language = kycCountryLanguage
	}

	var items []kundto.CountryItem
	if err := s.kunClient.Post(ctx, kunCountriesPath, &kundto.CountriesReq{
		RequestNo: kun.GenerateRequestNo(),
		Scene:     scene,
		Language:  language,
	}, &items); err != nil {
		slog.Error("KUN list countries failed", slog.Any("error", err), slog.String("scene", scene))
		return nil, bizerrors.ErrInternalError
	}

	resp := &dtoresp.CountryListResp{Items: make([]dtoresp.CountryOption, 0, len(items))}
	for _, item := range items {
		if item.CountryCode == "" {
			continue
		}
		label := item.CountryName
		if label == "" {
			label = item.CountryCode
		}
		resp.Items = append(resp.Items, dtoresp.CountryOption{
			CountryCode: item.CountryCode,
			CountryName: label,
		})
	}
	return resp, nil
}

// ListKycCountryAuthTypes returns document types for the given ISO country code.
func (s *OnboardingService) ListKycCountryAuthTypes(ctx context.Context, countryCode string) (*dtoresp.CountryAuthTypesResp, error) {
	countryCode = strings.TrimSpace(strings.ToUpper(countryCode))
	if countryCode == "" {
		return nil, bizerrors.NewBusinessError(http.StatusBadRequest, bizerrors.ErrValidation, "country_code is required")
	}

	var items []kundto.CountryAuthTypeItem
	if err := s.kunClient.Post(ctx, kunAuthTypesPath, &kundto.CountryAuthTypesReq{
		CountryCode: countryCode,
		RequestNo:   kun.GenerateRequestNo(),
	}, &items); err != nil {
		slog.Error("KUN list country auth types failed", slog.Any("error", err), slog.String("countryCode", countryCode))
		return nil, bizerrors.ErrInternalError
	}

	resp := &dtoresp.CountryAuthTypesResp{Items: make([]dtoresp.AuthTypeOption, 0, len(items))}
	for _, item := range items {
		if item.DocCode == "" {
			continue
		}
		label := item.DocName
		if label == "" {
			label = item.DocCode
		}
		resp.Items = append(resp.Items, dtoresp.AuthTypeOption{
			DocCode: item.DocCode,
			DocName: label,
		})
	}
	return resp, nil
}

func (s *OnboardingService) GetKycStatus(ctx context.Context, merchantID uint64) (*dtoresp.KycStatusResp, error) {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	if merchant.KycStatus == "AUTHING" && merchant.KycAuthID != nil && *merchant.KycAuthID != "" {
		if syncErr := s.syncKycStatusFromKun(ctx, merchant); syncErr != nil {
			slog.Error("KUN sub-merchant authentication result query failed", slog.Any("error", syncErr), slog.Uint64("merchant_id", merchantID))
		} else {
			merchant, err = s.merchantRepo.FindByID(ctx, merchantID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, bizerrors.ErrNotFoundError
				}
				return nil, bizerrors.ErrInternalError
			}
		}
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

	return s.syncKycStatusFromKun(ctx, merchant)
}

// syncKycStatusFromKun queries KUN sub-merchant authentication result and updates local merchant state.
// See: https://opendocs.kun.global/docs/api/sub-merchant-authentication-result-query
func (s *OnboardingService) syncKycStatusFromKun(ctx context.Context, merchant *model.Merchant) error {
	if merchant.KunSubCustomerNo == nil {
		return bizerrors.ErrMerchantNotRegisteredE
	}

	if merchant.KycStatus != "AUTHING" {
		return nil
	}

	if merchant.KycAuthID == nil || *merchant.KycAuthID == "" {
		return nil
	}

	var queryResp kundto.MerchantRegisterQueryResp
	if err := s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, kunSubMerchantAuthQueryPath, &kundto.MerchantRegisterQueryReq{
		RequestNo: kun.GenerateRequestNo(),
		AuthID:    *merchant.KycAuthID,
	}, &queryResp); err != nil {
		return err
	}

	now := time.Now()
	switch queryResp.AuthStatus {
	case "AUTH_SUC":
		if err := s.merchantRepo.UpdateFields(ctx, merchant.ID, map[string]interface{}{
			"status":           "ACTIVE",
			"kyc_status":       "AUTH_SUC",
			"kyc_completed_at": now,
		}); err != nil {
			return err
		}
		return s.walletSvc.InitializeWallets(ctx, merchant.ID)
	case "AUTH_FAIL":
		fields := map[string]interface{}{
			"status":           "KYC_FAILED",
			"kyc_status":       "AUTH_FAIL",
			"kyc_completed_at": now,
		}
		if queryResp.FailReason != "" {
			fields["kyc_fail_reason"] = queryResp.FailReason
		}
		return s.merchantRepo.UpdateFields(ctx, merchant.ID, fields)
	}

	return nil
}
