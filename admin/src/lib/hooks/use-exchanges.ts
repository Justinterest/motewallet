import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { exchangeApi, type ListExchangesParams } from "@/lib/api/exchanges";

export function useAdminExchanges(params: ListExchangesParams = {}) {
  return useQuery({
    queryKey: ["admin", "exchanges", params],
    queryFn: () => exchangeApi.list(params),
  });
}

export function useSyncExchangeStatus() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: exchangeApi.syncStatus,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "exchanges"] });
    },
  });
}
