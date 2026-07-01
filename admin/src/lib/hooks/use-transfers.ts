import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { transferApi, type ListTransfersParams } from "@/lib/api/transfers";

export function useAdminTransfers(params: ListTransfersParams = {}) {
  return useQuery({
    queryKey: ["admin", "transfers", params],
    queryFn: () => transferApi.list(params),
  });
}

export function useSyncTransferStatus() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: transferApi.syncStatus,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "transfers"] });
    },
  });
}
