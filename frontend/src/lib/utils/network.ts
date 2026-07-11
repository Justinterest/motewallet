/** Internal chain codes sent to the API (not shown in UI). */
export const DEPOSIT_NETWORKS: Record<string, { value: string; label: string }[]> = {
  USDT: [
    { value: "ETH_ERC20", label: "ERC20（以太坊）" },
    { value: "TRX_TRC20", label: "TRC20（波场）" },
    { value: "TON", label: "TON" },
    { value: "SOL_Solana", label: "Solana" },
    { value: "BSC_BEP20", label: "BEP20（BNB Chain）" },
  ],
  USDC: [
    { value: "ETH_ERC20", label: "ERC20（以太坊）" },
    { value: "TRX_TRC20", label: "TRC20（波场）" },
    { value: "SOL_Solana", label: "Solana" },
    { value: "BSC_BEP20", label: "BEP20（BNB Chain）" },
  ],
  BTC: [{ value: "BTC", label: "Bitcoin" }],
};

const KUN_WITHDRAWAL_CHAIN_TYPES = new Set([
  "ETH_ERC20",
  "TRX_TRC20",
  "SOL_Solana",
  "BSC_BEP20",
  "BTC",
]);

export const WITHDRAWAL_NETWORKS: Record<string, { value: string; label: string }[]> = {
  USDT: (DEPOSIT_NETWORKS.USDT ?? []).filter((item) => KUN_WITHDRAWAL_CHAIN_TYPES.has(item.value)),
  USDC: (DEPOSIT_NETWORKS.USDC ?? []).filter((item) => KUN_WITHDRAWAL_CHAIN_TYPES.has(item.value)),
  BTC: [{ value: "BTC", label: "Bitcoin" }],
};

function isChainSupported(chain: string, supported: string[]) {
  if (supported.includes(chain)) return true;
  // Accept legacy BTC_Bitcoin when BTC is enabled.
  if (chain === "BTC_Bitcoin" && supported.includes("BTC")) return true;
  if (chain === "BTC" && supported.includes("BTC_Bitcoin")) return true;
  return false;
}

export function filterNetworksBySupported(
  networks: { value: string; label: string }[],
  supportedChains?: string[]
) {
  if (!supportedChains || supportedChains.length === 0) {
    return networks;
  }
  return networks.filter((item) => isChainSupported(item.value, supportedChains));
}

export function getDepositNetworks(currency: string, supportedChains?: string[]) {
  return filterNetworksBySupported(DEPOSIT_NETWORKS[currency] ?? [], supportedChains);
}

export function getWithdrawalNetworks(currency: string, supportedChains?: string[]) {
  return filterNetworksBySupported(WITHDRAWAL_NETWORKS[currency] ?? [], supportedChains);
}

export function resolveDefaultChain(
  currency: string,
  networks: { value: string; label: string }[],
  defaultChains?: Record<string, string>
) {
  const preferred = defaultChains?.[currency];
  if (preferred && networks.some((item) => isChainSupported(item.value, [preferred]))) {
    const match = networks.find((item) => isChainSupported(item.value, [preferred]));
    return match?.value || networks[0]?.value || "";
  }
  return networks[0]?.value || "";
}

export function formatChainLabel(currency: string, chain: string) {
  const networks = [...(DEPOSIT_NETWORKS[currency] ?? []), ...(WITHDRAWAL_NETWORKS[currency] ?? [])];
  return networks.find((item) => item.value === chain)?.label ?? chain;
}

const DEPOSIT_STATUS_LABELS: Record<string, string> = {
  COMPLETED: "已到账",
  PROCESSING: "确认中",
  FAILED: "失败",
  SUCCESS: "已到账",
};

export function formatDepositStatus(status: string): string {
  return DEPOSIT_STATUS_LABELS[status] || "处理中";
}
