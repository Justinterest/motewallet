"use client";

import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import {
  Card,
  CardContent,
  CardHeader,
} from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { useDepositOrders, useWithdrawalOrders, useExchangeOrders, useTransferOrders } from "@/lib/hooks/use-trading";
import { useWalletLedger } from "@/lib/hooks/use-wallet";
import { formatAmount } from "@/lib/utils/format";
import type { WalletLedgerEntry } from "@/types/wallet";

type TabKey = "ledger" | "deposit" | "withdrawal" | "exchange" | "transfer";

const tabs: { key: TabKey; label: string }[] = [
  { key: "ledger", label: "资金变化" },
  { key: "deposit", label: "充值" },
  { key: "withdrawal", label: "提现" },
  { key: "exchange", label: "兑换" },
  { key: "transfer", label: "划转" },
];

function isTabKey(value: string | null): value is TabKey {
  return (
    value === "ledger" ||
    value === "deposit" ||
    value === "withdrawal" ||
    value === "exchange" ||
    value === "transfer"
  );
}

const statusColors: Record<string, string> = {
  COMPLETED: "bg-green-50 text-green-700",
  FAILED: "bg-red-50 text-red-700",
  PROCESSING: "bg-amber-50 text-amber-700",
  PENDING: "bg-slate-100 text-slate-600",
};

const entryTypeLabels: Record<string, string> = {
  CREDIT: "入账",
  FREEZE: "冻结",
  UNFREEZE: "解冻",
  DEDUCT_FROZEN: "扣款",
};

const bizTypeLabels: Record<string, string> = {
  DEPOSIT: "充值",
  WITHDRAWAL: "提现",
  EXCHANGE: "兑换",
  TRANSFER: "划转",
};

const accountTypeLabels: Record<string, string> = {
  FUNDING: "资金账户",
  TRADING: "交易账户",
};

function formatSignedAmount(entry: WalletLedgerEntry): { text: string; className: string } {
  const amount = formatAmount(entry.amount, entry.currency);
  switch (entry.entry_type) {
    case "CREDIT":
    case "UNFREEZE":
      return { text: `+${amount}`, className: "text-green-700" };
    case "DEDUCT_FROZEN":
    case "FREEZE":
      return { text: `-${amount}`, className: "text-red-700" };
    default:
      return { text: amount, className: "text-slate-900" };
  }
}

export default function TransactionsPage() {
  const searchParams = useSearchParams();
  const [activeTab, setActiveTab] = useState<TabKey>("ledger");

  useEffect(() => {
    const tab = searchParams.get("tab");
    if (isTabKey(tab)) {
      setActiveTab(tab);
    }
  }, [searchParams]);

  const { data: ledgerData, isLoading: ledgerLoading } = useWalletLedger(
    { page: 1, page_size: 50 },
    activeTab === "ledger",
  );
  const { data: depositData, isLoading: depositLoading } = useDepositOrders();
  const { data: withdrawalData, isLoading: withdrawalLoading } = useWithdrawalOrders();
  const { data: exchangeData, isLoading: exchangeLoading } = useExchangeOrders();
  const { data: transferData, isLoading: transferLoading } = useTransferOrders();

  function renderContent() {
    switch (activeTab) {
      case "ledger": {
        const entries = ledgerData?.entries || [];
        if (ledgerLoading) return <LoadingSkeleton />;
        if (entries.length === 0) return <EmptyState />;
        return entries.map((entry) => {
          const signed = formatSignedAmount(entry);
          const biz = entry.biz_type ? bizTypeLabels[entry.biz_type] || entry.biz_type : null;
          return (
            <div
              key={entry.id}
              className="flex items-start justify-between gap-4 rounded-lg border p-3"
            >
              <div className="min-w-0 space-y-1">
                <div className="flex flex-wrap items-center gap-2">
                  <p className="text-sm font-medium text-slate-900">
                    {entryTypeLabels[entry.entry_type] || entry.entry_type}
                  </p>
                  <Badge variant="secondary" className="bg-slate-100 text-slate-600">
                    {accountTypeLabels[entry.account_type] || entry.account_type}
                  </Badge>
                  {biz && (
                    <Badge variant="secondary" className="bg-blue-50 text-blue-700">
                      {biz}
                    </Badge>
                  )}
                </div>
                <p className="text-xs text-slate-500">
                  {entry.currency} · 余额 {formatAmount(entry.balance_before, entry.currency)} →{" "}
                  {formatAmount(entry.balance_after, entry.currency)}
                  {" · "}冻结 {formatAmount(entry.frozen_before, entry.currency)} →{" "}
                  {formatAmount(entry.frozen_after, entry.currency)}
                </p>
                <p className="text-xs text-slate-400">
                  {new Date(entry.created_at).toLocaleString("zh-CN")}
                  {entry.platform_order_id ? ` · ${entry.platform_order_id}` : ""}
                  {entry.remark ? ` · ${entry.remark}` : ""}
                </p>
              </div>
              <p className={`shrink-0 text-sm font-semibold ${signed.className}`}>
                {signed.text} {entry.currency}
              </p>
            </div>
          );
        });
      }
      case "deposit": {
        const orders = depositData?.orders || [];
        if (depositLoading) return <LoadingSkeleton />;
        if (orders.length === 0) return <EmptyState />;
        return orders.map((o) => (
          <TxRow
            key={o.id}
            title={`${formatAmount(o.amount, o.currency)} ${o.currency}`}
            sub={`${o.network} · ${new Date(o.created_at).toLocaleString("zh-CN")}`}
            status={o.status}
          />
        ));
      }
      case "withdrawal": {
        const orders = withdrawalData?.orders || [];
        if (withdrawalLoading) return <LoadingSkeleton />;
        if (orders.length === 0) return <EmptyState />;
        return orders.map((o) => (
          <TxRow
            key={o.id}
            title={`${formatAmount(o.amount, o.currency)} ${o.currency}`}
            sub={`${o.type === "CRYPTO" ? o.chain || "" : "法币"} · ${new Date(o.created_at).toLocaleString("zh-CN")}`}
            status={o.status}
          />
        ));
      }
      case "exchange": {
        const orders = exchangeData?.orders || [];
        if (exchangeLoading) return <LoadingSkeleton />;
        if (orders.length === 0) return <EmptyState />;
        return orders.map((o) => (
          <TxRow
            key={o.id}
            title={`${formatAmount(o.from_amount, o.from_currency)} ${o.from_currency} → ${o.to_amount ? formatAmount(o.to_amount, o.to_currency) : "—"} ${o.to_currency}`}
            sub={`${o.exchange_type} · ${new Date(o.created_at).toLocaleString("zh-CN")}${o.status === "FAILED" && o.fail_reason ? ` · 失败原因：${o.fail_reason}` : ""}`}
            status={o.status}
          />
        ));
      }
      case "transfer": {
        const orders = transferData?.orders || [];
        if (transferLoading) return <LoadingSkeleton />;
        if (orders.length === 0) return <EmptyState />;
        return orders.map((o) => (
          <TxRow
            key={o.id}
            title={`${formatAmount(o.amount, o.currency)} ${o.currency}`}
            sub={`${o.from_account_type === "FUNDING" ? "资金" : "交易"} → ${o.to_account_type === "FUNDING" ? "资金" : "交易"} · ${new Date(o.created_at).toLocaleString("zh-CN")}`}
            status={o.status}
          />
        ));
      }
    }
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-slate-900">交易记录</h1>

      <Card>
        <CardHeader className="pb-2">
          <div className="flex flex-wrap gap-1">
            {tabs.map((tab) => (
              <button
                key={tab.key}
                onClick={() => setActiveTab(tab.key)}
                className={`rounded-md px-3 py-1.5 text-sm font-medium transition-colors ${
                  activeTab === tab.key
                    ? "bg-blue-50 text-blue-700"
                    : "text-slate-500 hover:text-slate-700"
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">{renderContent()}</div>
        </CardContent>
      </Card>
    </div>
  );
}

function TxRow({ title, sub, status }: { title: string; sub: string; status: string }) {
  return (
    <div className="flex items-center justify-between rounded-lg border p-3">
      <div>
        <p className="text-sm font-medium text-slate-900">{title}</p>
        <p className="text-xs text-slate-500">{sub}</p>
      </div>
      <Badge variant="secondary" className={statusColors[status] || statusColors.PENDING}>
        {status === "COMPLETED" ? "已完成" : status === "FAILED" ? "失败" : status === "PROCESSING" ? "处理中" : "待处理"}
      </Badge>
    </div>
  );
}

function LoadingSkeleton() {
  return (
    <>
      <Skeleton className="h-14 w-full" />
      <Skeleton className="h-14 w-full" />
      <Skeleton className="h-14 w-full" />
    </>
  );
}

function EmptyState() {
  return <p className="py-8 text-center text-sm text-slate-400">暂无记录</p>;
}
