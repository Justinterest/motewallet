export interface AdminExchange {
  id: number;
  platform_order_id: string;
  merchant_id: number;
  merchant_email: string;
  exchange_type: string;
  from_currency: string;
  to_currency: string;
  from_amount: string;
  to_amount?: string;
  exchange_rate?: string;
  platform_fee: string;
  status: string;
  fail_reason?: string;
  kun_order_id?: string;
  kun_request_no?: string;
  created_at: string;
  completed_at?: string;
}

export interface AdminExchangeListResponse {
  exchanges: AdminExchange[];
  total: number;
}

export interface AdminExchangeSyncResponse {
  order_id: number;
  status: string;
  kun_status: string;
  updated: boolean;
  to_amount?: string;
  exchange_rate?: string;
  platform_fee: string;
  fail_reason?: string;
}
