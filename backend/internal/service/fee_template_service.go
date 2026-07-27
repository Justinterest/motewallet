package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	dtoreq "motewallet/internal/dto/request"
	dtoresp "motewallet/internal/dto/response"
	"motewallet/internal/model"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/repository"
)

type FeeTemplateService struct {
	db                   *gorm.DB
	feeTemplateRepo      repository.FeeTemplateRepository
	exchangeItemRepo     repository.FeeTemplateExchangeItemRepository
	cryptoWithdrawalRepo repository.FeeTemplateCryptoWithdrawalItemRepository
	fiatWithdrawalRepo   repository.FeeTemplateFiatWithdrawalItemRepository
	auditLogRepo         repository.AuditLogRepository
}

func NewFeeTemplateService(
	db *gorm.DB,
	feeTemplateRepo repository.FeeTemplateRepository,
	exchangeItemRepo repository.FeeTemplateExchangeItemRepository,
	cryptoWithdrawalRepo repository.FeeTemplateCryptoWithdrawalItemRepository,
	fiatWithdrawalRepo repository.FeeTemplateFiatWithdrawalItemRepository,
	auditLogRepo repository.AuditLogRepository,
) *FeeTemplateService {
	return &FeeTemplateService{
		db:                   db,
		feeTemplateRepo:      feeTemplateRepo,
		exchangeItemRepo:     exchangeItemRepo,
		cryptoWithdrawalRepo: cryptoWithdrawalRepo,
		fiatWithdrawalRepo:   fiatWithdrawalRepo,
		auditLogRepo:         auditLogRepo,
	}
}

func (s *FeeTemplateService) Create(ctx context.Context, adminID uint64, req *dtoreq.CreateFeeTemplateReq) (*dtoresp.FeeTemplateDetailResp, error) {
	if err := validateFeeTemplateDeductionMethods(req.ExchangeFeeDeductionMethod, req.CryptoWithdrawalFeeDeductionMethod, req.FiatWithdrawalFeeDeductionMethod); err != nil {
		return nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, err.Error())
	}
	template := &model.FeeTemplate{
		Name:                               req.Name,
		Description:                        req.Description,
		IsDefault:                          req.IsDefault,
		ExchangeFeeDeductionMethod:         normalizeFeeDeductionMethod(req.ExchangeFeeDeductionMethod),
		CryptoWithdrawalFeeDeductionMethod: normalizeFeeDeductionMethod(req.CryptoWithdrawalFeeDeductionMethod),
		FiatWithdrawalFeeDeductionMethod:   normalizeFeeDeductionMethod(req.FiatWithdrawalFeeDeductionMethod),
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if req.IsDefault {
			if err := tx.Model(&model.FeeTemplate{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
				return err
			}
		}

		if err := tx.Create(template).Error; err != nil {
			return err
		}

		exchangeItems, err := s.buildExchangeItems(template.ID, req.ExchangeItems)
		if err != nil {
			return err
		}
		if err := s.exchangeItemRepo.BatchReplace(ctx, tx, template.ID, exchangeItems); err != nil {
			return err
		}

		cryptoItems, err := s.buildCryptoWithdrawalItems(template.ID, req.CryptoWithdrawalItems)
		if err != nil {
			return err
		}
		if err := s.cryptoWithdrawalRepo.BatchReplace(ctx, tx, template.ID, cryptoItems); err != nil {
			return err
		}

		fiatItems, err := s.buildFiatWithdrawalItems(template.ID, req.FiatWithdrawalItems)
		if err != nil {
			return err
		}
		if err := s.fiatWithdrawalRepo.BatchReplace(ctx, tx, template.ID, fiatItems); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	s.logAudit(ctx, adminID, "CREATE_FEE_TEMPLATE", "FeeTemplate", fmt.Sprintf("%d", template.ID), nil)

	return s.GetByID(ctx, template.ID)
}

func (s *FeeTemplateService) List(ctx context.Context) (*dtoresp.FeeTemplateListResp, error) {
	templates, err := s.feeTemplateRepo.FindAll(ctx)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	var items []dtoresp.FeeTemplateListItem
	for _, t := range templates {
		items = append(items, dtoresp.FeeTemplateListItem{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			IsDefault:   t.IsDefault,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}

	return &dtoresp.FeeTemplateListResp{Templates: items}, nil
}

func (s *FeeTemplateService) GetByID(ctx context.Context, id uint64) (*dtoresp.FeeTemplateDetailResp, error) {
	template, err := s.feeTemplateRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	exchangeItems, err := s.exchangeItemRepo.FindByTemplateID(ctx, id)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	cryptoItems, err := s.cryptoWithdrawalRepo.FindByTemplateID(ctx, id)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	fiatItems, err := s.fiatWithdrawalRepo.FindByTemplateID(ctx, id)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	var exchangeResps []dtoresp.ExchangeItemResp
	for _, item := range exchangeItems {
		exchangeResps = append(exchangeResps, dtoresp.ExchangeItemResp{
			ID:             item.ID,
			FromCurrency:   item.FromCurrency,
			ToCurrency:     item.ToCurrency,
			FeeRate:        item.FeeRate.String(),
			MinFee:         item.MinFee.String(),
			MinFeeCurrency: item.MinFeeCurrency,
		})
	}

	var cryptoResps []dtoresp.CryptoWithdrawalItemResp
	for _, item := range cryptoItems {
		cryptoResps = append(cryptoResps, dtoresp.CryptoWithdrawalItemResp{
			ID:       item.ID,
			Currency: item.Currency,
			Chain:    item.Chain,
			FeeRate:  item.FeeRate.String(),
			FixedFee: item.FixedFee.String(),
		})
	}

	var fiatResps []dtoresp.FiatWithdrawalItemResp
	for _, item := range fiatItems {
		fiatResps = append(fiatResps, dtoresp.FiatWithdrawalItemResp{
			ID:           item.ID,
			Currency:     item.Currency,
			TransferType: item.TransferType,
			FeeRate:      item.FeeRate.String(),
			FixedFee:     item.FixedFee.String(),
		})
	}

	return &dtoresp.FeeTemplateDetailResp{
		ID:                                 template.ID,
		Name:                               template.Name,
		Description:                        template.Description,
		IsDefault:                          template.IsDefault,
		ExchangeFeeDeductionMethod:         normalizeFeeDeductionMethod(template.ExchangeFeeDeductionMethod),
		CryptoWithdrawalFeeDeductionMethod: normalizeFeeDeductionMethod(template.CryptoWithdrawalFeeDeductionMethod),
		FiatWithdrawalFeeDeductionMethod:   normalizeFeeDeductionMethod(template.FiatWithdrawalFeeDeductionMethod),
		ExchangeItems:                      exchangeResps,
		CryptoWithdrawalItems:              cryptoResps,
		FiatWithdrawalItems:                fiatResps,
		CreatedAt:                          template.CreatedAt,
		UpdatedAt:                          template.UpdatedAt,
	}, nil
}

func (s *FeeTemplateService) Update(ctx context.Context, adminID, id uint64, req *dtoreq.UpdateFeeTemplateReq) (*dtoresp.FeeTemplateDetailResp, error) {
	if err := validateFeeTemplateDeductionMethods(req.ExchangeFeeDeductionMethod, req.CryptoWithdrawalFeeDeductionMethod, req.FiatWithdrawalFeeDeductionMethod); err != nil {
		return nil, bizerrors.NewBusinessError(400, bizerrors.ErrValidation, err.Error())
	}
	template, err := s.feeTemplateRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	template.Name = req.Name
	template.Description = req.Description
	template.IsDefault = req.IsDefault
	template.ExchangeFeeDeductionMethod = normalizeFeeDeductionMethod(req.ExchangeFeeDeductionMethod)
	template.CryptoWithdrawalFeeDeductionMethod = normalizeFeeDeductionMethod(req.CryptoWithdrawalFeeDeductionMethod)
	template.FiatWithdrawalFeeDeductionMethod = normalizeFeeDeductionMethod(req.FiatWithdrawalFeeDeductionMethod)

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if req.IsDefault {
			if err := tx.Model(&model.FeeTemplate{}).Where("is_default = ? AND id != ?", true, id).Update("is_default", false).Error; err != nil {
				return err
			}
		}

		if err := tx.Save(template).Error; err != nil {
			return err
		}

		exchangeItems, err := s.buildExchangeItems(id, req.ExchangeItems)
		if err != nil {
			return err
		}
		if err := s.exchangeItemRepo.BatchReplace(ctx, tx, id, exchangeItems); err != nil {
			return err
		}

		cryptoItems, err := s.buildCryptoWithdrawalItems(id, req.CryptoWithdrawalItems)
		if err != nil {
			return err
		}
		if err := s.cryptoWithdrawalRepo.BatchReplace(ctx, tx, id, cryptoItems); err != nil {
			return err
		}

		fiatItems, err := s.buildFiatWithdrawalItems(id, req.FiatWithdrawalItems)
		if err != nil {
			return err
		}
		if err := s.fiatWithdrawalRepo.BatchReplace(ctx, tx, id, fiatItems); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	s.logAudit(ctx, adminID, "UPDATE_FEE_TEMPLATE", "FeeTemplate", fmt.Sprintf("%d", id), nil)

	return s.GetByID(ctx, id)
}

func validateFeeTemplateDeductionMethods(methods ...string) error {
	for _, method := range methods {
		if err := validateFeeDeductionMethod(method); err != nil {
			return err
		}
	}
	return nil
}

func (s *FeeTemplateService) Delete(ctx context.Context, adminID, id uint64) error {
	_, err := s.feeTemplateRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrNotFoundError
		}
		return bizerrors.ErrInternalError
	}

	referenced, err := s.feeTemplateRepo.IsReferencedByMerchant(ctx, id)
	if err != nil {
		return bizerrors.ErrInternalError
	}
	if referenced {
		return bizerrors.ErrFeeTemplateInUseE
	}

	if err := s.feeTemplateRepo.Delete(ctx, id); err != nil {
		return bizerrors.ErrInternalError
	}

	s.logAudit(ctx, adminID, "DELETE_FEE_TEMPLATE", "FeeTemplate", fmt.Sprintf("%d", id), nil)

	return nil
}

func (s *FeeTemplateService) buildExchangeItems(templateID uint64, items []dtoreq.ExchangeItemReq) ([]*model.FeeTemplateExchangeItem, error) {
	var result []*model.FeeTemplateExchangeItem
	for _, item := range items {
		feeRate, err := decimal.NewFromString(item.FeeRate)
		if err != nil {
			return nil, err
		}
		minFee, err := decimal.NewFromString(item.MinFee)
		if err != nil {
			return nil, err
		}
		result = append(result, &model.FeeTemplateExchangeItem{
			FeeTemplateID:  templateID,
			FromCurrency:   item.FromCurrency,
			ToCurrency:     item.ToCurrency,
			FeeRate:        feeRate,
			MinFee:         minFee,
			MinFeeCurrency: item.MinFeeCurrency,
		})
	}
	return result, nil
}

func (s *FeeTemplateService) buildCryptoWithdrawalItems(templateID uint64, items []dtoreq.CryptoWithdrawalItemReq) ([]*model.FeeTemplateCryptoWithdrawalItem, error) {
	var result []*model.FeeTemplateCryptoWithdrawalItem
	for _, item := range items {
		feeRate, err := decimal.NewFromString(item.FeeRate)
		if err != nil {
			return nil, err
		}
		fixedFee, err := decimal.NewFromString(item.FixedFee)
		if err != nil {
			return nil, err
		}
		result = append(result, &model.FeeTemplateCryptoWithdrawalItem{
			FeeTemplateID: templateID,
			Currency:      item.Currency,
			Chain:         item.Chain,
			FeeRate:       feeRate,
			FixedFee:      fixedFee,
		})
	}
	return result, nil
}

func (s *FeeTemplateService) buildFiatWithdrawalItems(templateID uint64, items []dtoreq.FiatWithdrawalItemReq) ([]*model.FeeTemplateFiatWithdrawalItem, error) {
	var result []*model.FeeTemplateFiatWithdrawalItem
	for _, item := range items {
		feeRate, err := decimal.NewFromString(item.FeeRate)
		if err != nil {
			return nil, err
		}
		fixedFee, err := decimal.NewFromString(item.FixedFee)
		if err != nil {
			return nil, err
		}
		result = append(result, &model.FeeTemplateFiatWithdrawalItem{
			FeeTemplateID: templateID,
			Currency:      item.Currency,
			TransferType:  item.TransferType,
			FeeRate:       feeRate,
			FixedFee:      fixedFee,
		})
	}
	return result, nil
}

func (s *FeeTemplateService) logAudit(ctx context.Context, operatorID uint64, action, targetType, targetID string, detail interface{}) {
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
