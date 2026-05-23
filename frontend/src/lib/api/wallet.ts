import apiClient from "./client";
import type { WalletBalancesResponse } from "@/types/wallet";

export const walletApi = {
  getBalances: () => apiClient.get<never, WalletBalancesResponse>("/api/v1/account/balances"),
};
