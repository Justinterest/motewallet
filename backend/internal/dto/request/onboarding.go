package request

import kundto "motewallet/internal/pkg/kun/dto"

// SubmitKycReq matches KUN sub-merchant onboarding authentication body.
// See: https://opendocs.kun.global/docs/api/sub-merchant-onboarding-authentication
type SubmitKycReq = kundto.SubMerchantRegisterReq
