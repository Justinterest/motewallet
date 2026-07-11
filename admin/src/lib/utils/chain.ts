/** Display labels for platform chain codes. */
export const CHAIN_LABELS: Record<string, string> = {
  ETH_ERC20: "ERC20（以太坊）",
  TRX_TRC20: "TRC20（波场）",
  TON: "TON",
  SOL_Solana: "Solana",
  BSC_BEP20: "BEP20（BNB Chain）",
  BTC: "Bitcoin",
  BTC_Bitcoin: "Bitcoin",
};

export function formatChainLabel(chain: string) {
  return CHAIN_LABELS[chain] ?? chain;
}
