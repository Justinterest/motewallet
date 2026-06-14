package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"gorm.io/gorm"
	dtoresp "motewallet/internal/dto/response"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/kun"
	kundto "motewallet/internal/pkg/kun/dto"
	"motewallet/internal/pkg/utils"
	"motewallet/internal/repository"
)

const (
	kunDepositAddressListPath = "/rest/v2.0/customer/crypto/deposit/addresses"
	kunDepositHistoryPath     = "/rest/v2.0/trade/digital/wallet/query/recharge"
	// KUN rejects queries spanning more than 90 days.
	kunDepositHistoryMaxDays = 89
)

type DepositService struct {
	kunClient         kun.KUNClient
	merchantRepo      repository.MerchantRepository
	currencyConfigSvc *CurrencyConfigService
}

func NewDepositService(
	kunClient kun.KUNClient,
	merchantRepo repository.MerchantRepository,
	currencyConfigSvc *CurrencyConfigService,
) *DepositService {
	return &DepositService{
		kunClient:         kunClient,
		merchantRepo:      merchantRepo,
		currencyConfigSvc: currencyConfigSvc,
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

	if err := s.currencyConfigSvc.EnsureCurrencySupported(ctx, merchant, currency); err != nil {
		return nil, err
	}

	chainType := utils.KUNDepositChain(currency, chain)

	var addresses []kundto.DepositAddressItem
	err = s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, kunDepositAddressListPath, &kundto.DepositAddressListReq{
		RequestNo: kun.GenerateRequestNo(),
		Currency:  currency,
		ChainType: chainType,
	}, &addresses)
	if err != nil {
		slog.Error("KUN get deposit addresses failed", slog.Any("error", err))
		return nil, bizerrors.ErrKUNAPIFailedE
	}

	if len(addresses) == 0 {
		return nil, bizerrors.ErrNotFoundError
	}

	item := addresses[0]
	for _, addr := range addresses {
		if addr.ChainType == chainType || addr.Chain == chain || addr.Chain == chainType {
			item = addr
			break
		}
	}

	return &dtoresp.DepositAddressResp{
		Address:  item.Address,
		Currency: item.Currency,
		Network:  utils.ChainDisplayName(item.ChainType),
	}, nil
}

func (s *DepositService) ListDepositOrders(ctx context.Context, merchantID uint64, currency, chain string, page, pageSize int) (*dtoresp.DepositOrderListResp, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

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

	if currency != "" {
		if err := s.currencyConfigSvc.EnsureCurrencySupported(ctx, merchant, currency); err != nil {
			return nil, err
		}
	}

	loc := kunDepositLocation()
	now := time.Now().In(loc)
	req := kundto.DepositHistoryQueryReq{
		StartTime: now.AddDate(0, 0, -kunDepositHistoryMaxDays).Format("2006-01-02 15:04:05"),
		EndTime:   now.Format("2006-01-02 15:04:05"),
		PageNo:    page,
		PageSize:  pageSize,
	}
	if currency != "" {
		req.OrderCurrency = currency
	}
	if chain != "" {
		req.Chain = utils.KUNDepositChain(currency, chain)
	}

	var kunResp kundto.DepositHistoryQueryResp
	err = s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, kunDepositHistoryPath, req, &kunResp)
	if err != nil {
		slog.Error("KUN list deposit history failed", slog.Any("error", err))
		return nil, bizerrors.ErrKUNAPIFailedE
	}

	var resp []dtoresp.DepositOrderResp
	for _, record := range kunResp.Items() {
		createdAt := parseKunDateTime(record.OrderTime)
		resp = append(resp, dtoresp.DepositOrderResp{
			ID:        record.OrderID,
			Currency:  record.OrderCurrency,
			Network:   utils.ChainDisplayName(record.Chain),
			Amount:    record.OrderAmount,
			TxHash:    nullableString(record.TxID),
			Status:    mapDepositOrderStatus(record.OrderStatus),
			CreatedAt: createdAt,
		})
	}

	return &dtoresp.DepositOrderListResp{
		Orders: resp,
		Total:  kunResp.Total(),
	}, nil
}
