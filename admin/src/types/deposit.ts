export interface AdminDeposit {
  id: number;
  platform_order_id: string;
  merchant_id: number;
  merchant_email: string;
  currency: string;
  network: string;
  amount: string;
  tx_hash?: string;
  to_address: string;
  from_address?: string;
  status: string;
  created_at: string;
  completed_at?: string;
}

export interface AdminDepositListResponse {
  deposits: AdminDeposit[];
  total: number;
}
