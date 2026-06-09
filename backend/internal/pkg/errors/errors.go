package errors

import "net/http"

// Business error codes
const (
	ErrInvalidCredentials          = 10001
	ErrEmailAlreadyExists          = 10002
	ErrInvalidVerificationCode     = 10008
	ErrVerificationCodeExpired     = 10009
	ErrVerificationSendTooFrequent = 10015
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
	Data       interface{}
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

func NewKycValidationError(errors []string) *BusinessError {
	return &BusinessError{
		HTTPStatus: http.StatusBadRequest,
		Code:       ErrValidation,
		Message:    "实名认证信息校验失败，请检查以下字段",
		Data: map[string]interface{}{
			"errors": errors,
		},
	}
}

// Pre-defined business errors
var (
	ErrInvalidCredentialsError      = NewBusinessError(http.StatusUnauthorized, ErrInvalidCredentials, "invalid email or password")
	ErrEmailAlreadyExistsError      = NewBusinessError(http.StatusConflict, ErrEmailAlreadyExists, "email already exists")
	ErrInvalidVerificationCodeE     = NewBusinessError(http.StatusBadRequest, ErrInvalidVerificationCode, "invalid verification code")
	ErrVerificationCodeExpiredE     = NewBusinessError(http.StatusBadRequest, ErrVerificationCodeExpired, "verification code expired")
	ErrVerificationSendTooFrequentE = NewBusinessError(http.StatusTooManyRequests, ErrVerificationSendTooFrequent, "verification code sent too frequently")
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
	ErrMerchantNotRegisteredE       = NewBusinessError(http.StatusBadRequest, ErrMerchantNotRegistered, "请先完成企业认证后再使用该功能")
	ErrKycAlreadyProcessingE        = NewBusinessError(http.StatusBadRequest, ErrKycAlreadyProcessing, "企业认证审核中，请耐心等待")
	ErrInsufficientBalanceE         = NewBusinessError(http.StatusBadRequest, ErrInsufficientBalance, "余额不足")
	ErrWithdrawalNotFoundE          = NewBusinessError(http.StatusNotFound, ErrWithdrawalNotFound, "提现记录不存在")
	ErrWithdrawalNotPendingE        = NewBusinessError(http.StatusBadRequest, ErrWithdrawalNotPending, "该提现申请当前无法审核")
	ErrStorageNotConfiguredE        = NewBusinessError(http.StatusServiceUnavailable, ErrStorageNotConfigured, "文件服务暂不可用，请稍后重试")
	ErrKUNAPIFailedE                = NewBusinessError(http.StatusBadGateway, ErrKUNAPIFailed, "服务暂时不可用，请稍后重试")
	ErrKUNSignatureFailedE          = NewBusinessError(http.StatusUnauthorized, ErrKUNSignatureFailed, "通知验签失败")
	ErrKUNTimestampExpiredE         = NewBusinessError(http.StatusUnauthorized, ErrKUNTimestampExpired, "通知已过期")
	ErrInternalError                = NewBusinessError(http.StatusInternalServerError, ErrInternal, "internal server error")
)
