import apiClient from "./client";
import type {
  DepositAddress,
  DepositOrderListResponse,
  WithdrawalFeePreview,
  WithdrawalOrderListResponse,
  ExchangePreview,
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
  submitCryptoWithdrawal: (data: { currency: string; crypto_address_id: number; amount: string }) =>
    apiClient.post<never, { order_id: number }>("/api/v1/withdraw/crypto", data),
  submitFiatWithdrawal: (data: {
    currency: string;
    amount: string;
    bank_account_id: number;
    purpose: string;
    postscript: string;
  }) =>
    apiClient.post<never, { order_id: number }>("/api/v1/withdraw/fiat", data),
  previewWithdrawalFee: (data: {
    type: "CRYPTO" | "FIAT";
    currency: string;
    amount: string;
    crypto_address_id?: number;
    bank_account_id?: number;
  }) => apiClient.post<never, WithdrawalFeePreview>("/api/v1/withdraw/fee-preview", data),
  listWithdrawalOrders: (page = 1, pageSize = 20) =>
    apiClient.get<never, WithdrawalOrderListResponse>("/api/v1/withdraw/orders", { params: { page, page_size: pageSize } }),

  // Exchange (spot)
  previewExchange: (data: { from_currency: string; to_currency: string; from_amount: string }) =>
    apiClient.post<never, ExchangePreview>("/api/v1/exchange/preview", data),
  createExchangeOrder: (data: {
    from_currency: string;
    to_currency: string;
    from_amount: string;
    quote_id: string;
  }) =>
    apiClient.post<never, { order_id: number }>("/api/v1/exchange/order", data),
  listExchangeOrders: (page = 1, pageSize = 20) =>
    apiClient.get<never, ExchangeOrderListResponse>("/api/v1/exchange/orders", { params: { page, page_size: pageSize } }),

  // Transfer
  transfer: (data: { from_account_type: string; to_account_type: string; currency: string; amount: string }) =>
    apiClient.post<never, { order_id: number }>("/api/v1/account/transfer", data),
  listTransfers: (page = 1, pageSize = 20) =>
    apiClient.get<never, TransferOrderListResponse>("/api/v1/account/transfers", { params: { page, page_size: pageSize } }),
};
