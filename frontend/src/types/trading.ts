export interface DepositAddress {
  address: string;
  currency: string;
  network: string;
}

export interface DepositOrder {
  id: string;
  currency: string;
  network: string;
  amount: string;
  tx_hash?: string;
  status: string;
  created_at: string;
}

export interface DepositOrderListResponse {
  orders: DepositOrder[];
  total: number;
}

export interface WithdrawalOrder {
  id: number;
  type: string;
  currency: string;
  chain?: string;
  amount: string;
  platform_fee: string;
  status: string;
  review_status: string;
  to_address?: string;
  tx_id?: string;
  created_at: string;
}

export interface WithdrawalOrderListResponse {
  orders: WithdrawalOrder[];
  total: number;
}

export interface WithdrawalFeePreview {
  currency: string;
  amount: string;
  platform_fee: string;
  total_deduction: string;
  net_amount: string;
}

export interface ExchangePreview {
  from_currency: string;
  to_currency: string;
  from_amount: string;
  to_amount: string;
  exchange_rate: string;
  quote_id: string;
  expire_time: number;
  kun_trade_fee: string;
  kun_fee_currency: string;
  platform_fee: string;
  fee_currency: string;
  net_to_amount: string;
  total_deduction: string;
}

export interface ExchangeOrder {
  id: number;
  exchange_type: string;
  from_currency: string;
  to_currency: string;
  from_amount: string;
  to_amount: string;
  exchange_rate: string;
  platform_fee: string;
  status: string;
  fail_reason?: string;
  created_at: string;
}

export interface ExchangeOrderListResponse {
  orders: ExchangeOrder[];
  total: number;
}

export interface TransferOrder {
  id: number;
  from_account_type: string;
  to_account_type: string;
  currency: string;
  amount: string;
  status: string;
  created_at: string;
}

export interface TransferOrderListResponse {
  orders: TransferOrder[];
  total: number;
}
