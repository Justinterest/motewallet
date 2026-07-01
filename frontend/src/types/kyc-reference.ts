export interface CountryOption {
  country_code: string;
  country_name: string;
}

export interface CountryListResponse {
  items: CountryOption[];
}

export interface AuthTypeOption {
  doc_code: string;
  doc_name: string;
}

export interface CountryAuthTypesResponse {
  items: AuthTypeOption[];
}

export type KycCountryScene =
  | "REGISTER"
  | "REGISTER_ADDRESS"
  | "WITHDRAWAL"
  | "BIND_ACCOUNT";
