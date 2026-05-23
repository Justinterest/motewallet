import apiClient from "./client";

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

export const withdrawalApi = {
  listPending: (page = 1, pageSize = 20) =>
    apiClient.get<never, WithdrawalListResponse>("/api/v1/admin/withdrawals/pending", {
      params: { page, page_size: pageSize },
    }),
  approve: (id: number) =>
    apiClient.post(`/api/v1/admin/withdrawals/${id}/approve`),
  reject: (id: number, reason: string) =>
    apiClient.post(`/api/v1/admin/withdrawals/${id}/reject`, { reason }),
};
