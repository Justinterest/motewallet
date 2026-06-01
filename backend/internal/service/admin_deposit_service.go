package service

import (
	"context"

	dtoreq "motewallet/internal/dto/request"
	dtoresp "motewallet/internal/dto/response"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/utils"
	"motewallet/internal/repository"
)

type AdminDepositService struct {
	depositOrderRepo repository.DepositOrderRepository
}

func NewAdminDepositService(depositOrderRepo repository.DepositOrderRepository) *AdminDepositService {
	return &AdminDepositService{depositOrderRepo: depositOrderRepo}
}

func (s *AdminDepositService) List(ctx context.Context, req *dtoreq.AdminListDepositsReq) (*dtoresp.AdminDepositListResp, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	rows, total, err := s.depositOrderRepo.ListForAdmin(ctx, repository.AdminDepositListFilter{
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

	deposits := make([]dtoresp.AdminDepositListItem, 0, len(rows))
	for _, row := range rows {
		completedAt := row.CompletedAt
		if completedAt == nil {
			completedAt = row.TxnCompletedAt
		}
		deposits = append(deposits, dtoresp.AdminDepositListItem{
			ID:              row.ID,
			PlatformOrderID: row.PlatformOrderID,
			MerchantID:      row.MerchantID,
			MerchantEmail:   row.MerchantEmail,
			Currency:        row.Currency,
			Network:         utils.ChainDisplayName(row.Chain),
			Amount:          row.Amount.String(),
			TxHash:          row.TxID,
			ToAddress:       row.ToAddress,
			FromAddress:     row.FromAddress,
			Status:          row.TransactionStatus,
			CreatedAt:       row.CreatedAt,
			CompletedAt:     completedAt,
		})
	}

	return &dtoresp.AdminDepositListResp{
		Deposits: deposits,
		Total:    total,
	}, nil
}
