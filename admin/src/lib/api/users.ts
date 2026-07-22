import apiClient from "./client";
import type {
  AdminEmployee,
  CreateAdminEmployeeRequest,
  CreateAdminEmployeeResponse,
  ResetAdminPasswordResponse,
} from "@/types/auth";

export const adminUsersApi = {
  list: () => apiClient.get<unknown, AdminEmployee[]>("/api/v1/admin/users"),

  create: (data: CreateAdminEmployeeRequest) =>
    apiClient.post<unknown, CreateAdminEmployeeResponse>(
      "/api/v1/admin/users",
      data
    ),

  resetPassword: (id: number) =>
    apiClient.post<unknown, ResetAdminPasswordResponse>(
      `/api/v1/admin/users/${id}/reset-password`
    ),

  reset2FA: (id: number) =>
    apiClient.post<never, void>(`/api/v1/admin/users/${id}/reset-2fa`),
};
