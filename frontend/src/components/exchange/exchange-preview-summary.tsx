import { Loader2 } from "lucide-react";
import { formatAmount } from "@/lib/utils/format";
import type { ExchangePreview } from "@/types/trading";

function hasPositiveFee(fee: string) {
  const value = parseFloat(fee);
  return !Number.isNaN(value) && value > 0;
}

interface ExchangePreviewSummaryProps {
  preview?: ExchangePreview;
  fromCurrency: string;
  toCurrency: string;
  isLoading?: boolean;
  isFetching?: boolean;
  isError?: boolean;
}

export function ExchangePreviewSummary({
  preview,
  fromCurrency,
  toCurrency,
  isLoading,
  isFetching,
  isError,
}: ExchangePreviewSummaryProps) {
  if (isLoading || isFetching) {
    return (
      <div className="flex items-center gap-2 rounded-lg border bg-slate-50 px-3 py-2 text-sm text-slate-500">
        <Loader2 className="h-4 w-4 animate-spin" />
        计算兑换信息中…
      </div>
    );
  }

  if (isError) {
    return (
      <div className="rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-sm text-amber-800">
        无法计算兑换信息，请检查金额和币对后重试
      </div>
    );
  }

  if (!preview) {
    return null;
  }

  return (
    <div className="space-y-2 rounded-lg border bg-slate-50 px-3 py-3 text-sm">
      <div className="flex items-center justify-between">
        <span className="text-slate-600">兑换比例</span>
        <span className="font-medium text-slate-900">
          1 {fromCurrency} = {formatAmount(preview.exchange_rate, toCurrency)}{" "}
          {toCurrency}
        </span>
      </div>
      <div className="flex items-center justify-between">
        <span className="text-slate-600">卖出金额</span>
        <span className="font-medium text-slate-900">
          {formatAmount(preview.from_amount, fromCurrency)} {fromCurrency}
        </span>
      </div>
      <div className="flex items-center justify-between">
        <span className="text-slate-600">平台手续费</span>
        <span className="font-medium text-slate-900">
          {hasPositiveFee(preview.platform_fee)
            ? `${formatAmount(preview.platform_fee, preview.fee_currency)} ${preview.fee_currency}`
            : "无"}
        </span>
      </div>
      <div className="flex items-center justify-between">
        <span className="text-slate-600">账户扣款总额</span>
        <span className="font-medium text-slate-900">
          {formatAmount(preview.total_deduction, fromCurrency)} {fromCurrency}
        </span>
      </div>
      <div className="flex items-center justify-between border-t border-slate-200 pt-2">
        <span className="text-slate-600">预估到账（交易账户）</span>
        <span className="font-semibold text-slate-900">
          {formatAmount(preview.net_to_amount, toCurrency)} {toCurrency}
        </span>
      </div>
    </div>
  );
}
