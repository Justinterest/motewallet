export const TRANSFER_TYPE_OPTIONS: Record<string, { value: string; label: string }[]> = {
  USD: [
    { value: "LOCAL", label: "本地转账" },
    { value: "TT", label: "电汇 (TT)" },
  ],
  HKD: [
    { value: "LOCAL", label: "本地转账" },
    { value: "CHATS", label: "CHATS 即时结算" },
    { value: "TT", label: "电汇 (TT)" },
  ],
  EUR: [
    { value: "LOCAL", label: "本地转账" },
    { value: "TT", label: "电汇 (TT)" },
  ],
};

export const TRANSFER_TYPE_LABELS: Record<string, string> = {
  LOCAL: "本地转账",
  CHATS: "CHATS",
  TT: "电汇",
};

export function getTransferTypesForCurrency(currency: string) {
  return TRANSFER_TYPE_OPTIONS[currency] ?? [{ value: "LOCAL", label: "本地转账" }];
}
