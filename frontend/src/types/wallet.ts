export interface WalletBalance {
  account_type: string;  // FUNDING or TRADING
  currency: string;      // USDT, USDC, BTC, USD, HKD, EUR
  balance: string;       // decimal string
  frozen_balance: string;
  available_balance: string;
}

export interface WalletBalancesResponse {
  wallets: WalletBalance[];
}

export type WalletLedgerEntryType = "CREDIT" | "FREEZE" | "UNFREEZE" | "DEDUCT_FROZEN";

export interface WalletLedgerEntry {
  id: number;
  account_type: string;
  currency: string;
  entry_type: WalletLedgerEntryType | string;
  amount: string;
  balance_before: string;
  balance_after: string;
  frozen_before: string;
  frozen_after: string;
  transaction_record_id?: number;
  platform_order_id?: string;
  biz_type?: string;
  remark?: string;
  created_at: string;
}

export interface WalletLedgerListResponse {
  entries: WalletLedgerEntry[];
  total: number;
  page: number;
  page_size: number;
}

export interface WalletLedgerQuery {
  page?: number;
  page_size?: number;
  account_type?: string;
  currency?: string;
  biz_type?: string;
  entry_type?: string;
}
