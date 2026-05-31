package errors

import "net/http"

// Business error codes
const (
	ErrInvalidCredentials          = 10001
	ErrEmailAlreadyExists          = 10002
	ErrInvalidToken                = 10003
	ErrUnauthorized                = 10004
	ErrForbidden                   = 10005
	ErrValidation                  = 10006
	ErrNotFound                    = 10007
	ErrMerchantNotPendingAgreement = 10010
	ErrMerchantNotPendingKyc       = 10011
	ErrMerchantFrozen              = 10012
	ErrFeeTemplateInUse            = 10013
	ErrInvalidStatusTransition     = 10014
	ErrWebhookDuplicate            = 10020
	ErrMerchantNotRegistered       = 10024
	ErrKycAlreadyProcessing        = 10025
	ErrInsufficientBalance         = 10030
	ErrWithdrawalNotFound          = 10031
	ErrWithdrawalNotPending        = 10032
	ErrInternal                    = 50000
	ErrStorageNotConfigured        = 10040
	ErrKUNAPIFailed                = 50001
	ErrKUNSignatureFailed          = 50002
	ErrKUNTimestampExpired         = 50003
)

// BusinessError represents a business logic error with HTTP status and error code.
type BusinessError struct {
	HTTPStatus int
	Code       int
	Message    string
}

func (e *BusinessError) Error() string {
	return e.Message
}

func NewBusinessError(httpStatus, code int, message string) *BusinessError {
	return &BusinessError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
	}
}

// Pre-defined business errors
var (
	ErrInvalidCredentialsError      = NewBusinessError(http.StatusUnauthorized, ErrInvalidCredentials, "invalid email or password")
	ErrEmailAlreadyExistsError      = NewBusinessError(http.StatusConflict, ErrEmailAlreadyExists, "email already exists")
	ErrInvalidTokenError            = NewBusinessError(http.StatusUnauthorized, ErrInvalidToken, "invalid or expired token")
	ErrUnauthorizedError            = NewBusinessError(http.StatusUnauthorized, ErrUnauthorized, "unauthorized")
	ErrForbiddenError               = NewBusinessError(http.StatusForbidden, ErrForbidden, "forbidden")
	ErrValidationError              = NewBusinessError(http.StatusBadRequest, ErrValidation, "validation failed")
	ErrNotFoundError                = NewBusinessError(http.StatusNotFound, ErrNotFound, "resource not found")
	ErrMerchantNotPendingAgreementE = NewBusinessError(http.StatusBadRequest, ErrMerchantNotPendingAgreement, "merchant is not in PENDING_AGREEMENT status")
	ErrMerchantNotPendingKycE       = NewBusinessError(http.StatusBadRequest, ErrMerchantNotPendingKyc, "merchant is not in PENDING_KYC status")
	ErrMerchantFrozenE              = NewBusinessError(http.StatusForbidden, ErrMerchantFrozen, "merchant account is frozen")
	ErrFeeTemplateInUseE            = NewBusinessError(http.StatusConflict, ErrFeeTemplateInUse, "fee template is referenced by merchants and cannot be deleted")
	ErrInvalidStatusTransitionE     = NewBusinessError(http.StatusBadRequest, ErrInvalidStatusTransition, "invalid status transition")
	ErrWebhookDuplicateE            = NewBusinessError(http.StatusOK, ErrWebhookDuplicate, "webhook event already processed")
	ErrMerchantNotRegisteredE       = NewBusinessError(http.StatusBadRequest, ErrMerchantNotRegistered, "merchant not registered with KUN")
	ErrKycAlreadyProcessingE        = NewBusinessError(http.StatusBadRequest, ErrKycAlreadyProcessing, "KYC verification is already in progress")
	ErrInsufficientBalanceE         = NewBusinessError(http.StatusBadRequest, ErrInsufficientBalance, "insufficient balance")
	ErrWithdrawalNotFoundE          = NewBusinessError(http.StatusNotFound, ErrWithdrawalNotFound, "withdrawal order not found")
	ErrWithdrawalNotPendingE        = NewBusinessError(http.StatusBadRequest, ErrWithdrawalNotPending, "withdrawal order is not in pending review status")
	ErrStorageNotConfiguredE        = NewBusinessError(http.StatusServiceUnavailable, ErrStorageNotConfigured, "file storage is not configured")
	ErrKUNAPIFailedE                = NewBusinessError(http.StatusBadGateway, ErrKUNAPIFailed, "KUN API call failed")
	ErrKUNSignatureFailedE          = NewBusinessError(http.StatusUnauthorized, ErrKUNSignatureFailed, "webhook signature verification failed")
	ErrKUNTimestampExpiredE         = NewBusinessError(http.StatusUnauthorized, ErrKUNTimestampExpired, "webhook timestamp expired")
	ErrInternalError                = NewBusinessError(http.StatusInternalServerError, ErrInternal, "internal server error")
)
