const EXCHANGE_1TO1_PAIRS: Record<string, string[]> = {
  USDT: ["USD"],
  USD: ["USDT", "USDC"],
  USDC: ["USD"],
};

export const EXCHANGE_1TO1_CURRENCIES = ["USD", "USDT", "USDC"] as const;

export function isSupportedExchangePair(fromCurrency: string, toCurrency: string) {
  if (fromCurrency === toCurrency) return false;
  return EXCHANGE_1TO1_PAIRS[fromCurrency]?.includes(toCurrency) ?? false;
}

export function getExchangeToOptions(fromCurrency: string, availableCurrencies: string[]) {
  const allowed = EXCHANGE_1TO1_PAIRS[fromCurrency] ?? [];
  return availableCurrencies.filter((currency) => allowed.includes(currency));
}

export function getExchangeFromOptions(availableCurrencies: string[]) {
  return availableCurrencies.filter((currency) => currency in EXCHANGE_1TO1_PAIRS);
}
