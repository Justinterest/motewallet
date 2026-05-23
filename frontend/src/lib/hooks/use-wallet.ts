import { useQuery } from "@tanstack/react-query";
import { walletApi } from "@/lib/api/wallet";

export function useWalletBalances() {
  return useQuery({
    queryKey: ["wallet", "balances"],
    queryFn: () => walletApi.getBalances(),
    staleTime: 30 * 1000,
  });
}
