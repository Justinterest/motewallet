package service

import (
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"
	dtoreq "motewallet-withdrawal/backend/internal/dto/request"
	dtoresp "motewallet-withdrawal/backend/internal/dto/response"
	bizerrors "motewallet-withdrawal/backend/internal/pkg/errors"
	"motewallet-withdrawal/backend/internal/pkg/kun"
	kundto "motewallet-withdrawal/backend/internal/pkg/kun/dto"
	"motewallet-withdrawal/backend/internal/model"
	"motewallet-withdrawal/backend/internal/repository"
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

	var kunResp kundto.FiatAddressAddResp
	err = s.kunClient.Post(ctx, "/rest/v2.0/customer/fiat/address/add", &kundto.FiatAddressAddReq{
		SubCustomerNo: *merchant.KunSubCustomerNo,
		Currency:      req.Currency,
		BankName:      req.BankName,
		BankCountry:   req.BankCountry,
		SwiftCode:     req.SwiftCode,
		AccountName:   req.AccountName,
		AccountNo:     req.AccountNo,
		TransferType:  req.TransferType,
		RequestNo:     kun.GenerateRequestNo(),
	}, &kunResp)
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
		SwiftCode:        &req.SwiftCode,
		PayeeCountryCode: &req.BankCountry,
		Area:             req.BankCountry,
		Status:           "ACTIVE",
	}
	if err := s.bankAccountRepo.Create(ctx, account); err != nil {
		return nil, bizerrors.ErrInternalError
	}

	return &dtoresp.BankAccountResp{
		ID:           account.ID,
		Currency:     account.CurrencyList,
		BankName:     req.BankName,
		BankCountry:  req.BankCountry,
		SwiftCode:    req.SwiftCode,
		AccountName:  account.AccountName,
		TransferType: account.TransferType,
		Status:       account.Status,
	}, nil
}

func (s *AddressService) ListBankAccounts(ctx context.Context, merchantID uint64) ([]dtoresp.BankAccountResp, error) {
	accounts, err := s.bankAccountRepo.ListByMerchant(ctx, merchantID)
	if err != nil {
		return nil, bizerrors.ErrInternalError
	}

	var resp []dtoresp.BankAccountResp
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
			ID:           a.ID,
			Currency:     a.CurrencyList,
			BankName:     bankName,
			BankCountry:  bankCountry,
			SwiftCode:    swiftCode,
			AccountName:  a.AccountName,
			TransferType: a.TransferType,
			Status:       a.Status,
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
	return s.bankAccountRepo.Delete(ctx, accountID)
}
