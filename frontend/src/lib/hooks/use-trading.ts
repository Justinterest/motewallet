import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { tradingApi } from "@/lib/api/trading";

export function useDepositAddresses(currency: string, chain: string) {
  return useQuery({
    queryKey: ["deposit", "addresses", currency, chain],
    queryFn: () => tradingApi.getDepositAddresses(currency, chain),
    enabled: !!currency && !!chain,
  });
}

export function useDepositOrders(currency?: string, chain?: string, page = 1) {
  return useQuery({
    queryKey: ["deposit", "orders", currency, chain, page],
    queryFn: () => tradingApi.listDepositOrders({ page, currency, chain }),
    refetchInterval: 30_000,
  });
}

export function useSubmitCryptoWithdrawal() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: tradingApi.submitCryptoWithdrawal,
    onSuccess: () => { qc.invalidateQueries({ queryKey: ["withdrawal"] }); qc.invalidateQueries({ queryKey: ["wallet"] }); },
  });
}

export function useSubmitFiatWithdrawal() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: tradingApi.submitFiatWithdrawal,
    onSuccess: () => { qc.invalidateQueries({ queryKey: ["withdrawal"] }); qc.invalidateQueries({ queryKey: ["wallet"] }); },
  });
}

export function useWithdrawalOrders(page = 1) {
  return useQuery({
    queryKey: ["withdrawal", "orders", page],
    queryFn: () => tradingApi.listWithdrawalOrders(page),
  });
}

export function useExchangeQuote() {
  return useMutation({
    mutationFn: tradingApi.getQuote,
  });
}

export function useCreateExchangeOrder() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: tradingApi.createExchangeOrder,
    onSuccess: () => { qc.invalidateQueries({ queryKey: ["exchange"] }); qc.invalidateQueries({ queryKey: ["wallet"] }); },
  });
}

export function useCreate1to1Order() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: tradingApi.create1to1Order,
    onSuccess: () => { qc.invalidateQueries({ queryKey: ["exchange"] }); qc.invalidateQueries({ queryKey: ["wallet"] }); },
  });
}

export function useExchangeOrders(page = 1) {
  return useQuery({
    queryKey: ["exchange", "orders", page],
    queryFn: () => tradingApi.listExchangeOrders(page),
  });
}

export function useTransfer() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: tradingApi.transfer,
    onSuccess: () => { qc.invalidateQueries({ queryKey: ["transfer"] }); qc.invalidateQueries({ queryKey: ["wallet"] }); },
  });
}

export function useTransferOrders(page = 1) {
  return useQuery({
    queryKey: ["transfer", "orders", page],
    queryFn: () => tradingApi.listTransfers(page),
  });
}
