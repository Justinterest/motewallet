import apiClient from "./client";
import type { AdminUser, AdminLoginRequest } from "@/types/auth";

export const adminAuthApi = {
  login: (data: AdminLoginRequest) =>
    apiClient.post<unknown, AdminUser>("/api/v1/admin/auth/login", data),

  logout: () => apiClient.post("/api/v1/admin/auth/logout"),

  getMe: () => apiClient.get<unknown, AdminUser>("/api/v1/admin/auth/me"),
};
