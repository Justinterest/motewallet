import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { merchantApi } from "@/lib/api/merchants";

export function useMerchants(params: { page?: number; page_size?: number; status?: string; kyc_status?: string; search?: string }) {
  return useQuery({
    queryKey: ["merchants", params],
    queryFn: () => merchantApi.list(params),
  });
}

export function useMerchant(id: number) {
  return useQuery({
    queryKey: ["merchants", id],
    queryFn: () => merchantApi.getById(id),
    enabled: !!id,
  });
}

export function useUpdateMerchantStatus() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status }: { id: number; status: string }) =>
      merchantApi.updateStatus(id, { status }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["merchants"] }),
  });
}

export function useUpdateMerchantFeeTemplate() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, fee_template_id }: { id: number; fee_template_id: number }) =>
      merchantApi.updateFeeTemplate(id, { fee_template_id }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["merchants"] }),
  });
}

export function useApproveKyc() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => merchantApi.approveKyc(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["merchants"] }),
  });
}

export function useRejectKyc() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, reason }: { id: number; reason: string }) =>
      merchantApi.rejectKyc(id, { reason }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["merchants"] }),
  });
}

export function useUpdateMerchantSupportedCurrencies() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      crypto_currencies,
      fiat_currencies,
    }: {
      id: number;
      crypto_currencies: string[];
      fiat_currencies: string[];
    }) => merchantApi.updateSupportedCurrencies(id, { crypto_currencies, fiat_currencies }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["merchants"] }),
  });
}
