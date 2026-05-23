import apiClient from "./client";
import type { LoginRequest, RegisterRequest, User } from "@/types/auth";

export const authApi = {
  login: (data: LoginRequest) =>
    apiClient.post<unknown, User>("/api/v1/auth/login", data),

  register: (data: RegisterRequest) =>
    apiClient.post<unknown, User>("/api/v1/auth/register", data),

  logout: () => apiClient.post("/api/v1/auth/logout"),

  getMe: () => apiClient.get<unknown, User>("/api/v1/auth/me"),
};
