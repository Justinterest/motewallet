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

	switch transferType {
	case "TT", "CHATS":
		if strings.TrimSpace(req.SwiftCode) == "" {
			return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "swift_code is required for this transfer type")
		}
		payeeCountry := strings.TrimSpace(req.PayeeCountryCode)
		if payeeCountry == "" {
			payeeCountry = strings.TrimSpace(req.BankCountry)
		}
		if payeeCountry == "" {
			return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "payee_country_code is required for this transfer type")
		}
		if strings.TrimSpace(req.PayeeAddress) == "" {
			return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "payee_address is required for this transfer type")
		}
		if strings.TrimSpace(req.PayeeAddressSecond) == "" {
			return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "payee_address_second is required for this transfer type")
		}
	case "LOCAL":
	default:
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "invalid transfer_type")
	}

	if transferType == "CHATS" && strings.TrimSpace(req.BankCode) == "" {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "bank_code is required for CHATS transfer type")
	}

	if transferType == "TT" && strings.TrimSpace(req.AccountType) == "" {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "account_type is required for TT transfer type")
	}

	return nil
}

func validateFiatWithdrawalPurpose(purpose string) error {
	purpose = strings.ToUpper(strings.TrimSpace(purpose))
	if !allowedFiatWithdrawalPurposes[purpose] {
		return bizerrors.NewBusinessError(400, bizerrors.ErrValidation, "invalid purpose; allowed: OTHER, TRAD, INVS, GDDS")
	}
	return nil
}
