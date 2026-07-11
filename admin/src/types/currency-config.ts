export interface SystemCurrencyConfig {
  crypto_currencies: string[];
  fiat_currencies: string[];
  crypto_chains: Record<string, string[]>;
  default_chains: Record<string, string>;
  catalog_chains: Record<string, string[]>;
  all_crypto: string[];
  all_fiat: string[];
}

export interface UpdateSystemCurrencyConfigPayload {
  crypto_currencies: string[];
  fiat_currencies: string[];
  crypto_chains: Record<string, string[]>;
  default_chains: Record<string, string>;
}
