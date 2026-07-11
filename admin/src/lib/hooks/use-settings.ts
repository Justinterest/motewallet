import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { settingsApi } from "@/lib/api/settings";
import type { UpdateSystemCurrencyConfigPayload } from "@/types/currency-config";

export function useSystemCurrencyConfig() {
  return useQuery({
    queryKey: ["system-currency-config"],
    queryFn: () => settingsApi.getCurrencyConfig(),
  });
}

export function useUpdateSystemCurrencyConfig() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: UpdateSystemCurrencyConfigPayload) => settingsApi.updateCurrencyConfig(data),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["system-currency-config"] }),
  });
}
