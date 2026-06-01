import apiClient from "./client";
import type {
  DepositAddress,
  DepositOrderListResponse,
  WithdrawalOrderListResponse,
  ExchangeQuote,
  ExchangeOrderListResponse,
  TransferOrderListResponse,
} from "@/types/trading";

export const tradingApi = {
  // Deposit
  getDepositAddresses: (currency: string, chain: string) =>
    apiClient.get<never, DepositAddress>("/api/v1/deposit/addresses", { params: { currency, chain } }),
  listDepositOrders: (params?: { page?: number; pageSize?: number; currency?: string; chain?: string }) => {
    const { page = 1, pageSize = 20, currency, chain } = params ?? {};
    return apiClient.get<never, DepositOrderListResponse>("/api/v1/deposit/orders", {
      params: { page, page_size: pageSize, currency, chain },
    });
  },

  // Withdrawal
  submitCryptoWithdrawal: (data: { currency: string; chain: string; amount: string; to_address: string }) =>
    apiClient.post<never, { order_id: number }>("/api/v1/withdraw/crypto", data),
  submitFiatWithdrawal: (data: { currency: string; amount: string; bank_account_id: number }) =>
    apiClient.post<never, { order_id: number }>("/api/v1/withdraw/fiat", data),
  listWithdrawalOrders: (page = 1, pageSize = 20) =>
    apiClient.get<never, WithdrawalOrderListResponse>("/api/v1/withdraw/orders", { params: { page, page_size: pageSize } }),

  // Exchange
  getQuote: (data: { from_currency: string; to_currency: string; from_amount: string }) =>
    apiClient.post<never, ExchangeQuote>("/api/v1/exchange/quote", data),
  createExchangeOrder: (data: { quote_id: string; from_currency: string; to_currency: string; from_amount: string }) =>
    apiClient.post<never, { order_id: number }>("/api/v1/exchange/order", data),
  create1to1Order: (data: { from_currency: string; to_currency: string; from_amount: string }) =>
    apiClient.post<never, { order_id: number }>("/api/v1/exchange/1to1", data),
  listExchangeOrders: (page = 1, pageSize = 20) =>
    apiClient.get<never, ExchangeOrderListResponse>("/api/v1/exchange/orders", { params: { page, page_size: pageSize } }),

  // Transfer
  transfer: (data: { from_account_type: string; to_account_type: string; currency: string; amount: string }) =>
    apiClient.post<never, { order_id: number }>("/api/v1/account/transfer", data),
  listTransfers: (page = 1, pageSize = 20) =>
    apiClient.get<never, TransferOrderListResponse>("/api/v1/account/transfers", { params: { page, page_size: pageSize } }),
};
