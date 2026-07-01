package service

import (
	"context"

	dtoreq "motewallet/internal/dto/request"
	dtoresp "motewallet/internal/dto/response"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/utils"
	"motewallet/internal/repository"
)

type AdminWithdrawalService struct {
	withdrawalOrderRepo repository.WithdrawalOrderRepository
}

func NewAdminWithdrawalService(withdrawalOrderRepo repository.WithdrawalOrderRepository) *AdminWithdrawalService {
	return &AdminWithdrawalService{withdrawalOrderRepo: withdrawalOrderRepo}
}

func (s *AdminWithdrawalService) List(ctx context.Context, req *dtoreq.AdminListWithdrawalsReq) (*dtoresp.AdminWithdrawalListResp, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	rows, total, err := s.withdrawalOrderRepo.ListForAdmin(ctx, repository.AdminWithdrawalListFilter{
		Page:          page,
		PageSize:      pageSize,
		MerchantID:    req.MerchantID,
		Currency:      req.Currency,
		Status:        req.Status,
		ReviewStatus:  req.ReviewStatus,
		Type:          req.Type,
		MerchantEmail: req.MerchantEmail,
	})
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	withdrawals := make([]dtoresp.AdminWithdrawalListItem, 0, len(rows))
	for _, row := range rows {
		network := ""
		if row.Chain != nil {
			network = utils.ChainDisplayName(*row.Chain)
		}
		withdrawals = append(withdrawals, dtoresp.AdminWithdrawalListItem{
			ID:              row.ID,
			PlatformOrderID: row.PlatformOrderID,
			MerchantID:      row.MerchantID,
			MerchantEmail:   row.MerchantEmail,
			Type:            row.WithdrawalType,
			Currency:        row.Currency,
			Network:         network,
			Amount:          row.Amount.String(),
			PlatformFee:     row.PlatformFee.String(),
			Status:          row.TransactionStatus,
			ReviewStatus:    row.ReviewStatus,
			ToAddress:       row.ToAddress,
			TxID:            row.TxID,
			CreatedAt:       row.CreatedAt,
			CompletedAt:     row.CompletedAt,
		})
	}

	return &dtoresp.AdminWithdrawalListResp{
		Withdrawals: withdrawals,
		Total:       total,
	}, nil
}
