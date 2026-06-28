import apiClient from "./client";
import type { AddBankAccountRequest, BankAccount } from "@/types/address";

export const addressesApi = {
  listBankAccounts: () => apiClient.get<never, BankAccount[]>("/api/v1/addresses/bank"),
  addBankAccount: (data: AddBankAccountRequest) =>
    apiClient.post<never, BankAccount>("/api/v1/addresses/bank", data),
  deleteBankAccount: (id: number) => apiClient.delete<never, void>(`/api/v1/addresses/bank/${id}`),
};
