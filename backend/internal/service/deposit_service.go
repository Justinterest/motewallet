package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
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

	chainType := normalizeDepositChainType(currency, chain)

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

	now := time.Now()
	req := kundto.DepositHistoryQueryReq{
		StartTime: now.AddDate(0, -3, 0).Format("2006-01-02 15:04:05"),
		EndTime:   now.Format("2006-01-02 15:04:05"),
		PageNo:    page,
		PageSize:  pageSize,
	}
	if currency != "" {
		req.OrderCurrency = currency
	}
	if chain != "" {
		req.Chain = normalizeDepositChainType(currency, chain)
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

func normalizeDepositChainType(currency, chain string) string {
	switch strings.ToUpper(chain) {
	case "TRC20", "TRX_TRC20":
		return "TRX_TRC20"
	case "ERC20", "ETH_ERC20":
		return "ETH_ERC20"
	case "BTC", "BTC_BITCOIN":
		return "BTC_Bitcoin"
	case "SOL", "SOL_SOLANA":
		return "SOL_Solana"
	case "BSC_BEP20", "BEP20", "BNB_BEP20":
		return "BSC_BEP20"
	case "TON":
		return "TON"
	default:
		if chain != "" {
			return chain
		}
		if strings.ToUpper(currency) == "BTC" {
			return "BTC_Bitcoin"
		}
		return "TRX_TRC20"
	}
}

func mapDepositOrderStatus(status string) string {
	switch strings.ToUpper(status) {
	case "SUCCESS":
		return "COMPLETED"
	case "FAIL", "FAILED":
		return "FAILED"
	case "PROCESSING", "PENDING":
		return "PROCESSING"
	default:
		return status
	}
}

func parseKunDateTime(value string) time.Time {
	if value == "" {
		return time.Now()
	}
	for _, layout := range []string{"2006-01-02 15:04:05", time.RFC3339} {
		if t, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return t
		}
	}
	return time.Now()
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
