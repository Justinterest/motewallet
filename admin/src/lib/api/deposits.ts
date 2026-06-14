import apiClient from "./client";
import type { AdminDepositListResponse } from "@/types/deposit";

export interface ListDepositsParams {
  page?: number;
  pageSize?: number;
  merchantId?: number;
  merchantEmail?: string;
  currency?: string;
  status?: string;
}

export const depositApi = {
  list: (params: ListDepositsParams = {}) => {
    const { page = 1, pageSize = 20, merchantId, merchantEmail, currency, status } = params;
    return apiClient.get<never, AdminDepositListResponse>("/api/v1/admin/deposits", {
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
};
