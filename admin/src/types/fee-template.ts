export interface FeeTemplateExchangeItem {
  id?: number;
  from_currency: string;
  to_currency: string;
  fee_rate: string;
  min_fee: string;
  min_fee_currency: string;
}

export interface FeeTemplateCryptoWithdrawalItem {
  id?: number;
  currency: string;
  chain: string;
  fee_rate: string;
  fixed_fee: string;
}

export interface FeeTemplateFiatWithdrawalItem {
  id?: number;
  currency: string;
  transfer_type: string;
  fee_rate: string;
  fixed_fee: string;
}

export interface FeeTemplate {
  id: number;
  name: string;
  description?: string;
  is_default: boolean;
  exchange_items?: FeeTemplateExchangeItem[];
  crypto_withdrawal_items?: FeeTemplateCryptoWithdrawalItem[];
  fiat_withdrawal_items?: FeeTemplateFiatWithdrawalItem[];
  created_at: string;
  updated_at: string;
}

export interface CreateFeeTemplateRequest {
  name: string;
  description?: string;
  is_default?: boolean;
  exchange_items?: Omit<FeeTemplateExchangeItem, "id">[];
  crypto_withdrawal_items?: Omit<FeeTemplateCryptoWithdrawalItem, "id">[];
  fiat_withdrawal_items?: Omit<FeeTemplateFiatWithdrawalItem, "id">[];
}

export type UpdateFeeTemplateRequest = CreateFeeTemplateRequest;
