export function formatAmount(amount: string, currency?: string): string {
  const num = parseFloat(amount);
  if (isNaN(num)) return "0.00";
  const decimals = currency === "BTC" ? 8 : 2;
  return num.toLocaleString("en-US", {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  });
}
