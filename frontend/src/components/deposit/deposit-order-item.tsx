import { formatAmount } from "@/lib/utils/format";
import {
  formatDepositOrderMeta,
  formatDepositStatus,
  getDepositStatusClassName,
  maskTxHash,
} from "@/lib/utils/network";
import type { DepositOrder } from "@/types/trading";

interface DepositOrderItemProps {
  order: DepositOrder;
}

export function DepositOrderItem({ order }: DepositOrderItemProps) {
  const meta = formatDepositOrderMeta(order);
  const isCredited =
    order.status === "COMPLETED" || order.status === "SUCCESS";

  return (
    <div className="rounded-lg border p-3">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0 space-y-1">
          <p className="text-sm font-medium text-slate-900">
            {isCredited ? "充值到账" : "充值"}{" "}
            {formatAmount(order.amount, order.currency)} {order.currency}
          </p>
          <p className="text-xs text-slate-500">{meta}</p>
          {order.tx_hash && (
            <p
              className="truncate font-mono text-xs text-slate-400"
              title={order.tx_hash}
            >
              交易哈希 {maskTxHash(order.tx_hash)}
            </p>
          )}
          {!isCredited && order.status !== "FAILED" && (
            <p className="text-xs text-slate-500">链上确认中，到账后余额将自动更新</p>
          )}
        </div>
        <span
          className={`shrink-0 text-xs font-medium ${getDepositStatusClassName(order.status)}`}
        >
          {formatDepositStatus(order.status)}
        </span>
      </div>
    </div>
  );
}
