export const FIAT_WITHDRAWAL_PURPOSES = [
  { value: "OTHER", label: "其他" },
  { value: "TRAD", label: "服务贸易" },
  { value: "INVS", label: "投资" },
  { value: "GDDS", label: "货物贸易" },
] as const;

export type FiatWithdrawalPurpose = (typeof FIAT_WITHDRAWAL_PURPOSES)[number]["value"];
