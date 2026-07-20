export const BANK_ACCOUNT_TRANSFER_TYPE = "TT";
export const BANK_ACCOUNT_TYPE = "ENTERPRISE";

export const TRANSFER_TYPE_LABELS: Record<string, string> = {
  TT: "电汇",
};

export function getTransferTypesForCurrency(currency: string) {
  void currency;
  return [{ value: BANK_ACCOUNT_TRANSFER_TYPE, label: "电汇 (TT)" }];
}
