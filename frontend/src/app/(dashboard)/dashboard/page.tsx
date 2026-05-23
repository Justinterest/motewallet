"use client";

import Link from "next/link";
import {
  ArrowDownToLine,
  ArrowRightLeft,
  RefreshCw,
  Wallet,
} from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { useWalletBalances } from "@/lib/hooks/use-wallet";
import { formatAmount, getCurrencySymbol } from "@/lib/utils/format";
import type { WalletBalance } from "@/types/wallet";

function BalanceList({
  balances,
  isLoading,
}: {
  balances: WalletBalance[];
  isLoading: boolean;
}) {
  if (isLoading) {
    return (
      <div className="space-y-3">
        <div>
          <p className="text-sm text-slate-500">可用余额</p>
          <Skeleton className="mt-1 h-8 w-40" />
        </div>
        <div>
          <p className="text-sm text-slate-500">冻结金额</p>
          <Skeleton className="mt-1 h-6 w-32" />
        </div>
      </div>
    );
  }

  if (balances.length === 0) {
    return (
      <p className="py-4 text-center text-sm text-slate-400">暂无余额信息</p>
    );
  }

  return (
    <div className="space-y-3">
      {balances.map((wallet, index) => (
        <div key={`${wallet.currency}-${index}`}>
          {index > 0 && <Separator className="mb-3" />}
          <div className="flex items-baseline justify-between">
            <span className="text-sm font-medium text-slate-600">
              {getCurrencySymbol(wallet.currency)} {wallet.currency}
            </span>
          </div>
          <div className="mt-1 flex items-baseline gap-1">
            <span className="text-2xl font-bold text-slate-900">
              {formatAmount(wallet.available_balance, wallet.currency)}
            </span>
            <span className="text-xs text-slate-400">可用</span>
          </div>
          <div className="mt-0.5 text-sm text-slate-500">
            冻结：{formatAmount(wallet.frozen_balance, wallet.currency)}
          </div>
        </div>
      ))}
    </div>
  );
}

export default function DashboardPage() {
  const { data, isLoading } = useWalletBalances();

  const wallets = data?.wallets || [];
  const fundingBalances = wallets.filter((w) => w.account_type === "FUNDING");
  const tradingBalances = wallets.filter((w) => w.account_type === "TRADING");

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-slate-900">概览</h1>

      {/* Account cards */}
      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <Wallet className="h-5 w-5 text-blue-700" />
              资金账户
            </CardTitle>
          </CardHeader>
          <CardContent>
            <BalanceList balances={fundingBalances} isLoading={isLoading} />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <ArrowRightLeft className="h-5 w-5 text-blue-700" />
              交易账户
            </CardTitle>
          </CardHeader>
          <CardContent>
            <BalanceList balances={tradingBalances} isLoading={isLoading} />
          </CardContent>
        </Card>
      </div>

      {/* Quick actions */}
      <div className="flex flex-wrap gap-3">
        <Button variant="outline" asChild>
          <Link href="/deposit">
            <ArrowDownToLine className="mr-2 h-4 w-4" />
            充值
          </Link>
        </Button>
        <Button variant="outline" asChild>
          <Link href="/transfer">
            <ArrowRightLeft className="mr-2 h-4 w-4" />
            划转
          </Link>
        </Button>
        <Button variant="outline" asChild>
          <Link href="/exchange">
            <RefreshCw className="mr-2 h-4 w-4" />
            兑换
          </Link>
        </Button>
        <Button variant="outline" asChild>
          <Link href="/withdraw">
            <Wallet className="mr-2 h-4 w-4" />
            提现
          </Link>
        </Button>
      </div>

      {/* Recent transactions */}
      <div>
        <h2 className="mb-4 text-lg font-semibold text-slate-900">最近交易</h2>
        <Card>
          <CardContent className="flex items-center justify-center py-12">
            <p className="text-sm text-slate-400">暂无交易记录</p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
