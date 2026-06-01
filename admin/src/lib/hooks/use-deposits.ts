import { useQuery } from "@tanstack/react-query";
import { depositApi, type ListDepositsParams } from "@/lib/api/deposits";

export function useAdminDeposits(params: ListDepositsParams) {
  return useQuery({
    queryKey: ["admin", "deposits", params],
    queryFn: () => depositApi.list(params),
  });
}
