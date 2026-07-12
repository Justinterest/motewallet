import apiClient from "./client";
import type { AdminMerchant, AdminMerchantDetail, SyncDepositsResponse, SyncKUNBalancesResponse, WalletLedgerListResponse } from "@/types/merchant";
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
  updateSupportedCurrencies: (
    id: number,
    data: {
      crypto_currencies: string[];
      fiat_currencies: string[];
      crypto_chains: Record<string, string[]>;
      default_chains: Record<string, string>;
    }
  ) => apiClient.put<never, void>(`/api/v1/admin/merchants/${id}/supported-currencies`, data),
  syncKUNBalances: (id: number) =>
    apiClient.post<never, SyncKUNBalancesResponse>(`/api/v1/admin/merchants/${id}/sync-kun-balances`),
  syncDeposits: (id: number) =>
    apiClient.post<never, SyncDepositsResponse>(`/api/v1/admin/merchants/${id}/sync-deposits`),
  getLedger: (
    id: number,
    params?: {
      page?: number;
      page_size?: number;
      account_type?: string;
      currency?: string;
      biz_type?: string;
      entry_type?: string;
    },
  ) =>
    apiClient.get<never, WalletLedgerListResponse>(`/api/v1/admin/merchants/${id}/ledger`, { params }),
};
