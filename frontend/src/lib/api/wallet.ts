import apiClient from "./client";
import type {
  WalletBalancesResponse,
  WalletLedgerListResponse,
  WalletLedgerQuery,
} from "@/types/wallet";

export const walletApi = {
  getBalances: () => apiClient.get<never, WalletBalancesResponse>("/api/v1/account/balances"),
  getLedger: (params?: WalletLedgerQuery) =>
    apiClient.get<never, WalletLedgerListResponse>("/api/v1/account/ledger", { params }),
};
