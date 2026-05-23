import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { withdrawalApi } from "@/lib/api/withdrawals";

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
