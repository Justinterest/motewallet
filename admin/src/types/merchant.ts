export interface MerchantWallet {
  account_type: string;
  currency: string;
  balance: string;
  frozen_balance: string;
  available_balance: string;
}

export interface AdminMerchant {
  id: number;
  email: string;
  status: string;
  kyc_status: string;
  fee_template_id?: number;
  fee_template_name?: string;
  kyc_submitted_at?: string;
  kyc_completed_at?: string;
  agreement_signed_at?: string;
  frozen_at?: string;
  created_at: string;
}

export interface AdminMerchantDetail extends AdminMerchant {
  wallets: MerchantWallet[];
  kyc_fail_reason?: string;
  supported_crypto_currencies: string[];
  supported_fiat_currencies: string[];
  available_crypto_currencies: string[];
  available_fiat_currencies: string[];
}
