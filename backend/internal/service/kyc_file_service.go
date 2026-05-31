package service

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
	"motewallet/internal/config"
	dtoreq "motewallet/internal/dto/request"
	dtoresp "motewallet/internal/dto/response"
	"motewallet/internal/model"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/kun"
	kundto "motewallet/internal/pkg/kun/dto"
	"motewallet/internal/pkg/storage"
	"motewallet/internal/repository"
)

const maxKycFileBytes = 10 << 20 // 10 MB per KUN docs

var allowedKycContentTypes = map[string]bool{
	"application/pdf": true,
	"image/jpeg":      true,
	"image/jpg":       true,
	"image/png":       true,
	"image/gif":       true,
	"image/bmp":       true,
	"image/webp":      true,
}

type KycFileService struct {
	cfg          *config.Config
	merchantRepo repository.MerchantRepository
	storage      *storage.S3Storage
	kunClient    kun.KUNClient
}

func NewKycFileService(
	cfg *config.Config,
	merchantRepo repository.MerchantRepository,
	s3Storage *storage.S3Storage,
	kunClient kun.KUNClient,
) *KycFileService {
	return &KycFileService{
		cfg:          cfg,
		merchantRepo: merchantRepo,
		storage:      s3Storage,
		kunClient:    kunClient,
	}
}

func (s *KycFileService) PresignUpload(ctx context.Context, merchantID uint64, req *dtoreq.PresignKycFileReq) (*dtoresp.PresignKycFileResp, error) {
	if !s.storage.Enabled() {
		return nil, bizerrors.ErrStorageNotConfiguredE
	}

	if err := validateKycFileMeta(req.Filename, req.ContentType, 0); err != nil {
		return nil, err
	}

	objectKey := s.storage.ObjectKey(merchantID, req.Filename)
	uploadURL, err := s.storage.PresignPut(ctx, objectKey, req.ContentType, s.cfg.S3.PresignExpiry)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	return &dtoresp.PresignKycFileResp{
		UploadURL: uploadURL,
		ObjectKey: objectKey,
		ExpiresIn: int(s.cfg.S3.PresignExpiry.Seconds()),
	}, nil
}

func (s *KycFileService) PresignAccess(
	ctx context.Context,
	merchantID uint64,
	req *dtoreq.PresignKycFileAccessReq,
) (*dtoresp.PresignKycFileAccessResp, error) {
	if !s.storage.Enabled() {
		return nil, bizerrors.ErrStorageNotConfiguredE
	}

	if !s.storage.ValidateObjectKey(merchantID, req.ObjectKey) {
		return nil, bizerrors.ErrValidationError
	}

	accessURL, err := s.storage.PresignGet(ctx, req.ObjectKey, s.cfg.S3.PresignExpiry)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	return &dtoresp.PresignKycFileAccessResp{
		AccessURL: accessURL,
		ExpiresIn: int(s.cfg.S3.PresignExpiry.Seconds()),
	}, nil
}

// ResolveSubmitKycFiles uploads staged S3 objects to KUN and replaces object keys with KUN file URLs
// in the submit payload. Already-resolved http(s) paths are left unchanged.
func (s *KycFileService) ResolveSubmitKycFiles(
	ctx context.Context,
	merchantID uint64,
	req *kundto.SubMerchantRegisterReq,
) error {
	if !s.storage.Enabled() {
		return bizerrors.ErrStorageNotConfiguredE
	}

	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrNotFoundError
		}
		return bizerrors.ErrInternalError
	}

	if merchant.KunSubCustomerNo == nil {
		return bizerrors.ErrMerchantNotRegisteredE
	}

	cache := make(map[string]string)

	resolveRefs := func(refs []kundto.FileRef) error {
		for i := range refs {
			staged := strings.TrimSpace(refs[i].Path)
			if staged == "" {
				continue
			}
			if isHTTPURL(staged) {
				continue
			}
			if !s.storage.ValidateObjectKey(merchantID, staged) {
				return bizerrors.NewBusinessError(
					bizerrors.ErrValidationError.HTTPStatus,
					bizerrors.ErrValidation,
					"invalid staged file reference",
				)
			}
			if kunPath, ok := cache[staged]; ok {
				refs[i].Path = kunPath
				continue
			}
			kunPath, err := s.uploadObjectKeyToKun(ctx, merchant, staged)
			if err != nil {
				return err
			}
			cache[staged] = kunPath
			refs[i].Path = kunPath
		}
		return nil
	}

	e := &req.EnterpriseInfo
	for _, refs := range [][]kundto.FileRef{
		e.IncorporationCertificate,
		e.BusinessRegistration,
		e.Incumbency,
		e.AssociationRules,
		e.AuthenticMaterials,
		e.ManagerIdHolding,
		e.AuthorizationLetter,
		e.EquityStructure,
		e.Nnc1,
	} {
		if err := resolveRefs(refs); err != nil {
			return err
		}
	}

	for i := range req.ShareholdersInfo {
		if err := resolveRefs(req.ShareholdersInfo[i].IdHolding); err != nil {
			return err
		}
	}
	for i := range req.DirectorInfo {
		if err := resolveRefs(req.DirectorInfo[i].IdHolding); err != nil {
			return err
		}
	}

	return nil
}

func (s *KycFileService) uploadObjectKeyToKun(
	ctx context.Context,
	merchant *model.Merchant,
	objectKey string,
) (string, error) {
	content, contentType, err := s.storage.GetObject(ctx, objectKey)
	if err != nil {
		return "", bizerrors.ErrInternalError
	}

	if err := validateKycFileMeta(filepath.Base(objectKey), contentType, len(content)); err != nil {
		return "", err
	}

	filename := filepath.Base(objectKey)
	uploadResp, err := s.kunClient.UploadFileAsCustomer(
		ctx,
		*merchant.KunSubCustomerNo,
		filename,
		content,
		contentType,
	)
	if err != nil {
		return "", bizerrors.ErrKUNAPIFailedE
	}

	return uploadResp.URL, nil
}

func isHTTPURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func validateKycFileMeta(filename, contentType string, size int) error {
	ct := strings.ToLower(strings.TrimSpace(contentType))
	if ct == "" || !allowedKycContentTypes[ct] {
		return bizerrors.NewBusinessError(
			bizerrors.ErrValidationError.HTTPStatus,
			bizerrors.ErrValidation,
			"unsupported file type; allowed: PDF, JPEG, PNG, GIF, BMP, WEBP",
		)
	}
	if size > maxKycFileBytes {
		return bizerrors.NewBusinessError(
			bizerrors.ErrValidationError.HTTPStatus,
			bizerrors.ErrValidation,
			"file exceeds 10 MB limit",
		)
	}
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf", ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
	default:
		return bizerrors.NewBusinessError(
			bizerrors.ErrValidationError.HTTPStatus,
			bizerrors.ErrValidation,
			"unsupported file extension",
		)
	}
	return nil
}
