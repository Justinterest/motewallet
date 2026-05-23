export interface AdminUser {
  id: number;
  username: string;
  email: string;
  role: string;
}

export interface AdminLoginRequest {
  username: string;
  password: string;
}
