export interface Agreement {
  id: string;
  title: string;
  content: string;
  required: boolean;
}

export interface AgreementListResponse {
  agreements: Agreement[];
  signed: boolean;
}

export interface KycStatusResponse {
  status: string;        // merchant status: PENDING_AGREEMENT, PENDING_KYC, PROCESSING, ACTIVE, FROZEN
  kyc_status: string;    // NONE, PENDING, PROCESSING, AUTH_SUC, AUTH_FAIL
  kyc_fail_reason?: string;
  kyc_submitted_at?: string;
  kyc_completed_at?: string;
  agreement_signed_at?: string;
}

export interface SubmitKycRequest {
  company_name: string;
  country: string;
  registration_number: string;
  contact_name: string;
  contact_phone: string;
}
