"use client";

import Link from "next/link";
import {
  ArrowDownToLine,
  ArrowRight,
  ArrowRightLeft,
  RefreshCw,
  Wallet,
} from "lucide-react";

import { CurrencyIcon } from "@/components/currency-icon";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { useWalletBalances } from "@/lib/hooks/use-wallet";
import { useSupportedCurrencies } from "@/lib/hooks/use-supported-currencies";
import { formatAmount } from "@/lib/utils/format";
import { getAllSupportedCurrencies } from "@/lib/utils/currency";
import type { WalletBalance } from "@/types/wallet";
import type { LucideIcon } from "lucide-react";

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
            { label: "划转到交易账户", href: "/transfer?from=FUNDING&to=TRADING", icon: ArrowRightLeft },
          ]}
        />
        <AccountPanel
          title="交易账户"
          icon={ArrowRightLeft}
          balances={tradingBalances}
          isLoading={isLoading}
          actions={[
            { label: "去兑换", href: "/exchange", icon: RefreshCw },
            { label: "划转回资金账户", href: "/transfer?from=TRADING&to=FUNDING", icon: ArrowRightLeft },
          ]}
        />
      </div>

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
          <div className="hidden grid-cols-[1fr_auto_auto_auto] gap-4 border-b border-border bg-muted/40 px-6 py-3 text-xs font-medium uppercase tracking-wider text-muted-foreground sm:grid">
            <span>类型</span>
            <span className="text-right">金额</span>
            <span>状态</span>
            <span className="text-right">时间</span>
          </div>
          <div className="flex items-center justify-center px-6 py-14">
            <p className="text-sm text-muted-foreground">暂无交易记录</p>
          </div>
        </div>
      </section>
    </div>
  );
}
