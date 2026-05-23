"use client";

import { useState } from "react";
import { Copy, Check } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { NativeSelect } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { useDepositAddresses, useDepositOrders } from "@/lib/hooks/use-trading";
import { formatAmount } from "@/lib/utils/format";

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

export default function DepositPage() {
  const [currency, setCurrency] = useState("USDT");
  const [chain, setChain] = useState("TRC20");
  const [copied, setCopied] = useState(false);

  const { data: addressData, isLoading: addressLoading } = useDepositAddresses(currency, chain);
  const { data: ordersData, isLoading: ordersLoading } = useDepositOrders();

  const orders = ordersData?.orders || [];

  function handleCopy(text: string) {
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-slate-900">充值</h1>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="text-base">获取充值地址</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
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

            {addressLoading ? (
              <Skeleton className="h-20 w-full" />
            ) : addressData?.address ? (
              <div className="rounded-lg border bg-slate-50 p-4">
                <p className="mb-1 text-xs text-slate-500">充值地址</p>
                <p className="break-all font-mono text-sm text-slate-900">{addressData.address}</p>
                <Button
                  variant="outline"
                  size="sm"
                  className="mt-2"
                  onClick={() => handleCopy(addressData.address)}
                >
                  {copied ? <Check className="mr-1 h-3 w-3" /> : <Copy className="mr-1 h-3 w-3" />}
                  {copied ? "已复制" : "复制地址"}
                </Button>
              </div>
            ) : (
              <p className="text-sm text-slate-400">暂无可用地址</p>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">充值记录</CardTitle>
          </CardHeader>
          <CardContent>
            {ordersLoading ? (
              <div className="space-y-3">
                <Skeleton className="h-12 w-full" />
                <Skeleton className="h-12 w-full" />
              </div>
            ) : orders.length === 0 ? (
              <p className="py-8 text-center text-sm text-slate-400">暂无充值记录</p>
            ) : (
              <div className="space-y-3">
                {orders.map((order) => (
                  <div key={order.id} className="flex items-center justify-between rounded-lg border p-3">
                    <div>
                      <p className="text-sm font-medium text-slate-900">
                        {formatAmount(order.amount, order.currency)} {order.currency}
                      </p>
                      <p className="text-xs text-slate-500">{order.chain} · {new Date(order.created_at).toLocaleString("zh-CN")}</p>
                    </div>
                    <span className="text-xs font-medium text-green-600">{order.status}</span>
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
