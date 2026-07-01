package service

import (
	"context"

	dtoreq "motewallet/internal/dto/request"
	dtoresp "motewallet/internal/dto/response"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/repository"
)

type AdminExchangeService struct {
	exchangeOrderRepo repository.ExchangeOrderRepository
	exchangeService   *ExchangeService
}

func NewAdminExchangeService(
	exchangeOrderRepo repository.ExchangeOrderRepository,
	exchangeService *ExchangeService,
) *AdminExchangeService {
	return &AdminExchangeService{
		exchangeOrderRepo: exchangeOrderRepo,
		exchangeService:   exchangeService,
	}
}

func (s *AdminExchangeService) List(ctx context.Context, req *dtoreq.AdminListExchangesReq) (*dtoresp.AdminExchangeListResp, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	rows, total, err := s.exchangeOrderRepo.ListForAdmin(ctx, repository.AdminExchangeListFilter{
		Page:          page,
		PageSize:      pageSize,
		MerchantID:    req.MerchantID,
		Currency:      req.Currency,
		Status:        req.Status,
		MerchantEmail: req.MerchantEmail,
	})
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	exchanges := make([]dtoresp.AdminExchangeListItem, 0, len(rows))
	for _, row := range rows {
		item := dtoresp.AdminExchangeListItem{
			ID:              row.ID,
			PlatformOrderID: row.PlatformOrderID,
			MerchantID:      row.MerchantID,
			MerchantEmail:   row.MerchantEmail,
			ExchangeType:    row.ExchangeType,
			FromCurrency:    row.FromCurrency,
			ToCurrency:      row.ToCurrency,
			FromAmount:      row.FromAmount.String(),
			PlatformFee:     row.PlatformFee.String(),
			Status:          row.TransactionStatus,
			KunOrderID:      row.KunOrderID,
			KunRequestNo:    row.KunRequestNo,
			CreatedAt:       row.CreatedAt,
			CompletedAt:     row.CompletedAt,
		}
		if row.ToAmount != nil {
			item.ToAmount = row.ToAmount.String()
		}
		if row.ExchangeRate != nil {
			item.ExchangeRate = row.ExchangeRate.String()
		}
		if row.FailReason != nil {
			item.FailReason = *row.FailReason
		}
		exchanges = append(exchanges, item)
	}

	return &dtoresp.AdminExchangeListResp{
		Exchanges: exchanges,
		Total:     total,
	}, nil
}

func (s *AdminExchangeService) SyncStatus(ctx context.Context, orderID uint64) (*dtoresp.AdminExchangeSyncResp, error) {
	return s.exchangeService.SyncOrderStatus(ctx, orderID)
}
