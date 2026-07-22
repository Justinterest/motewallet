export interface AdminUser {
  id: number;
  username: string;
  email: string;
  role: string;
}

export type AuthChallengeStatus =
  | "SUCCESS"
  | "REQUIRES_2FA"
  | "REQUIRES_2FA_SETUP"
  | "REQUIRES_PASSWORD_CHANGE";

export interface AdminAuthChallenge {
  status: AuthChallengeStatus;
  temp_token?: string;
  totp_secret?: string;
  totp_uri?: string;
  admin?: AdminUser;
}

export interface AdminLoginRequest {
  username: string;
  password: string;
}

export interface TotpVerifyRequest {
  temp_token: string;
  code: string;
}

export interface AdminChangePasswordRequest {
  temp_token: string;
  new_password: string;
}

export interface AdminEmployee {
  id: number;
  username: string;
  email: string;
  role: string;
  status: string;
  totp_enabled: boolean;
  last_login_at?: string | null;
  created_at: string;
}

export interface CreateAdminEmployeeRequest {
  username: string;
  email: string;
  role: string;
  password?: string;
}

export interface CreateAdminEmployeeResponse {
  user: AdminEmployee;
  initial_password?: string;
}

export interface ResetAdminPasswordResponse {
  new_password: string;
}
