package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	dtoresp "motewallet/internal/dto/response"
	"motewallet/internal/model"
	bizerrors 	"motewallet/internal/pkg/errors"
	"motewallet/internal/repository"
)

var (
	accountTypes = []string{"FUNDING", "TRADING"}
	currencies   = []string{"USDT", "USDC", "BTC", "USD", "HKD", "EUR"}
)

type WalletService struct {
	merchantWalletRepo repository.MerchantWalletRepository
	merchantRepo       repository.MerchantRepository
	currencyConfigSvc  *CurrencyConfigService
}

func NewWalletService(
	merchantWalletRepo repository.MerchantWalletRepository,
	merchantRepo repository.MerchantRepository,
	currencyConfigSvc *CurrencyConfigService,
) *WalletService {
	return &WalletService{
		merchantWalletRepo: merchantWalletRepo,
		merchantRepo:       merchantRepo,
		currencyConfigSvc:  currencyConfigSvc,
	}
}

func (s *WalletService) InitializeWallets(ctx context.Context, merchantID uint64) error {
	var wallets []*model.MerchantWallet
	for _, accountType := range accountTypes {
		for _, currency := range currencies {
			wallets = append(wallets, &model.MerchantWallet{
				MerchantID:    merchantID,
				AccountType:   accountType,
				Currency:      currency,
				Balance:       decimal.Zero,
				FrozenBalance: decimal.Zero,
				Version:       0,
			})
		}
	}
	return s.merchantWalletRepo.CreateBatch(ctx, wallets)
}

func (s *WalletService) FreezeBalance(ctx context.Context, tx *gorm.DB, merchantID uint64, accountType, currency string, amount decimal.Decimal) error {
	wallet, err := s.merchantWalletRepo.FindByMerchantAccountCurrencyWithDB(ctx, tx, merchantID, accountType, currency)
	if err != nil {
		return fmt.Errorf("find wallet: %w", err)
	}
	available := wallet.Balance.Sub(wallet.FrozenBalance)
	if available.LessThan(amount) {
		return bizerrors.ErrInsufficientBalanceE
	}
	wallet.FrozenBalance = wallet.FrozenBalance.Add(amount)
	return s.merchantWalletRepo.UpdateBalanceWithVersion(ctx, tx, wallet)
}

func (s *WalletService) UnfreezeBalance(ctx context.Context, tx *gorm.DB, merchantID uint64, accountType, currency string, amount decimal.Decimal) error {
	wallet, err := s.merchantWalletRepo.FindByMerchantAccountCurrencyWithDB(ctx, tx, merchantID, accountType, currency)
	if err != nil {
		return fmt.Errorf("find wallet: %w", err)
	}
	if wallet.FrozenBalance.LessThan(amount) {
		return bizerrors.ErrInsufficientBalanceE
	}
	wallet.FrozenBalance = wallet.FrozenBalance.Sub(amount)
	return s.merchantWalletRepo.UpdateBalanceWithVersion(ctx, tx, wallet)
}

func (s *WalletService) DeductFrozen(ctx context.Context, tx *gorm.DB, merchantID uint64, accountType, currency string, amount decimal.Decimal) error {
	wallet, err := s.merchantWalletRepo.FindByMerchantAccountCurrencyWithDB(ctx, tx, merchantID, accountType, currency)
	if err != nil {
		return fmt.Errorf("find wallet: %w", err)
	}
	if wallet.FrozenBalance.LessThan(amount) {
		return bizerrors.ErrInsufficientBalanceE
	}
	wallet.Balance = wallet.Balance.Sub(amount)
	wallet.FrozenBalance = wallet.FrozenBalance.Sub(amount)
	return s.merchantWalletRepo.UpdateBalanceWithVersion(ctx, tx, wallet)
}

func (s *WalletService) CreditBalance(ctx context.Context, tx *gorm.DB, merchantID uint64, accountType, currency string, amount decimal.Decimal) error {
	wallet, err := s.merchantWalletRepo.FindByMerchantAccountCurrencyWithDB(ctx, tx, merchantID, accountType, currency)
	if err != nil {
		return fmt.Errorf("find wallet: %w", err)
	}
	wallet.Balance = wallet.Balance.Add(amount)
	return s.merchantWalletRepo.UpdateBalanceWithVersion(ctx, tx, wallet)
}

func (s *WalletService) GetBalances(ctx context.Context, merchantID uint64) (*dtoresp.WalletBalancesResp, error) {
	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerrors.ErrNotFoundError
		}
		return nil, bizerrors.ErrInternalError
	}

	supportedCrypto, err := s.currencyConfigSvc.GetSupportedCrypto(ctx, merchant)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}
	supportedFiat, err := s.currencyConfigSvc.GetSupportedFiat(ctx, merchant)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}
	supported := make(map[string]struct{}, len(supportedCrypto)+len(supportedFiat))
	for _, code := range append(supportedCrypto, supportedFiat...) {
		supported[code] = struct{}{}
	}

	wallets, err := s.merchantWalletRepo.FindByMerchantID(ctx, merchantID)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	walletByKey := make(map[string]dtoresp.WalletBalanceResp, len(wallets))
	for _, w := range wallets {
		if _, ok := supported[w.Currency]; !ok {
			continue
		}
		available := w.Balance.Sub(w.FrozenBalance)
		key := w.AccountType + ":" + w.Currency
		walletByKey[key] = dtoresp.WalletBalanceResp{
			AccountType:      w.AccountType,
			Currency:         w.Currency,
			Balance:          w.Balance.String(),
			FrozenBalance:    w.FrozenBalance.String(),
			AvailableBalance: available.String(),
		}
	}

	orderedCurrencies := append(supportedFiat, supportedCrypto...)
	var resp []dtoresp.WalletBalanceResp
	for _, accountType := range accountTypes {
		for _, code := range orderedCurrencies {
			key := accountType + ":" + code
			if item, ok := walletByKey[key]; ok {
				resp = append(resp, item)
			}
		}
	}

	return &dtoresp.WalletBalancesResp{Wallets: resp}, nil
}
