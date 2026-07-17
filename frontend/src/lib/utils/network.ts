/** Standard display labels for internal chain codes (API values stay unchanged). */
const CHAIN_LABELS: Record<string, string> = {
  ETH_ERC20: "ERC20（以太坊）",
  TRX_TRC20: "TRC20（波场）",
  TON: "TON",
  SOL_Solana: "Solana",
  BSC_BEP20: "BEP20（BNB Chain）",
  BTC: "Bitcoin",
  BTC_Bitcoin: "Bitcoin",
};

function chainOption(value: string) {
  return { value, label: CHAIN_LABELS[value] || "未知网络" };
}

/** Internal chain codes sent to the API (not shown in UI). */
export const DEPOSIT_NETWORKS: Record<string, { value: string; label: string }[]> = {
  USDT: [
    chainOption("ETH_ERC20"),
    chainOption("TRX_TRC20"),
    chainOption("TON"),
    chainOption("SOL_Solana"),
    chainOption("BSC_BEP20"),
  ],
  USDC: [
    chainOption("ETH_ERC20"),
    chainOption("TRX_TRC20"),
    chainOption("SOL_Solana"),
    chainOption("BSC_BEP20"),
  ],
  BTC: [chainOption("BTC")],
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
  BTC: [chainOption("BTC")],
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

export function formatChainLabel(currency: string, chain?: string) {
  if (!chain) return "";
  const networks = [...(DEPOSIT_NETWORKS[currency] ?? []), ...(WITHDRAWAL_NETWORKS[currency] ?? [])];
  return (
    networks.find((item) => item.value === chain)?.label ||
    CHAIN_LABELS[chain] ||
    "未知网络"
  );
}

const DEPOSIT_STATUS_LABELS: Record<string, string> = {
  COMPLETED: "已到账",
  PROCESSING: "确认中",
  FAILED: "失败",
  SUCCESS: "已到账",
};

const DEPOSIT_STATUS_CLASS: Record<string, string> = {
  COMPLETED: "text-green-600",
  SUCCESS: "text-green-600",
  FAILED: "text-red-600",
  PROCESSING: "text-amber-600",
};

export function formatDepositStatus(status: string): string {
  return DEPOSIT_STATUS_LABELS[status] || "处理中";
}

export function getDepositStatusClassName(status: string): string {
  return DEPOSIT_STATUS_CLASS[status] || "text-amber-600";
}

export function formatDepositOrderMeta(order: {
  currency?: string;
  network?: string;
  created_at: string;
}) {
  const parts: string[] = [];
  const networkLabel = formatChainLabel(order.currency || "", order.network);
  if (networkLabel) parts.push(networkLabel);
  parts.push(new Date(order.created_at).toLocaleString("zh-CN"));
  return parts.join(" · ");
}

export function maskTxHash(txHash: string) {
  if (txHash.length <= 16) return txHash;
  return `${txHash.slice(0, 8)}…${txHash.slice(-8)}`;
}
