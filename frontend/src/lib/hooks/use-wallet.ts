import { useQuery } from "@tanstack/react-query";
import { walletApi } from "@/lib/api/wallet";
import type { WalletLedgerQuery } from "@/types/wallet";

export function useWalletBalances() {
  return useQuery({
    queryKey: ["wallet", "balances"],
    queryFn: () => walletApi.getBalances(),
    staleTime: 30 * 1000,
  });
}

export function useWalletLedger(params?: WalletLedgerQuery, enabled = true) {
  return useQuery({
    queryKey: ["wallet", "ledger", params],
    queryFn: () => walletApi.getLedger(params),
    enabled,
    staleTime: 30 * 1000,
  });
}
