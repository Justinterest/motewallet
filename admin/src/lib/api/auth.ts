import apiClient from "./client";
import type {
  AdminAuthChallenge,
  AdminChangePasswordRequest,
  AdminLoginRequest,
  AdminUser,
  TotpVerifyRequest,
} from "@/types/auth";

export const adminAuthApi = {
  login: (data: AdminLoginRequest) =>
    apiClient.post<unknown, AdminAuthChallenge>("/api/v1/admin/auth/login", data),

  verify2FA: (data: TotpVerifyRequest) =>
    apiClient.post<unknown, AdminAuthChallenge>(
      "/api/v1/admin/auth/2fa/verify",
      data
    ),

  confirm2FASetup: (data: TotpVerifyRequest) =>
    apiClient.post<unknown, AdminAuthChallenge>(
      "/api/v1/admin/auth/2fa/setup/confirm",
      data
    ),

  changePassword: (data: AdminChangePasswordRequest) =>
    apiClient.post<unknown, AdminAuthChallenge>(
      "/api/v1/admin/auth/change-password",
      data
    ),

  logout: () => apiClient.post("/api/v1/admin/auth/logout"),

  getMe: () => apiClient.get<unknown, AdminUser>("/api/v1/admin/auth/me"),
};
