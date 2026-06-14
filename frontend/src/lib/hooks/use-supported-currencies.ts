import { useQuery } from "@tanstack/react-query";
import { currencyApi } from "@/lib/api/currency";

export function useSupportedCurrencies() {
  return useQuery({
    queryKey: ["supported-currencies"],
    queryFn: () => currencyApi.getSupported(),
    staleTime: 5 * 60 * 1000,
  });
}
