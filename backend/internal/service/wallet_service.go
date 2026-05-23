package service

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	dtoresp "motewallet-withdrawal/backend/internal/dto/response"
	"motewallet-withdrawal/backend/internal/model"
	bizerrors "motewallet-withdrawal/backend/internal/pkg/errors"
	"motewallet-withdrawal/backend/internal/repository"
)

var (
	accountTypes = []string{"FUNDING", "TRADING"}
	currencies   = []string{"USDT", "USDC", "BTC", "USD", "HKD", "EUR"}
)

type WalletService struct {
	merchantWalletRepo repository.MerchantWalletRepository
}

func NewWalletService(merchantWalletRepo repository.MerchantWalletRepository) *WalletService {
	return &WalletService{
		merchantWalletRepo: merchantWalletRepo,
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
	wallets, err := s.merchantWalletRepo.FindByMerchantID(ctx, merchantID)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	var resp []dtoresp.WalletBalanceResp
	for _, w := range wallets {
		available := w.Balance.Sub(w.FrozenBalance)
		resp = append(resp, dtoresp.WalletBalanceResp{
			AccountType:      w.AccountType,
			Currency:         w.Currency,
			Balance:          w.Balance.String(),
			FrozenBalance:    w.FrozenBalance.String(),
			AvailableBalance: available.String(),
		})
	}

	return &dtoresp.WalletBalancesResp{Wallets: resp}, nil
}
