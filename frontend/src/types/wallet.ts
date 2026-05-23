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
