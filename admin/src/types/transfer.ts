export interface AdminTransfer {
  id: number;
  platform_order_id: string;
  merchant_id: number;
  merchant_email: string;
  from_account_type: string;
  to_account_type: string;
  currency: string;
  amount: string;
  status: string;
  kun_order_id?: string;
  kun_request_no?: string;
  created_at: string;
}

export interface AdminTransferListResponse {
  transfers: AdminTransfer[];
  total: number;
}

export interface AdminTransferSyncResponse {
  order_id: number;
  status: string;
  kun_status: string;
  updated: boolean;
}
