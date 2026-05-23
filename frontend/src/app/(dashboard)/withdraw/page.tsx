"use client";

import { useState } from "react";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { NativeSelect } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { useSubmitCryptoWithdrawal, useWithdrawalOrders } from "@/lib/hooks/use-trading";
import { formatAmount } from "@/lib/utils/format";
import { toast } from "@/hooks/use-toast";

const currencies = [
  { value: "USDT", label: "USDT" },
  { value: "USDC", label: "USDC" },
  { value: "BTC", label: "BTC" },
];

const chains: Record<string, { value: string; label: string }[]> = {
  USDT: [
    { value: "TRC20", label: "TRC20 (Tron)" },
    { value: "ERC20", label: "ERC20 (Ethereum)" },
  ],
  USDC: [
    { value: "ERC20", label: "ERC20 (Ethereum)" },
    { value: "TRC20", label: "TRC20 (Tron)" },
  ],
  BTC: [{ value: "BTC", label: "Bitcoin" }],
};

const statusLabels: Record<string, string> = {
  PENDING: "待审核",
  PROCESSING: "处理中",
  COMPLETED: "已完成",
  FAILED: "失败",
};

const reviewLabels: Record<string, string> = {
  PENDING_REVIEW: "待审核",
  APPROVED: "已批准",
  REJECTED: "已拒绝",
};

export default function WithdrawPage() {
  const [currency, setCurrency] = useState("USDT");
  const [chain, setChain] = useState("TRC20");
  const [amount, setAmount] = useState("");
  const [toAddress, setToAddress] = useState("");

  const submitMutation = useSubmitCryptoWithdrawal();
  const { data: ordersData, isLoading: ordersLoading } = useWithdrawalOrders();
  const orders = ordersData?.orders || [];

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!amount || !toAddress) return;

    submitMutation.mutate(
      { currency, chain, amount, to_address: toAddress },
      {
        onSuccess: () => {
          toast({ title: "提交成功", description: "提现申请已提交，等待审核。" });
          setAmount("");
          setToAddress("");
        },
        onError: (err) => {
          toast({ variant: "destructive", title: "提交失败", description: err.message });
        },
      }
    );
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-slate-900">提现</h1>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="text-base">加密货币提现</CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">币种</label>
                <NativeSelect value={currency} onChange={(e) => { setCurrency(e.target.value); setChain(chains[e.target.value]?.[0]?.value || ""); }}>
                  {currencies.map((c) => (
                    <option key={c.value} value={c.value}>{c.label}</option>
                  ))}
                </NativeSelect>
              </div>

              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">链/网络</label>
                <NativeSelect value={chain} onChange={(e) => setChain(e.target.value)}>
                  {(chains[currency] || []).map((c) => (
                    <option key={c.value} value={c.value}>{c.label}</option>
                  ))}
                </NativeSelect>
              </div>

              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">提现金额</label>
                <Input
                  type="text"
                  placeholder="0.00"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                />
              </div>

              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">收款地址</label>
                <Input
                  type="text"
                  placeholder="请输入收款地址"
                  value={toAddress}
                  onChange={(e) => setToAddress(e.target.value)}
                />
              </div>

              <Button
                type="submit"
                className="w-full bg-blue-700 hover:bg-blue-800 text-white"
                disabled={submitMutation.isPending || !amount || !toAddress}
              >
                {submitMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                提交提现
              </Button>
            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">提现记录</CardTitle>
          </CardHeader>
          <CardContent>
            {ordersLoading ? (
              <div className="space-y-3">
                <Skeleton className="h-14 w-full" />
                <Skeleton className="h-14 w-full" />
              </div>
            ) : orders.length === 0 ? (
              <p className="py-8 text-center text-sm text-slate-400">暂无提现记录</p>
            ) : (
              <div className="space-y-3">
                {orders.map((order) => (
                  <div key={order.id} className="rounded-lg border p-3">
                    <div className="flex items-center justify-between">
                      <p className="text-sm font-medium text-slate-900">
                        {formatAmount(order.amount, order.currency)} {order.currency}
                      </p>
                      <span className={`text-xs font-medium ${order.status === "COMPLETED" ? "text-green-600" : order.status === "FAILED" ? "text-red-600" : "text-amber-600"}`}>
                        {statusLabels[order.status] || order.status}
                      </span>
                    </div>
                    <div className="mt-1 flex items-center justify-between text-xs text-slate-500">
                      <span>{order.chain} · {reviewLabels[order.review_status] || order.review_status}</span>
                      <span>{new Date(order.created_at).toLocaleString("zh-CN")}</span>
                    </div>
                    {order.platform_fee && order.platform_fee !== "0" && (
                      <p className="mt-1 text-xs text-slate-400">手续费：{order.platform_fee} {order.currency}</p>
                    )}
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
