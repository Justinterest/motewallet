package service

import (
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"
	dtoresp "motewallet-withdrawal/backend/internal/dto/response"
	bizerrors "motewallet-withdrawal/backend/internal/pkg/errors"
	"motewallet-withdrawal/backend/internal/pkg/kun"
	kundto "motewallet-withdrawal/backend/internal/pkg/kun/dto"
	"motewallet-withdrawal/backend/internal/repository"
)

type DepositService struct {
	kunClient        kun.KUNClient
	merchantRepo     repository.MerchantRepository
	depositOrderRepo repository.DepositOrderRepository
}

func NewDepositService(
	kunClient kun.KUNClient,
	merchantRepo repository.MerchantRepository,
	depositOrderRepo repository.DepositOrderRepository,
) *DepositService {
	return &DepositService{
		kunClient:        kunClient,
		merchantRepo:     merchantRepo,
		depositOrderRepo: depositOrderRepo,
	}
}

func (s *DepositService) GetDepositAddresses(ctx context.Context, merchantID uint64, currency, chain string) (*dtoresp.DepositAddressResp, error) {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	if merchant.KunSubCustomerNo == nil {
		return nil, bizerrors.ErrMerchantNotRegisteredE
	}

	var kunResp kundto.DepositAddressResp
	err = s.kunClient.Post(ctx, "/rest/v2.0/customer/crypto/deposit/addresses", &kundto.DepositAddressReq{
		SubCustomerNo: *merchant.KunSubCustomerNo,
		Currency:      currency,
		Chain:         chain,
	}, &kunResp)
	if err != nil {
		slog.Error("KUN get deposit addresses failed", slog.Any("error", err))
		return nil, bizerrors.ErrKUNAPIFailedE
	}

	return &dtoresp.DepositAddressResp{
		Address:  kunResp.Address,
		Currency: kunResp.Currency,
		Chain:    kunResp.Chain,
	}, nil
}

func (s *DepositService) ListDepositOrders(ctx context.Context, merchantID uint64, page, pageSize int) (*dtoresp.DepositOrderListResp, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	orders, total, err := s.depositOrderRepo.ListByMerchant(ctx, merchantID, page, pageSize)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	var resp []dtoresp.DepositOrderResp
	for _, o := range orders {
		resp = append(resp, dtoresp.DepositOrderResp{
			ID:        o.ID,
			Currency:  o.Currency,
			Chain:     o.Chain,
			Amount:    "0",
			TxID:      o.TxID,
			Status:    "COMPLETED",
			CreatedAt: o.CreatedAt,
		})
	}

	return &dtoresp.DepositOrderListResp{
		Orders: resp,
		Total:  total,
	}, nil
}
