package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"gorm.io/gorm"
	dtoreq "motewallet/internal/dto/request"
	dtoresp "motewallet/internal/dto/response"
	bizerrors "motewallet/internal/pkg/errors"
	"motewallet/internal/pkg/kun"
	kundto "motewallet/internal/pkg/kun/dto"
	"motewallet/internal/model"
	"motewallet/internal/repository"
)

type AddressService struct {
	kunClient       kun.KUNClient
	merchantRepo    repository.MerchantRepository
	cryptoAddrRepo  repository.CryptoAddressRepository
	bankAccountRepo repository.BankAccountRepository
}

func NewAddressService(
	kunClient kun.KUNClient,
	merchantRepo repository.MerchantRepository,
	cryptoAddrRepo repository.CryptoAddressRepository,
	bankAccountRepo repository.BankAccountRepository,
) *AddressService {
	return &AddressService{
		kunClient:       kunClient,
		merchantRepo:    merchantRepo,
		cryptoAddrRepo:  cryptoAddrRepo,
		bankAccountRepo: bankAccountRepo,
	}
}

func (s *AddressService) AddCryptoAddress(ctx context.Context, merchantID uint64, req *dtoreq.AddCryptoAddressReq) (*dtoresp.CryptoAddressResp, error) {
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

	var kunResp kundto.CryptoAddressAddResp
	err = s.kunClient.Post(ctx, "/rest/v2.0/customer/crypto/address/add", &kundto.CryptoAddressAddReq{
		SubCustomerNo: *merchant.KunSubCustomerNo,
		Currency:      req.Currency,
		Chain:         req.Chain,
		Address:       req.Address,
		Alias:         req.Alias,
		RequestNo:     kun.GenerateRequestNo(),
	}, &kunResp)
	if err != nil {
		slog.Error("KUN add crypto address failed", slog.Any("error", err))
		return nil, bizerrors.ErrKUNAPIFailedE
	}

	addr := &model.CryptoAddress{
		MerchantID:   merchantID,
		Currency:     req.Currency,
		Chain:        req.Chain,
		Address:      req.Address,
		Alias:        req.Alias,
		KunAccountID: &kunResp.AccountId,
		Status:       "ACTIVE",
	}
	if err := s.cryptoAddrRepo.Create(ctx, addr); err != nil {
		return nil, bizerrors.ErrInternalError
	}

	return &dtoresp.CryptoAddressResp{
		ID:       addr.ID,
		Currency: addr.Currency,
		Chain:    addr.Chain,
		Address:  addr.Address,
		Alias:    addr.Alias,
		Status:   addr.Status,
	}, nil
}

func (s *AddressService) ListCryptoAddresses(ctx context.Context, merchantID uint64) ([]dtoresp.CryptoAddressResp, error) {
	addrs, err := s.cryptoAddrRepo.ListByMerchant(ctx, merchantID)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	var resp []dtoresp.CryptoAddressResp
	for _, a := range addrs {
		resp = append(resp, dtoresp.CryptoAddressResp{
			ID:       a.ID,
			Currency: a.Currency,
			Chain:    a.Chain,
			Address:  a.Address,
			Alias:    a.Alias,
			Status:   a.Status,
		})
	}
	return resp, nil
}

func (s *AddressService) DeleteCryptoAddress(ctx context.Context, merchantID uint64, addressID uint64) error {
	addr, err := s.cryptoAddrRepo.FindByID(ctx, addressID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrNotFoundError
		}
		return bizerrors.ErrInternalError
	}
	if addr.MerchantID != merchantID {
		return bizerrors.ErrForbiddenError
	}
	return s.cryptoAddrRepo.Delete(ctx, addressID)
}

func (s *AddressService) AddBankAccount(ctx context.Context, merchantID uint64, req *dtoreq.AddBankAccountReq) (*dtoresp.BankAccountResp, error) {
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

	if err := validateBankAccountBindReq(req); err != nil {
		return nil, err
	}

	payeeCountryCode := strings.TrimSpace(req.PayeeCountryCode)
	if payeeCountryCode == "" {
		payeeCountryCode = strings.TrimSpace(req.BankCountry)
	}

	kunReq := &kundto.FiatAddressAddReq{
		RequestNo:          kun.GenerateRequestNo(),
		AccountCategory:    "2",
		AccountTypes:       "2",
		CurrencyList:       []string{req.Currency},
		Area:               req.BankCountry,
		TransferType:       req.TransferType,
		AccountNo:          req.AccountNo,
		AccountName:        req.AccountName,
		SwiftCode:          req.SwiftCode,
		PayeeCountryCode:   payeeCountryCode,
		Address:            req.PayeeAddress,
		PayeeAddressSecond: req.PayeeAddressSecond,
		MiddleSwiftCode:    req.MiddleSwiftCode,
		BankName:           req.BankName,
		BankCode:           req.BankCode,
		BankAddress:        req.BankAddress,
		AccountType:        req.AccountType,
	}

	var kunResp kundto.FiatAddressAddResp
	err = s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, kun.FiatAddressAddPath, kunReq, &kunResp)
	if err != nil {
		slog.Error("KUN add bank account failed", slog.Any("error", err))
		return nil, bizerrors.ErrKUNAPIFailedE
	}

	account := &model.BankAccount{
		MerchantID:       merchantID,
		KunAccountID:     &kunResp.AccountId,
		CurrencyList:     req.Currency,
		TransferType:     req.TransferType,
		AccountNo:        req.AccountNo,
		AccountName:      req.AccountName,
		BankName:         &req.BankName,
		SwiftCode:        stringPtr(strings.TrimSpace(req.SwiftCode)),
		BankCode:         stringPtr(strings.TrimSpace(req.BankCode)),
		PayeeCountryCode: &payeeCountryCode,
		PayeeAddress:     stringPtr(strings.TrimSpace(req.PayeeAddress)),
		MiddleSwiftCode:  stringPtr(strings.TrimSpace(req.MiddleSwiftCode)),
		Area:             req.BankCountry,
		Status:           "ACTIVE",
	}
	if err := s.bankAccountRepo.Create(ctx, account); err != nil {
		return nil, bizerrors.ErrInternalError
	}

	return &dtoresp.BankAccountResp{
		ID:              account.ID,
		Currency:        account.CurrencyList,
		BankName:        req.BankName,
		BankCountry:     req.BankCountry,
		SwiftCode:       req.SwiftCode,
		AccountName:     account.AccountName,
		AccountNoMasked: maskAccountNo(account.AccountNo),
		TransferType:    account.TransferType,
		Status:          account.Status,
	}, nil
}

func (s *AddressService) ListBankAccounts(ctx context.Context, merchantID uint64) ([]dtoresp.BankAccountResp, error) {
	accounts, err := s.bankAccountRepo.ListByMerchant(ctx, merchantID)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	var resp = make([]dtoresp.BankAccountResp, 0)
	for _, a := range accounts {
		bankName := ""
		if a.BankName != nil {
			bankName = *a.BankName
		}
		swiftCode := ""
		if a.SwiftCode != nil {
			swiftCode = *a.SwiftCode
		}
		bankCountry := ""
		if a.PayeeCountryCode != nil {
			bankCountry = *a.PayeeCountryCode
		}
		resp = append(resp, dtoresp.BankAccountResp{
			ID:              a.ID,
			Currency:        a.CurrencyList,
			BankName:        bankName,
			BankCountry:     bankCountry,
			SwiftCode:       swiftCode,
			AccountName:     a.AccountName,
			AccountNoMasked: maskAccountNo(a.AccountNo),
			TransferType:    a.TransferType,
			Status:          a.Status,
		})
	}
	return resp, nil
}

func (s *AddressService) DeleteBankAccount(ctx context.Context, merchantID uint64, accountID uint64) error {
	account, err := s.bankAccountRepo.FindByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrNotFoundError
		}
		return bizerrors.ErrInternalError
	}
	if account.MerchantID != merchantID {
		return bizerrors.ErrForbiddenError
	}
	if account.KunAccountID == nil || strings.TrimSpace(*account.KunAccountID) == "" {
		return bizerrors.ErrKUNAPIFailedE
	}

	merchant, err := s.merchantRepo.FindByID(ctx, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.ErrNotFoundError
		}
		return bizerrors.ErrInternalError
	}
	if merchant.KunSubCustomerNo == nil {
		return bizerrors.ErrMerchantNotRegisteredE
	}

	err = s.kunClient.PostAsCustomer(ctx, *merchant.KunSubCustomerNo, kun.FiatAddressDelPath, &kundto.FiatAddressDelReq{
		RequestNo: kun.GenerateRequestNo(),
		AccountId: *account.KunAccountID,
		Currency:  account.CurrencyList,
	}, &struct{}{})
	if err != nil {
		slog.Error("KUN unbind bank account failed", slog.Any("error", err))
		return bizerrors.ErrKUNAPIFailedE
	}

	return s.bankAccountRepo.Delete(ctx, accountID)
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func maskAccountNo(accountNo string) string {
	accountNo = strings.TrimSpace(accountNo)
	if accountNo == "" {
		return ""
	}
	if len(accountNo) <= 4 {
		return strings.Repeat("*", len(accountNo))
	}
	return strings.Repeat("*", len(accountNo)-4) + accountNo[len(accountNo)-4:]
}
