import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { addressesApi } from "@/lib/api/addresses";
import type { AddBankAccountRequest, AddCryptoAddressRequest } from "@/types/address";

export function useBankAccounts() {
  return useQuery({
    queryKey: ["bank-accounts"],
    queryFn: async () => (await addressesApi.listBankAccounts()) ?? [],
  });
}

export function useAddBankAccount() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: AddBankAccountRequest) => addressesApi.addBankAccount(data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["bank-accounts"] });
    },
  });
}

export function useDeleteBankAccount() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => addressesApi.deleteBankAccount(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["bank-accounts"] });
    },
  });
}

export function useCryptoAddresses() {
  return useQuery({
    queryKey: ["crypto-addresses"],
    queryFn: async () => (await addressesApi.listCryptoAddresses()) ?? [],
  });
}

export function useAddCryptoAddress() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: AddCryptoAddressRequest) => addressesApi.addCryptoAddress(data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["crypto-addresses"] });
    },
  });
}

export function useDeleteCryptoAddress() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => addressesApi.deleteCryptoAddress(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["crypto-addresses"] });
    },
  });
}
