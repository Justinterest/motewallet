import apiClient from "./client";
import type {
  AuthChallenge,
  LoginRequest,
  RegisterRequest,
  SendVerificationCodeRequest,
  TotpSetup,
  TotpStatus,
  TotpVerifyRequest,
  User,
} from "@/types/auth";

export const authApi = {
  login: (data: LoginRequest) =>
    apiClient.post<unknown, AuthChallenge>("/api/v1/auth/login", data),

  sendVerificationCode: (data: SendVerificationCodeRequest) =>
    apiClient.post<unknown, null>("/api/v1/auth/send-verification-code", data),

  register: (data: RegisterRequest) =>
    apiClient.post<unknown, AuthChallenge>("/api/v1/auth/register", data),

  verify2FA: (data: TotpVerifyRequest) =>
    apiClient.post<unknown, AuthChallenge>("/api/v1/auth/2fa/verify", data),

  confirm2FASetup: (data: TotpVerifyRequest) =>
    apiClient.post<unknown, AuthChallenge>("/api/v1/auth/2fa/setup/confirm", data),

  getTotpStatus: () =>
    apiClient.get<unknown, TotpStatus>("/api/v1/auth/2fa/status"),

  prepareTotpRebind: (currentCode: string) =>
    apiClient.post<unknown, TotpSetup>("/api/v1/auth/2fa/rebind/prepare", {
      current_code: currentCode,
    }),

  confirmTotpRebind: (code: string) =>
    apiClient.post<unknown, null>("/api/v1/auth/2fa/rebind/confirm", { code }),

  logout: () => apiClient.post("/api/v1/auth/logout"),

  getMe: () => apiClient.get<unknown, User>("/api/v1/auth/me"),
};
