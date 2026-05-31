"use client";

import { useState } from "react";
import { Loader2, ArrowRightLeft } from "lucide-react";
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
import { useExchangeQuote, useCreateExchangeOrder, useExchangeOrders } from "@/lib/hooks/use-trading";
import { formatAmount } from "@/lib/utils/format";
import { toast } from "@/hooks/use-toast";
import type { ExchangeQuote } from "@/types/trading";

const allCurrencies = ["USDT", "USDC", "BTC", "USD", "HKD", "EUR"];

export default function ExchangePage() {
  const [fromCurrency, setFromCurrency] = useState("USDT");
  const [toCurrency, setToCurrency] = useState("USD");
  const [fromAmount, setFromAmount] = useState("");
  const [quote, setQuote] = useState<ExchangeQuote | null>(null);

  const quoteMutation = useExchangeQuote();
  const orderMutation = useCreateExchangeOrder();
  const { data: ordersData, isLoading: ordersLoading } = useExchangeOrders();
  const orders = ordersData?.orders || [];

  function handleGetQuote() {
    if (!fromAmount || fromCurrency === toCurrency) return;
    quoteMutation.mutate(
      { from_currency: fromCurrency, to_currency: toCurrency, from_amount: fromAmount },
      {
        onSuccess: (data) => setQuote(data),
        onError: (err) => {
          toast({ variant: "destructive", title: "询价失败", description: err.message });
          setQuote(null);
        },
      }
    );
  }

  function handleConfirmOrder() {
    if (!quote) return;
    orderMutation.mutate(
      {
        quote_id: quote.quote_id,
        from_currency: fromCurrency,
        to_currency: toCurrency,
        from_amount: fromAmount,
      },
      {
        onSuccess: () => {
          toast({ title: "下单成功", description: "兑换订单已创建。" });
          setFromAmount("");
          setQuote(null);
        },
        onError: (err) => {
          toast({ variant: "destructive", title: "下单失败", description: err.message });
        },
      }
    );
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-slate-900">兑换</h1>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="text-base">币种兑换</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-[1fr_auto_1fr] items-end gap-2">
              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">卖出</label>
                <SimpleSelect
                  value={fromCurrency}
                  onValueChange={(value) => {
                    setFromCurrency(value);
                    setQuote(null);
                  }}
                  options={allCurrencies.filter((c) => c !== toCurrency)}
                />
              </div>
              <ArrowRightLeft className="mb-2 h-5 w-5 text-slate-400" />
              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">买入</label>
                <SimpleSelect
                  value={toCurrency}
                  onValueChange={(value) => {
                    setToCurrency(value);
                    setQuote(null);
                  }}
                  options={allCurrencies.filter((c) => c !== fromCurrency)}
                />
              </div>
            </div>

            <div>
              <label className="mb-1.5 block text-sm font-medium text-slate-700">卖出金额</label>
              <Input
                type="text"
                placeholder="0.00"
                value={fromAmount}
                onChange={(e) => { setFromAmount(e.target.value); setQuote(null); }}
              />
            </div>

            <Button
              variant="outline"
              className="w-full"
              onClick={handleGetQuote}
              disabled={quoteMutation.isPending || !fromAmount}
            >
              {quoteMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              获取报价
            </Button>

            {quote && (
              <div className="rounded-lg border bg-blue-50 p-4 space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="text-slate-600">汇率</span>
                  <span className="font-medium text-slate-900">1 {fromCurrency} = {quote.exchange_rate} {toCurrency}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-slate-600">预计获得</span>
                  <span className="font-medium text-slate-900">{quote.to_amount} {toCurrency}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-slate-600">手续费</span>
                  <span className="text-slate-900">{quote.platform_fee} {quote.fee_currency}</span>
                </div>
                <Button
                  className="mt-2 w-full bg-blue-700 hover:bg-blue-800 text-white"
                  onClick={handleConfirmOrder}
                  disabled={orderMutation.isPending}
                >
                  {orderMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  确认兑换
                </Button>
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">兑换记录</CardTitle>
          </CardHeader>
          <CardContent>
            {ordersLoading ? (
              <div className="space-y-3">
                <Skeleton className="h-14 w-full" />
                <Skeleton className="h-14 w-full" />
              </div>
            ) : orders.length === 0 ? (
              <p className="py-8 text-center text-sm text-slate-400">暂无兑换记录</p>
            ) : (
              <div className="space-y-3">
                {orders.map((order) => (
                  <div key={order.id} className="rounded-lg border p-3">
                    <div className="flex items-center justify-between">
                      <p className="text-sm font-medium text-slate-900">
                        {formatAmount(order.from_amount, order.from_currency)} {order.from_currency} → {order.to_amount ? formatAmount(order.to_amount, order.to_currency) : "—"} {order.to_currency}
                      </p>
                      <span className={`text-xs font-medium ${order.status === "COMPLETED" ? "text-green-600" : order.status === "FAILED" ? "text-red-600" : "text-amber-600"}`}>
                        {order.status === "COMPLETED" ? "已完成" : order.status === "FAILED" ? "失败" : "处理中"}
                      </span>
                    </div>
                    <div className="mt-1 flex items-center justify-between text-xs text-slate-500">
                      <span>{order.exchange_rate ? `汇率 ${order.exchange_rate}` : order.exchange_type}</span>
                      <span>{new Date(order.created_at).toLocaleString("zh-CN")}</span>
                    </div>
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
