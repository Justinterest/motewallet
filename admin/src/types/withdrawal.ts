export interface AdminWithdrawal {
  id: number;
  platform_order_id: string;
  merchant_id: number;
  merchant_email: string;
  type: string;
  currency: string;
  network?: string;
  amount: string;
  platform_fee: string;
  status: string;
  review_status: string;
  to_address?: string;
  tx_id?: string;
  created_at: string;
  completed_at?: string;
}

export interface AdminWithdrawalListResponse {
  withdrawals: AdminWithdrawal[];
  total: number;
}
