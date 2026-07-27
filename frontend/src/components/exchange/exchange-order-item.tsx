import { formatAmount } from "@/lib/utils/format";
import {
  formatExchangeOrderHint,
  formatExchangeOrderMeta,
  formatExchangeStatus,
  getExchangeStatusClassName,
} from "@/lib/utils/exchange";
import type { ExchangeOrder } from "@/types/trading";

interface ExchangeOrderItemProps {
  order: ExchangeOrder;
}

export function ExchangeOrderItem({ order }: ExchangeOrderItemProps) {
  const hint = formatExchangeOrderHint(order);
  const meta = formatExchangeOrderMeta(order, formatAmount);
  const isCompleted = order.status === "COMPLETED" && order.to_amount;
  const creditedAmount = order.net_to_amount || order.to_amount;

  return (
    <div className="rounded-lg border p-3">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0 space-y-1">
          <p className="text-sm font-medium text-slate-900">
            卖出 {formatAmount(order.from_amount, order.from_currency)} {order.from_currency}
          </p>
          <p className="text-sm text-slate-700">
            {isCompleted ? "到账" : "买入"}{" "}
            {isCompleted
              ? `${formatAmount(creditedAmount, order.to_currency)} ${order.to_currency}`
              : order.to_currency}
          </p>
          {hint && <p className="text-xs text-slate-500">{hint}</p>}
          <p className="text-xs text-slate-500">{meta}</p>
          {order.status === "FAILED" && order.fail_reason && (
            <p className="text-xs text-red-600">失败原因：{order.fail_reason}</p>
          )}
        </div>
        <span
          className={`shrink-0 text-xs font-medium ${getExchangeStatusClassName(order.status)}`}
        >
          {formatExchangeStatus(order.status)}
        </span>
      </div>
    </div>
  );
}
