"use client";

import { useEffect, useState } from "react";
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
import { useSupportedCurrencies } from "@/lib/hooks/use-supported-currencies";
import { formatAmount } from "@/lib/utils/format";
import { getAllSupportedCurrencies } from "@/lib/utils/currency";
import { toast } from "@/hooks/use-toast";

export default function TransferPage() {
  const { data: supportedCurrencies } = useSupportedCurrencies();
  const currencies = getAllSupportedCurrencies(supportedCurrencies);
  const [fromAccount, setFromAccount] = useState("FUNDING");
  const [toAccount, setToAccount] = useState("TRADING");
  const [currency, setCurrency] = useState("");
  const [amount, setAmount] = useState("");

  useEffect(() => {
    const nextCurrency = currencies[0];
    if (!nextCurrency) return;
    if (!currency || !currencies.includes(currency)) {
      setCurrency(nextCurrency);
    }
  }, [currencies, currency]);

  const transferMutation = useTransfer();
  const { data: ordersData, isLoading: ordersLoading } = useTransferOrders();
  const orders = ordersData?.orders || [];

  function handleSwapAccounts() {
    setFromAccount(toAccount);
    setToAccount(fromAccount);
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!amount) return;

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
      }
    );
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-slate-900">划转</h1>

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
                    {fromAccount === "FUNDING" ? "资金账户" : "交易账户"}
                  </div>
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="mt-6"
                  onClick={handleSwapAccounts}
                >
                  <ArrowDownUp className="h-4 w-4" />
                </Button>
                <div className="flex-1">
                  <label className="mb-1.5 block text-sm font-medium text-slate-700">到</label>
                  <div className="rounded-md border bg-slate-50 px-3 py-2 text-sm font-medium text-slate-900">
                    {toAccount === "FUNDING" ? "资金账户" : "交易账户"}
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
                <label className="mb-1.5 block text-sm font-medium text-slate-700">金额</label>
                <Input
                  type="text"
                  placeholder="0.00"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                />
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
