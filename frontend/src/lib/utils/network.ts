/** Internal chain codes sent to the API (not shown in UI). */
export const DEPOSIT_NETWORKS: Record<string, { value: string; label: string }[]> = {
  USDT: [
    { value: "TRX_TRC20", label: "TRC20（波场）" },
    { value: "ETH_ERC20", label: "ERC20（以太坊）" },
  ],
  USDC: [
    { value: "ETH_ERC20", label: "ERC20（以太坊）" },
    { value: "TRX_TRC20", label: "TRC20（波场）" },
  ],
  BTC: [{ value: "BTC_Bitcoin", label: "Bitcoin" }],
};

const DEPOSIT_STATUS_LABELS: Record<string, string> = {
  COMPLETED: "已到账",
  PROCESSING: "确认中",
  FAILED: "失败",
  SUCCESS: "已到账",
};

export function formatDepositStatus(status: string): string {
  return DEPOSIT_STATUS_LABELS[status] || "处理中";
}
