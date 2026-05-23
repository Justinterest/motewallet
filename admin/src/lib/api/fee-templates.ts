import apiClient from "./client";
import type { FeeTemplate, CreateFeeTemplateRequest, UpdateFeeTemplateRequest } from "@/types/fee-template";

interface FeeTemplateListResponse {
  templates: FeeTemplate[];
}

export const feeTemplateApi = {
  list: () => apiClient.get<never, FeeTemplateListResponse>("/api/v1/admin/fee-templates").then((res) => res.templates),
  getById: (id: number) => apiClient.get<never, FeeTemplate>(`/api/v1/admin/fee-templates/${id}`),
  create: (data: CreateFeeTemplateRequest) => apiClient.post<never, FeeTemplate>("/api/v1/admin/fee-templates", data),
  update: (id: number, data: UpdateFeeTemplateRequest) => apiClient.put<never, FeeTemplate>(`/api/v1/admin/fee-templates/${id}`, data),
  delete: (id: number) => apiClient.delete<never, void>(`/api/v1/admin/fee-templates/${id}`),
};
