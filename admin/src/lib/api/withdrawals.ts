import apiClient from "./client";
import type { AdminWithdrawalListResponse } from "@/types/withdrawal";

export interface WithdrawalOrder {
  id: number;
  type: string;
  currency: string;
  chain?: string;
  amount: string;
  platform_fee: string;
  status: string;
  review_status: string;
  to_address?: string;
  tx_id?: string;
  created_at: string;
}

export interface WithdrawalListResponse {
  orders: WithdrawalOrder[];
  total: number;
}

export interface ListWithdrawalsParams {
  page?: number;
  pageSize?: number;
  merchantId?: number;
  merchantEmail?: string;
  currency?: string;
  status?: string;
  reviewStatus?: string;
  type?: string;
}

export const withdrawalApi = {
  list: (params: ListWithdrawalsParams = {}) => {
    const {
      page = 1,
      pageSize = 20,
      merchantId,
      merchantEmail,
      currency,
      status,
      reviewStatus,
      type,
    } = params;
    return apiClient.get<never, AdminWithdrawalListResponse>("/api/v1/admin/withdrawals", {
      params: {
        page,
        page_size: pageSize,
        merchant_id: merchantId || undefined,
        merchant_email: merchantEmail || undefined,
        currency: currency && currency !== "ALL" ? currency : undefined,
        status: status && status !== "ALL" ? status : undefined,
        review_status: reviewStatus && reviewStatus !== "ALL" ? reviewStatus : undefined,
        type: type && type !== "ALL" ? type : undefined,
      },
    });
  },
  listPending: (page = 1, pageSize = 20) =>
    apiClient.get<never, WithdrawalListResponse>("/api/v1/admin/withdrawals/pending", {
      params: { page, page_size: pageSize },
    }),
  approve: (id: number) =>
    apiClient.post(`/api/v1/admin/withdrawals/${id}/approve`),
  reject: (id: number, reason: string) =>
    apiClient.post(`/api/v1/admin/withdrawals/${id}/reject`, { reason }),
};
