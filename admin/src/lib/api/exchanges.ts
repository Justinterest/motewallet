import apiClient from "./client";
import type { AdminExchangeListResponse, AdminExchangeSyncResponse } from "@/types/exchange";

export interface ListExchangesParams {
  page?: number;
  pageSize?: number;
  merchantId?: number;
  merchantEmail?: string;
  currency?: string;
  status?: string;
}

export const exchangeApi = {
  list: (params: ListExchangesParams = {}) => {
    const { page = 1, pageSize = 20, merchantId, merchantEmail, currency, status } = params;
    return apiClient.get<never, AdminExchangeListResponse>("/api/v1/admin/exchanges", {
      params: {
        page,
        page_size: pageSize,
        merchant_id: merchantId || undefined,
        merchant_email: merchantEmail || undefined,
        currency: currency && currency !== "ALL" ? currency : undefined,
        status: status && status !== "ALL" ? status : undefined,
      },
    });
  },
  syncStatus: (id: number) =>
    apiClient.post<never, AdminExchangeSyncResponse>(`/api/v1/admin/exchanges/${id}/sync`),
};
