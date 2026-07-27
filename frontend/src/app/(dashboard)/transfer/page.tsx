"use client";

import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "next/navigation";
import { Loader2, ArrowDownUp } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { SimpleSelect } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { useTransfer, useTransferOrders } from "@/lib/hooks/use-trading";
import { useWalletBalances } from "@/lib/hooks/use-wallet";
import { useSupportedCurrencies } from "@/lib/hooks/use-supported-currencies";
import { formatAmount } from "@/lib/utils/format";
import { getAllSupportedCurrencies } from "@/lib/utils/currency";
import { toast } from "@/hooks/use-toast";

type AccountType = "FUNDING" | "TRADING";

const ACCOUNT_LABELS: Record<AccountType, string> = {
  FUNDING: "资金账户",
  TRADING: "交易账户",
};

function parseAccountType(value: string | null, fallback: AccountType): AccountType {
  return value === "FUNDING" || value === "TRADING" ? value : fallback;
}

export default function TransferPage() {
  const searchParams = useSearchParams();
  const { data: supportedCurrencies } = useSupportedCurrencies();
  const currencies = getAllSupportedCurrencies(supportedCurrencies);
  const [fromAccount, setFromAccount] = useState<AccountType>(() =>
    parseAccountType(searchParams.get("from"), "FUNDING"),
  );
  const [toAccount, setToAccount] = useState<AccountType>(() =>
    parseAccountType(searchParams.get("to"), "TRADING"),
  );
  const [currency, setCurrency] = useState("");
  const [amount, setAmount] = useState("");

  useEffect(() => {
    const fromParam = searchParams.get("from");
    const toParam = searchParams.get("to");
    if (!fromParam && !toParam) return;
    const from = parseAccountType(fromParam, "FUNDING");
    const to = parseAccountType(toParam, "TRADING");
    if (from !== to) {
      setFromAccount(from);
      setToAccount(to);
      setAmount("");
    }
  }, [searchParams]);

  useEffect(() => {
    const currencyParam = searchParams.get("currency")?.toUpperCase();
    const nextCurrency =
      currencyParam && currencies.includes(currencyParam)
        ? currencyParam
        : currencies[0];
    if (!nextCurrency) return;
    if (!currency || !currencies.includes(currency)) {
      setCurrency(nextCurrency);
    }
  }, [currencies, currency, searchParams]);

  const { data: walletData } = useWalletBalances();
  const transferMutation = useTransfer();
  const { data: ordersData, isLoading: ordersLoading } = useTransferOrders();
  const orders = ordersData?.orders || [];

  const wallets = useMemo(() => walletData?.wallets ?? [], [walletData?.wallets]);

  const sourceAvailable = useMemo(() => {
    const wallet = wallets.find(
      (item) => item.account_type === fromAccount && item.currency === currency,
    );
    return wallet?.available_balance ?? "0";
  }, [wallets, fromAccount, currency]);

  const fundingBalances = useMemo(
    () => wallets.filter((w) => w.account_type === "FUNDING"),
    [wallets],
  );
  const tradingBalances = useMemo(
    () => wallets.filter((w) => w.account_type === "TRADING"),
    [wallets],
  );

  function handleSwapAccounts() {
    setFromAccount(toAccount);
    setToAccount(fromAccount);
    setAmount("");
  }

  function handleMaxAmount() {
    setAmount(sourceAvailable);
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!amount) return;

    const parsed = parseFloat(amount);
    const available = parseFloat(sourceAvailable);
    if (Number.isNaN(parsed) || parsed <= 0) {
      toast({ variant: "destructive", title: "划转失败", description: "请输入有效金额。" });
      return;
    }
    if (parsed > available) {
      toast({
        variant: "destructive",
        title: "划转失败",
        description: `${ACCOUNT_LABELS[fromAccount]}可用余额不足。`,
      });
      return;
    }

    transferMutation.mutate(
      {
        from_account_type: fromAccount,
        to_account_type: toAccount,
        currency,
        amount,
      },
      {
        onSuccess: () => {
          toast({ title: "划转成功", description: "资金划转已完成。" });
          setAmount("");
        },
        onError: (err) => {
          toast({ variant: "destructive", title: "划转失败", description: err.message });
        },
      },
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-slate-900">划转</h1>
        <p className="mt-1 text-sm text-slate-600">
          支持资金账户与交易账户双向划转
        </p>
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-slate-600">资金账户余额</CardTitle>
          </CardHeader>
          <CardContent className="space-y-1">
            {fundingBalances.length === 0 ? (
              <p className="text-sm text-slate-400">暂无余额</p>
            ) : (
              fundingBalances.map((wallet) => (
                <div key={wallet.currency} className="flex justify-between text-sm">
                  <span className="text-slate-600">{wallet.currency}</span>
                  <span className="font-medium tabular-nums text-slate-900">
                    {formatAmount(wallet.available_balance, wallet.currency)}
                  </span>
                </div>
              ))
            )}
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-slate-600">交易账户余额</CardTitle>
          </CardHeader>
          <CardContent className="space-y-1">
            {tradingBalances.length === 0 ? (
              <p className="text-sm text-slate-400">暂无余额</p>
            ) : (
              tradingBalances.map((wallet) => (
                <div key={wallet.currency} className="flex justify-between text-sm">
                  <span className="text-slate-600">{wallet.currency}</span>
                  <span className="font-medium tabular-nums text-slate-900">
                    {formatAmount(wallet.available_balance, wallet.currency)}
                  </span>
                </div>
              ))
            )}
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="text-base">账户划转</CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="flex items-center gap-3">
                <div className="flex-1">
                  <label className="mb-1.5 block text-sm font-medium text-slate-700">从</label>
                  <div className="rounded-md border bg-slate-50 px-3 py-2 text-sm font-medium text-slate-900">
                    {ACCOUNT_LABELS[fromAccount]}
                  </div>
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="mt-6"
                  onClick={handleSwapAccounts}
                  title="切换划转方向"
                >
                  <ArrowDownUp className="h-4 w-4" />
                </Button>
                <div className="flex-1">
                  <label className="mb-1.5 block text-sm font-medium text-slate-700">到</label>
                  <div className="rounded-md border bg-slate-50 px-3 py-2 text-sm font-medium text-slate-900">
                    {ACCOUNT_LABELS[toAccount]}
                  </div>
                </div>
              </div>

              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">币种</label>
                <SimpleSelect
                  value={currency}
                  onValueChange={setCurrency}
                  options={currencies}
                />
              </div>

              <div>
                <div className="mb-1.5 flex items-center justify-between">
                  <label className="text-sm font-medium text-slate-700">金额</label>
                  <span className="text-xs text-slate-500">
                    可用 {formatAmount(sourceAvailable, currency)} {currency}
                  </span>
                </div>
                <div className="flex gap-2">
                  <Input
                    type="text"
                    placeholder="0.00"
                    value={amount}
                    onChange={(e) => setAmount(e.target.value)}
                    className="flex-1"
                  />
                  <Button type="button" variant="outline" onClick={handleMaxAmount}>
                    MAX
                  </Button>
                </div>
              </div>

              <Button
                type="submit"
                className="w-full bg-blue-700 hover:bg-blue-800 text-white"
                disabled={transferMutation.isPending || !amount}
              >
                {transferMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                确认划转
              </Button>
            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">划转记录</CardTitle>
          </CardHeader>
          <CardContent>
            {ordersLoading ? (
              <div className="space-y-3">
                <Skeleton className="h-12 w-full" />
                <Skeleton className="h-12 w-full" />
              </div>
            ) : orders.length === 0 ? (
              <p className="py-8 text-center text-sm text-slate-400">暂无划转记录</p>
            ) : (
              <div className="space-y-3">
                {orders.map((order) => (
                  <div key={order.id} className="flex items-center justify-between rounded-lg border p-3">
                    <div>
                      <p className="text-sm font-medium text-slate-900">
                        {formatAmount(order.amount, order.currency)} {order.currency}
                      </p>
                      <p className="text-xs text-slate-500">
                        {order.from_account_type === "FUNDING" ? "资金" : "交易"} → {order.to_account_type === "FUNDING" ? "资金" : "交易"}
                        {" · "}
                        {new Date(order.created_at).toLocaleString("zh-CN")}
                      </p>
                    </div>
                    <span className={`text-xs font-medium ${order.status === "COMPLETED" ? "text-green-600" : order.status === "FAILED" ? "text-red-600" : "text-amber-600"}`}>
                      {order.status === "COMPLETED" ? "已完成" : order.status === "FAILED" ? "失败" : "处理中"}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
