export function toCurrencyOptions(currencies: string[]) {
  return currencies.map((currency) => ({ value: currency, label: currency }));
}

export function getAllSupportedCurrencies(data?: {
  crypto_currencies?: string[];
  fiat_currencies?: string[];
}) {
  return [...(data?.crypto_currencies ?? []), ...(data?.fiat_currencies ?? [])];
}
