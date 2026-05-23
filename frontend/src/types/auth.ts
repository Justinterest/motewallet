export interface User {
  id: number;
  email: string;
  status: string;
  kyc_status: string;
  created_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
}
