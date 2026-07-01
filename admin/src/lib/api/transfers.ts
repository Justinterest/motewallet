import apiClient from "./client";
import type { AdminTransferListResponse, AdminTransferSyncResponse } from "@/types/transfer";

export interface ListTransfersParams {
  page?: number;
  pageSize?: number;
  merchantId?: number;
  merchantEmail?: string;
  currency?: string;
  status?: string;
}

export const transferApi = {
  list: (params: ListTransfersParams = {}) => {
    const { page = 1, pageSize = 20, merchantId, merchantEmail, currency, status } = params;
    return apiClient.get<never, AdminTransferListResponse>("/api/v1/admin/transfers", {
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
    apiClient.post<never, AdminTransferSyncResponse>(`/api/v1/admin/transfers/${id}/sync`),
};
