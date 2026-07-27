"use client";

import { useMemo } from "react";
import Link from "next/link";
import {
  ArrowDownToLine,
  ArrowRight,
  ArrowRightLeft,
  RefreshCw,
  Wallet,
} from "lucide-react";

import { CurrencyIcon } from "@/components/currency-icon";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import {
  useDepositOrders,
  useExchangeOrders,
  useTransferOrders,
  useWithdrawalOrders,
} from "@/lib/hooks/use-trading";
import { useWalletBalances } from "@/lib/hooks/use-wallet";
import { useSupportedCurrencies } from "@/lib/hooks/use-supported-currencies";
import { formatAmount } from "@/lib/utils/format";
import { getAllSupportedCurrencies } from "@/lib/utils/currency";
import type { WalletBalance } from "@/types/wallet";
import type { LucideIcon } from "lucide-react";

const RECENT_TX_LIMIT = 8;

const statusLabels: Record<string, string> = {
  COMPLETED: "已完成",
  FAILED: "失败",
  PROCESSING: "处理中",
  PENDING: "待处理",
};

const statusColors: Record<string, string> = {
  COMPLETED: "bg-green-50 text-green-700",
  FAILED: "bg-red-50 text-red-700",
  PROCESSING: "bg-amber-50 text-amber-700",
  PENDING: "bg-slate-100 text-slate-600",
};

type RecentTx = {
  id: string;
  typeLabel: string;
  amountText: string;
  status: string;
  createdAt: string;
};

function BalanceRows({
  balances,
  isLoading,
}: {
  balances: WalletBalance[];
  isLoading: boolean;
}) {
  if (isLoading) {
    return (
      <div className="space-y-4">
        {Array.from({ length: 3 }).map((_, i) => (
          <div key={i} className="flex items-center justify-between">
            <Skeleton className="h-4 w-24" />
            <Skeleton className="h-6 w-32" />
          </div>
        ))}
      </div>
    );
  }

  if (balances.length === 0) {
    return (
      <p className="py-6 text-center text-sm text-muted-foreground">
        暂无余额信息
      </p>
    );
  }

  return (
    <div className="divide-y divide-border">
      {balances.map((wallet, index) => (
        <div
          key={`${wallet.currency}-${index}`}
          className="flex items-center justify-between py-3 first:pt-0 last:pb-0"
        >
          <div className="flex items-center gap-2.5">
            <CurrencyIcon currency={wallet.currency} />
            <div>
              <p className="text-sm font-medium text-foreground">
                {wallet.currency}
              </p>
              {parseFloat(wallet.frozen_balance) > 0 && (
                <p className="text-xs text-muted-foreground">
                  冻结 {formatAmount(wallet.frozen_balance, wallet.currency)}
                </p>
              )}
            </div>
          </div>
          <span className="text-right text-base font-semibold tabular-nums text-foreground">
            {formatAmount(wallet.available_balance, wallet.currency)}
          </span>
        </div>
      ))}
    </div>
  );
}

function AccountPanel({
  title,
  icon: Icon,
  balances,
  isLoading,
  actions,
}: {
  title: string;
  icon: LucideIcon;
  balances: WalletBalance[];
  isLoading: boolean;
  actions: { label: string; href: string; icon: LucideIcon }[];
}) {
  return (
    <div className="panel flex flex-col">
      <div className="flex items-center gap-2.5 border-b border-border px-6 py-4">
        <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-accent">
          <Icon className="h-4 w-4 text-primary" />
        </div>
        <h2 className="text-base font-semibold text-foreground">{title}</h2>
      </div>
      <div className="flex-1 px-6 py-4">
        <BalanceRows balances={balances} isLoading={isLoading} />
      </div>
      <div className="flex flex-wrap gap-2 border-t border-border px-6 py-4">
        {actions.map((action) => (
          <Button key={action.href} variant="outline" size="sm" asChild>
            <Link href={action.href}>
              <action.icon className="mr-1.5 h-3.5 w-3.5" />
              {action.label}
            </Link>
          </Button>
        ))}
      </div>
    </div>
  );
}

function RecentTransactions() {
  const { data: depositData, isLoading: depositLoading } = useDepositOrders();
  const { data: withdrawalData, isLoading: withdrawalLoading } = useWithdrawalOrders();
  const { data: exchangeData, isLoading: exchangeLoading } = useExchangeOrders();
  const { data: transferData, isLoading: transferLoading } = useTransferOrders();

  const isLoading =
    depositLoading || withdrawalLoading || exchangeLoading || transferLoading;

  const recentTxs = useMemo(() => {
    const items: RecentTx[] = [];

    for (const order of depositData?.orders || []) {
      items.push({
        id: `deposit-${order.id}`,
        typeLabel: "充值",
        amountText: `+${formatAmount(order.amount, order.currency)} ${order.currency}`,
        status: order.status,
        createdAt: order.created_at,
      });
    }

    for (const order of withdrawalData?.orders || []) {
      items.push({
        id: `withdrawal-${order.id}`,
        typeLabel: "提现",
        amountText: `-${formatAmount(order.amount, order.currency)} ${order.currency}`,
        status: order.status,
        createdAt: order.created_at,
      });
    }

    for (const order of exchangeData?.orders || []) {
      const creditedAmount = order.net_to_amount || order.to_amount;
      const toAmount = creditedAmount
        ? formatAmount(creditedAmount, order.to_currency)
        : "—";
      items.push({
        id: `exchange-${order.id}`,
        typeLabel: "兑换",
        amountText: `${formatAmount(order.from_amount, order.from_currency)} ${order.from_currency} → ${toAmount} ${order.to_currency}`,
        status: order.status,
        createdAt: order.created_at,
      });
    }

    for (const order of transferData?.orders || []) {
      const fromLabel = order.from_account_type === "FUNDING" ? "资金" : "交易";
      const toLabel = order.to_account_type === "FUNDING" ? "资金" : "交易";
      items.push({
        id: `transfer-${order.id}`,
        typeLabel: `划转（${fromLabel}→${toLabel}）`,
        amountText: `${formatAmount(order.amount, order.currency)} ${order.currency}`,
        status: order.status,
        createdAt: order.created_at,
      });
    }

    return items
      .sort(
        (a, b) =>
          new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
      )
      .slice(0, RECENT_TX_LIMIT);
  }, [depositData, withdrawalData, exchangeData, transferData]);

  return (
    <section>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold text-foreground">最近交易</h2>
        <Link
          href="/transactions"
          className="inline-flex items-center text-sm font-medium text-primary hover:underline"
        >
          查看全部
          <ArrowRight className="ml-1 h-3.5 w-3.5" />
        </Link>
      </div>
      <div className="panel overflow-hidden">
        {isLoading ? (
          <div className="space-y-3 px-6 py-4">
            {Array.from({ length: 4 }).map((_, i) => (
              <Skeleton key={i} className="h-10 w-full" />
            ))}
          </div>
        ) : recentTxs.length === 0 ? (
          <div className="flex items-center justify-center px-6 py-14">
            <p className="text-sm text-muted-foreground">暂无交易记录</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[36rem] border-collapse text-sm">
              <thead>
                <tr className="border-b border-border bg-muted/40 text-xs font-medium uppercase tracking-wider text-muted-foreground">
                  <th className="px-6 py-3 text-left font-medium">类型</th>
                  <th className="px-6 py-3 text-right font-medium">金额</th>
                  <th className="px-6 py-3 text-left font-medium">状态</th>
                  <th className="px-6 py-3 text-right font-medium">时间</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {recentTxs.map((tx) => (
                  <tr key={tx.id}>
                    <td className="px-6 py-3.5 text-left font-medium text-foreground">
                      {tx.typeLabel}
                    </td>
                    <td className="px-6 py-3.5 text-right font-medium tabular-nums text-foreground">
                      {tx.amountText}
                    </td>
                    <td className="px-6 py-3.5 text-left">
                      <Badge
                        variant="secondary"
                        className={statusColors[tx.status] || statusColors.PENDING}
                      >
                        {statusLabels[tx.status] || tx.status}
                      </Badge>
                    </td>
                    <td className="px-6 py-3.5 text-right text-muted-foreground">
                      {new Date(tx.createdAt).toLocaleString("zh-CN")}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </section>
  );
}

export default function DashboardPage() {
  const { data, isLoading: balancesLoading } = useWalletBalances();
  const { data: supportedCurrencies, isLoading: currenciesLoading } =
    useSupportedCurrencies();
  const isLoading = balancesLoading || currenciesLoading;

  const supported = new Set(getAllSupportedCurrencies(supportedCurrencies));
  const wallets = (data?.wallets || []).filter((wallet) =>
    supported.has(wallet.currency),
  );
  const fundingBalances = wallets.filter((w) => w.account_type === "FUNDING");
  const tradingBalances = wallets.filter((w) => w.account_type === "TRADING");

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-[32px] font-bold tracking-tight text-foreground">
          概览
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          查看账户余额与最近交易
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        <AccountPanel
          title="资金账户"
          icon={Wallet}
          balances={fundingBalances}
          isLoading={isLoading}
          actions={[
            { label: "充值", href: "/deposit", icon: ArrowDownToLine },
            { label: "提现", href: "/withdraw", icon: Wallet },
            { label: "兑换", href: "/exchange", icon: RefreshCw },
          ]}
        />
        <AccountPanel
          title="交易账户"
          icon={ArrowRightLeft}
          balances={tradingBalances}
          isLoading={isLoading}
          actions={[
            { label: "划转", href: "/transfer?from=TRADING&to=FUNDING", icon: ArrowRightLeft },
          ]}
        />
      </div>

      <RecentTransactions />
    </div>
  );
}
