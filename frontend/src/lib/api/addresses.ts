import apiClient from "./client";
import type { AddBankAccountRequest, AddCryptoAddressRequest, BankAccount, CryptoAddress } from "@/types/address";

export const addressesApi = {
  listBankAccounts: () => apiClient.get<never, BankAccount[]>("/api/v1/addresses/bank"),
  addBankAccount: (data: AddBankAccountRequest) =>
    apiClient.post<never, BankAccount>("/api/v1/addresses/bank", data),
  deleteBankAccount: (id: number) => apiClient.delete<never, void>(`/api/v1/addresses/bank/${id}`),

  listCryptoAddresses: () => apiClient.get<never, CryptoAddress[]>("/api/v1/addresses/crypto"),
  addCryptoAddress: (data: AddCryptoAddressRequest) =>
    apiClient.post<never, CryptoAddress>("/api/v1/addresses/crypto", data),
  deleteCryptoAddress: (id: number) => apiClient.delete<never, void>(`/api/v1/addresses/crypto/${id}`),
};
