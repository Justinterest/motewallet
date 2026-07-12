package service

import (
	"context"
	"errors"
	"fmt"

	dtoresp "motewallet/internal/dto/response"
	"motewallet/internal/model"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/repository"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var (
	accountTypes = []string{"FUNDING", "TRADING"}
	currencies   = []string{"USDT", "USDC", "BTC", "USD", "HKD", "EUR"}
)

// WalletChangeRef links a wallet mutation to a business transaction_record.
// TransactionRecordID must be set before mutating the wallet so ledger rows are auditable.
type WalletChangeRef struct {
	TransactionRecordID uint64
	BizType             string // DEPOSIT / WITHDRAWAL / EXCHANGE / TRANSFER
	Remark              string
}

type WalletService struct {
	merchantWalletRepo repository.MerchantWalletRepository
	walletLedgerRepo   repository.WalletLedgerRepository
	merchantRepo       repository.MerchantRepository
	currencyConfigSvc  *CurrencyConfigService
}

func NewWalletService(
	merchantWalletRepo repository.MerchantWalletRepository,
	walletLedgerRepo repository.WalletLedgerRepository,
	merchantRepo repository.MerchantRepository,
	currencyConfigSvc *CurrencyConfigService,
) *WalletService {
	return &WalletService{
		merchantWalletRepo: merchantWalletRepo,
		walletLedgerRepo:   walletLedgerRepo,
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

func (s *WalletService) FreezeBalance(ctx context.Context, tx *gorm.DB, merchantID uint64, accountType, currency string, amount decimal.Decimal, ref WalletChangeRef) error {
	if err := validatePositiveAmount(amount); err != nil {
		return err
	}
	wallet, err := s.merchantWalletRepo.FindByMerchantAccountCurrencyWithDB(ctx, tx, merchantID, accountType, currency)
	if err != nil {
		return fmt.Errorf("find wallet: %w", err)
	}
	available := wallet.Balance.Sub(wallet.FrozenBalance)
	if available.LessThan(amount) {
		return bizerrors.ErrInsufficientBalanceE
	}

	balanceBefore := wallet.Balance
	frozenBefore := wallet.FrozenBalance
	wallet.FrozenBalance = wallet.FrozenBalance.Add(amount)

	if err := s.merchantWalletRepo.UpdateBalanceWithVersion(ctx, tx, wallet); err != nil {
		return err
	}
	return s.appendLedger(ctx, tx, wallet, model.WalletLedgerFreeze, amount, balanceBefore, wallet.Balance, frozenBefore, wallet.FrozenBalance, ref)
}

func (s *WalletService) UnfreezeBalance(ctx context.Context, tx *gorm.DB, merchantID uint64, accountType, currency string, amount decimal.Decimal, ref WalletChangeRef) error {
	if err := validatePositiveAmount(amount); err != nil {
		return err
	}
	wallet, err := s.merchantWalletRepo.FindByMerchantAccountCurrencyWithDB(ctx, tx, merchantID, accountType, currency)
	if err != nil {
		return fmt.Errorf("find wallet: %w", err)
	}
	if wallet.FrozenBalance.LessThan(amount) {
		return bizerrors.ErrInsufficientBalanceE
	}

	balanceBefore := wallet.Balance
	frozenBefore := wallet.FrozenBalance
	wallet.FrozenBalance = wallet.FrozenBalance.Sub(amount)

	if err := s.merchantWalletRepo.UpdateBalanceWithVersion(ctx, tx, wallet); err != nil {
		return err
	}
	return s.appendLedger(ctx, tx, wallet, model.WalletLedgerUnfreeze, amount, balanceBefore, wallet.Balance, frozenBefore, wallet.FrozenBalance, ref)
}

func (s *WalletService) DeductFrozen(ctx context.Context, tx *gorm.DB, merchantID uint64, accountType, currency string, amount decimal.Decimal, ref WalletChangeRef) error {
	if err := validatePositiveAmount(amount); err != nil {
		return err
	}
	wallet, err := s.merchantWalletRepo.FindByMerchantAccountCurrencyWithDB(ctx, tx, merchantID, accountType, currency)
	if err != nil {
		return fmt.Errorf("find wallet: %w", err)
	}
	if wallet.FrozenBalance.LessThan(amount) {
		return bizerrors.ErrInsufficientBalanceE
	}

	balanceBefore := wallet.Balance
	frozenBefore := wallet.FrozenBalance
	wallet.Balance = wallet.Balance.Sub(amount)
	wallet.FrozenBalance = wallet.FrozenBalance.Sub(amount)

	if err := s.merchantWalletRepo.UpdateBalanceWithVersion(ctx, tx, wallet); err != nil {
		return err
	}
	return s.appendLedger(ctx, tx, wallet, model.WalletLedgerDeductFrozen, amount, balanceBefore, wallet.Balance, frozenBefore, wallet.FrozenBalance, ref)
}

func (s *WalletService) CreditBalance(ctx context.Context, tx *gorm.DB, merchantID uint64, accountType, currency string, amount decimal.Decimal, ref WalletChangeRef) error {
	if err := validatePositiveAmount(amount); err != nil {
		return err
	}
	wallet, err := s.merchantWalletRepo.FindByMerchantAccountCurrencyWithDB(ctx, tx, merchantID, accountType, currency)
	if err != nil {
		return fmt.Errorf("find wallet: %w", err)
	}

	balanceBefore := wallet.Balance
	frozenBefore := wallet.FrozenBalance
	wallet.Balance = wallet.Balance.Add(amount)

	if err := s.merchantWalletRepo.UpdateBalanceWithVersion(ctx, tx, wallet); err != nil {
		return err
	}
	return s.appendLedger(ctx, tx, wallet, model.WalletLedgerCredit, amount, balanceBefore, wallet.Balance, frozenBefore, wallet.FrozenBalance, ref)
}

func validatePositiveAmount(amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "amount must be positive")
	}
	return nil
}

func (s *WalletService) appendLedger(
	ctx context.Context,
	tx *gorm.DB,
	wallet *model.MerchantWallet,
	entryType string,
	amount decimal.Decimal,
	balanceBefore, balanceAfter, frozenBefore, frozenAfter decimal.Decimal,
	ref WalletChangeRef,
) error {
	entry := &model.WalletLedger{
		MerchantID:    wallet.MerchantID,
		WalletID:      wallet.ID,
		AccountType:   wallet.AccountType,
		Currency:      wallet.Currency,
		EntryType:     entryType,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		FrozenBefore:  frozenBefore,
		FrozenAfter:   frozenAfter,
	}
	if ref.TransactionRecordID > 0 {
		id := ref.TransactionRecordID
		entry.TransactionRecordID = &id
	}
	if ref.BizType != "" {
		bizType := ref.BizType
		entry.BizType = &bizType
	}
	if ref.Remark != "" {
		remark := ref.Remark
		entry.Remark = &remark
	}
	return s.walletLedgerRepo.Create(ctx, tx, entry)
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

func (s *WalletService) ListLedger(
	ctx context.Context,
	merchantID uint64,
	accountType, currency, bizType, entryType string,
	page, pageSize int,
) (*dtoresp.WalletLedgerListResp, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	items, total, err := s.walletLedgerRepo.ListByMerchant(ctx, merchantID, repository.WalletLedgerListFilter{
		AccountType: accountType,
		Currency:    currency,
		BizType:     bizType,
		EntryType:   entryType,
		Page:        page,
		PageSize:    pageSize,
	})
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	entries := make([]dtoresp.WalletLedgerEntryResp, 0, len(items))
	for _, item := range items {
		entry := dtoresp.WalletLedgerEntryResp{
			ID:                  item.ID,
			AccountType:         item.AccountType,
			Currency:            item.Currency,
			EntryType:           item.EntryType,
			Amount:              item.Amount.String(),
			BalanceBefore:       item.BalanceBefore.String(),
			BalanceAfter:        item.BalanceAfter.String(),
			FrozenBefore:        item.FrozenBefore.String(),
			FrozenAfter:         item.FrozenAfter.String(),
			TransactionRecordID: item.TransactionRecordID,
			PlatformOrderID:     item.PlatformOrderID,
			BizType:             item.BizType,
			Remark:              item.Remark,
			CreatedAt:           item.CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"),
		}
		entries = append(entries, entry)
	}

	return &dtoresp.WalletLedgerListResp{
		Entries:  entries,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
