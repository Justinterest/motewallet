import apiClient from "./client";
import type { AdminMerchant, AdminMerchantDetail, SyncKUNBalancesResponse } from "@/types/merchant";
import type { PaginatedData } from "@/types/api";

export const merchantApi = {
  list: (params: { page?: number; page_size?: number; status?: string; kyc_status?: string; search?: string }) =>
    apiClient.get<never, PaginatedData<AdminMerchant>>("/api/v1/admin/merchants", { params }),
  getById: (id: number) =>
    apiClient.get<never, AdminMerchantDetail>(`/api/v1/admin/merchants/${id}`),
  updateStatus: (id: number, data: { status: string }) =>
    apiClient.put<never, void>(`/api/v1/admin/merchants/${id}/status`, data),
  updateFeeTemplate: (id: number, data: { fee_template_id: number }) =>
    apiClient.put<never, void>(`/api/v1/admin/merchants/${id}/fee-template`, data),
  approveKyc: (id: number) =>
    apiClient.post<never, void>(`/api/v1/admin/merchants/${id}/kyc/approve`),
  rejectKyc: (id: number, data: { reason: string }) =>
    apiClient.post<never, void>(`/api/v1/admin/merchants/${id}/kyc/reject`, data),
  updateSupportedCurrencies: (id: number, data: { crypto_currencies: string[]; fiat_currencies: string[] }) =>
    apiClient.put<never, void>(`/api/v1/admin/merchants/${id}/supported-currencies`, data),
  syncKUNBalances: (id: number) =>
    apiClient.post<never, SyncKUNBalancesResponse>(`/api/v1/admin/merchants/${id}/sync-kun-balances`),
};
