export interface SupportedCurrencies {
  crypto_currencies: string[];
  fiat_currencies: string[];
  crypto_chains: Record<string, string[]>;
  default_chains: Record<string, string>;
}
