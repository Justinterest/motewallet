import apiClient from "./client";
import type {
  LoginRequest,
  RegisterRequest,
  SendVerificationCodeRequest,
  User,
} from "@/types/auth";

export const authApi = {
  login: (data: LoginRequest) =>
    apiClient.post<unknown, User>("/api/v1/auth/login", data),

  sendVerificationCode: (data: SendVerificationCodeRequest) =>
    apiClient.post<unknown, null>("/api/v1/auth/send-verification-code", data),

  register: (data: RegisterRequest) =>
    apiClient.post<unknown, User>("/api/v1/auth/register", data),

  logout: () => apiClient.post("/api/v1/auth/logout"),

  getMe: () => apiClient.get<unknown, User>("/api/v1/auth/me"),
};
