package verification

import "errors"

var (
	ErrInvalidCode      = errors.New("invalid verification code")
	ErrExpiredCode      = errors.New("verification code expired")
	ErrSendTooFrequent  = errors.New("verification code sent too frequently")
)
