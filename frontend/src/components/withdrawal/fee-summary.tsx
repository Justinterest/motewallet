import { Loader2 } from "lucide-react";
import { formatAmount } from "@/lib/utils/format";
import type { WithdrawalFeePreview } from "@/types/trading";

function hasPositiveFee(fee: string) {
  const value = parseFloat(fee);
  return !Number.isNaN(value) && value > 0;
}

interface WithdrawalFeeSummaryProps {
  preview?: WithdrawalFeePreview;
  isLoading?: boolean;
  isError?: boolean;
  showCryptoNetworkNote?: boolean;
}

export function WithdrawalFeeSummary({
  preview,
  isLoading,
  isError,
  showCryptoNetworkNote,
}: WithdrawalFeeSummaryProps) {
  if (isLoading) {
    return (
      <div className="flex items-center gap-2 rounded-lg border bg-slate-50 px-3 py-2 text-sm text-slate-500">
        <Loader2 className="h-4 w-4 animate-spin" />
        计算手续费中…
      </div>
    );
  }

  if (isError) {
    return (
      <div className="rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-sm text-amber-800">
        手续费计算失败，请检查金额后重试
      </div>
    );
  }

  if (!preview) {
    return null;
  }

  return (
    <div className="space-y-2 rounded-lg border bg-slate-50 px-3 py-3 text-sm">
      <div className="flex items-center justify-between">
        <span className="text-slate-600">提现金额</span>
        <span className="font-medium text-slate-900">
          {formatAmount(preview.amount, preview.currency)} {preview.currency}
        </span>
      </div>
      <div className="flex items-center justify-between">
        <span className="text-slate-600">平台手续费</span>
        <span className="font-medium text-slate-900">
          {hasPositiveFee(preview.platform_fee)
            ? `${formatAmount(preview.platform_fee, preview.currency)} ${preview.currency}`
            : "无"}
        </span>
      </div>
      <div className="flex items-center justify-between border-t border-slate-200 pt-2">
        <span className="text-slate-600">账户冻结总额</span>
        <span className="font-semibold text-slate-900">
          {formatAmount(preview.total_deduction, preview.currency)} {preview.currency}
        </span>
      </div>
      <div className="flex items-center justify-between">
        <span className="text-slate-600">预计到账</span>
        <span className="font-medium text-slate-900">
          {formatAmount(preview.net_amount, preview.currency)} {preview.currency}
        </span>
      </div>
      {showCryptoNetworkNote && (
        <p className="text-xs text-slate-500">链上网络手续费另行扣除，以实际到账为准。</p>
      )}
    </div>
  );
}
