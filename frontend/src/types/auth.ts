export interface User {
  id: number;
  email: string;
  status: string;
  kyc_status: string;
  totp_enabled: boolean;
  created_at: string;
}

export type AuthChallengeStatus =
  | "SUCCESS"
  | "REQUIRES_2FA"
  | "REQUIRES_2FA_SETUP";

export interface AuthChallenge {
  status: AuthChallengeStatus;
  temp_token?: string;
  totp_secret?: string;
  totp_uri?: string;
  merchant?: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  verification_code: string;
}

export interface SendVerificationCodeRequest {
  email: string;
}

export interface TotpVerifyRequest {
  temp_token: string;
  code: string;
}

export interface TotpStatus {
  enabled: boolean;
}

export interface TotpSetup {
  totp_secret: string;
  totp_uri: string;
}
