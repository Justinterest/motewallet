package service

import (
	dtoreq "motewallet/internal/dto/request"
	bizerrors "motewallet/internal/pkg/errors"
	"strings"
)

var allowedFiatWithdrawalPurposes = map[string]bool{
	"OTHER": true,
	"TRAD":  true,
	"INVS":  true,
	"GDDS":  true,
}

func validateBankAccountBindReq(req *dtoreq.AddBankAccountReq) error {
	transferType := strings.ToUpper(strings.TrimSpace(req.TransferType))
	req.TransferType = transferType

	if transferType != "TT" {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "only TT transfer type is supported")
	}

	if strings.TrimSpace(req.SwiftCode) == "" {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "swift_code is required for TT transfer type")
	}
	payeeCountry := strings.TrimSpace(req.PayeeCountryCode)
	if payeeCountry == "" {
		payeeCountry = strings.TrimSpace(req.BankCountry)
	}
	if payeeCountry == "" {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "payee_country_code is required for TT transfer type")
	}
	if strings.TrimSpace(req.PayeeAddress) == "" {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "payee_address is required for TT transfer type")
	}
	if strings.TrimSpace(req.PayeeAddressSecond) == "" {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "payee_address_second is required for TT transfer type")
	}

	req.AccountType = "ENTERPRISE"

	return nil
}

func validateFiatWithdrawalPurpose(purpose string) error {
	purpose = strings.ToUpper(strings.TrimSpace(purpose))
	if !allowedFiatWithdrawalPurposes[purpose] {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "invalid purpose; allowed: OTHER, TRAD, INVS, GDDS")
	}
	return nil
}
