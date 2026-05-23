export function formatAmount(amount: string, currency?: string): string {
  const num = parseFloat(amount);
  if (isNaN(num)) return "0.00";

  const isCrypto = ["USDT", "USDC", "BTC"].includes(currency || "");
  const decimals = currency === "BTC" ? 8 : isCrypto ? 2 : 2;

  return num.toLocaleString("en-US", {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  });
}

export function getCurrencySymbol(currency: string): string {
  const symbols: Record<string, string> = {
    USD: "$",
    HKD: "HK$",
    EUR: "€",
    USDT: "USDT",
    USDC: "USDC",
    BTC: "₿",
  };
  return symbols[currency] || currency;
}
