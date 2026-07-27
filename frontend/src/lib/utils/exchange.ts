const EXCHANGE_1TO1_PAIRS: Record<string, string[]> = {
  USDT: ["USD"],
  USD: ["USDT", "USDC"],
  USDC: ["USD"],
};

export const EXCHANGE_1TO1_CURRENCIES = ["USD", "USDT", "USDC"] as const;

const EXCHANGE_STATUS_LABELS: Record<string, string> = {
  COMPLETED: "已完成",
  FAILED: "失败",
  PROCESSING: "处理中",
  PENDING: "待处理",
};

const EXCHANGE_STATUS_CLASS: Record<string, string> = {
  COMPLETED: "text-green-600",
  FAILED: "text-red-600",
  PROCESSING: "text-amber-600",
  PENDING: "text-slate-500",
};

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

export function formatExchangeStatus(status: string) {
  return EXCHANGE_STATUS_LABELS[status] || status;
}

export function getExchangeStatusClassName(status: string) {
  return EXCHANGE_STATUS_CLASS[status] || EXCHANGE_STATUS_CLASS.PENDING;
}

function hasPositiveAmount(value: string) {
  const parsed = parseFloat(value);
  return !Number.isNaN(parsed) && parsed > 0;
}

export function formatExchangeOrderHint(order: { status: string; to_amount?: string }) {
  if (order.status === "COMPLETED") return null;
  if (order.status === "FAILED") return "兑换未成功，资金已退回";
  return "到账金额确认中";
}

export function formatExchangeOrderMeta(
  order: {
    from_currency: string;
    to_currency: string;
    exchange_rate?: string;
    platform_fee?: string;
    fee_deduction_method?: "WALLET" | "RECEIVED_AMOUNT";
    created_at: string;
    status: string;
  },
  formatAmount: (amount: string, currency?: string) => string,
) {
  const parts: string[] = [];

  if (order.status === "COMPLETED" && order.exchange_rate) {
    parts.push(
      `汇率 1 ${order.from_currency} = ${formatAmount(order.exchange_rate, order.to_currency)} ${order.to_currency}`,
    );
  }

  if (order.platform_fee && hasPositiveAmount(order.platform_fee)) {
    const feeCurrency =
      order.fee_deduction_method === "RECEIVED_AMOUNT"
        ? order.to_currency
        : order.from_currency;
    parts.push(
      `手续费 ${formatAmount(order.platform_fee, feeCurrency)} ${feeCurrency}`,
    );
  }

  parts.push(new Date(order.created_at).toLocaleString("zh-CN"));
  return parts.join(" · ");
}
