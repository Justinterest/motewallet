export interface BankAccount {
  id: number;
  currency: string;
  bank_name: string;
  bank_country: string;
  swift_code: string;
  account_name: string;
  account_no_masked: string;
  transfer_type: string;
  status: string;
}

export interface AddBankAccountRequest {
  currency: string;
  bank_name: string;
  bank_country: string;
  swift_code?: string;
  account_name: string;
  account_no: string;
  transfer_type: string;
  account_type?: string;
  payee_country_code?: string;
  payee_address?: string;
  payee_address_second?: string;
  bank_code?: string;
  bank_address?: string;
  middle_swift_code?: string;
}
