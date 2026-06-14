import { useQuery } from "@tanstack/react-query";
import { depositApi, type ListDepositsParams } from "@/lib/api/deposits";

export function useAdminDeposits(params: ListDepositsParams) {
  return useQuery({
    queryKey: ["admin", "deposits", params],
    queryFn: () => depositApi.list(params),
  });
}

export function useMerchantDeposits(merchantId: number, page = 1, pageSize = 10) {
  return useQuery({
    queryKey: ["admin", "deposits", merchantId, page, pageSize],
    queryFn: () => depositApi.list({ merchantId, page, pageSize }),
    enabled: !!merchantId,
  });
}
