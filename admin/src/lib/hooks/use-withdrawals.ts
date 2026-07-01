import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { withdrawalApi, type ListWithdrawalsParams } from "@/lib/api/withdrawals";

export function useAdminWithdrawals(params: ListWithdrawalsParams) {
  return useQuery({
    queryKey: ["admin", "withdrawals", params],
    queryFn: () => withdrawalApi.list(params),
  });
}

export function useMerchantWithdrawals(merchantId: number, page = 1, pageSize = 10) {
  return useQuery({
    queryKey: ["admin", "withdrawals", merchantId, page, pageSize],
    queryFn: () => withdrawalApi.list({ merchantId, page, pageSize }),
    enabled: !!merchantId,
  });
}

export function usePendingWithdrawals(page = 1) {
  return useQuery({
    queryKey: ["withdrawals", "pending", page],
    queryFn: () => withdrawalApi.listPending(page),
  });
}

export function useApproveWithdrawal() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => withdrawalApi.approve(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["withdrawals"] });
    },
  });
}

export function useRejectWithdrawal() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, reason }: { id: number; reason: string }) =>
      withdrawalApi.reject(id, reason),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["withdrawals"] });
    },
  });
}
