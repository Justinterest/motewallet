package service

import (
	"context"

	dtoreq "motewallet/internal/dto/request"
	dtoresp "motewallet/internal/dto/response"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/repository"
)

type AdminTransferService struct {
	transferOrderRepo repository.TransferOrderRepository
	transferService   *TransferService
}

func NewAdminTransferService(
	transferOrderRepo repository.TransferOrderRepository,
	transferService *TransferService,
) *AdminTransferService {
	return &AdminTransferService{
		transferOrderRepo: transferOrderRepo,
		transferService:   transferService,
	}
}

func (s *AdminTransferService) List(ctx context.Context, req *dtoreq.AdminListTransfersReq) (*dtoresp.AdminTransferListResp, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	rows, total, err := s.transferOrderRepo.ListForAdmin(ctx, repository.AdminTransferListFilter{
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

	transfers := make([]dtoresp.AdminTransferListItem, 0, len(rows))
	for _, row := range rows {
		transfers = append(transfers, dtoresp.AdminTransferListItem{
			ID:              row.ID,
			PlatformOrderID: row.PlatformOrderID,
			MerchantID:      row.MerchantID,
			MerchantEmail:   row.MerchantEmail,
			FromAccountType: row.FromAccountType,
			ToAccountType:   row.ToAccountType,
			Currency:        row.Currency,
			Amount:          row.Amount.String(),
			Status:          row.TransactionStatus,
			KunOrderID:      row.KunOrderID,
			KunRequestNo:    row.KunRequestNo,
			CreatedAt:       row.CreatedAt,
		})
	}

	return &dtoresp.AdminTransferListResp{
		Transfers: transfers,
		Total:     total,
	}, nil
}

func (s *AdminTransferService) SyncStatus(ctx context.Context, orderID uint64) (*dtoresp.AdminTransferSyncResp, error) {
	return s.transferService.SyncOrderStatus(ctx, orderID)
}
