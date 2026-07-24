"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { ArrowRight, ArrowRightLeft, Loader2 } from "lucide-react";
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
import { ExchangePreviewSummary } from "@/components/exchange/exchange-preview-summary";
import { ExchangeOrderItem } from "@/components/exchange/exchange-order-item";
import {
  useCreateExchangeOrder,
  useExchangeOrders,
  useExchangePreview,
} from "@/lib/hooks/use-trading";
import { useWalletBalances } from "@/lib/hooks/use-wallet";
import { useSupportedCurrencies } from "@/lib/hooks/use-supported-currencies";
import { formatAmount } from "@/lib/utils/format";
import { getAllSupportedCurrencies } from "@/lib/utils/currency";
import {
  getExchangeFromOptions,
  getExchangeToOptions,
  isSupportedExchangePair,
} from "@/lib/utils/exchange";
import { toast } from "@/hooks/use-toast";

const DEFAULT_FROM_CURRENCY = "USDT";
const DEFAULT_TO_CURRENCY = "USD";

function isPositiveAmount(value: string) {
  const trimmed = value.trim();
  if (!trimmed) return false;
  const parsed = parseFloat(trimmed);
  return !Number.isNaN(parsed) && parsed > 0;
}

export default function ExchangePage() {
  const { data: supportedCurrencies } = useSupportedCurrencies();
  const allCurrencies = getAllSupportedCurrencies(supportedCurrencies);
  const exchangeCurrencies = useMemo(
    () => getExchangeFromOptions(allCurrencies),
    [allCurrencies],
  );
  const [fromCurrency, setFromCurrency] = useState(DEFAULT_FROM_CURRENCY);
  const [toCurrency, setToCurrency] = useState(DEFAULT_TO_CURRENCY);
  const [fromAmount, setFromAmount] = useState("");
  const [debouncedFromAmount, setDebouncedFromAmount] = useState("");

  useEffect(() => {
    const timer = setTimeout(() => setDebouncedFromAmount(fromAmount), 300);
    return () => clearTimeout(timer);
  }, [fromAmount]);

  useEffect(() => {
    if (exchangeCurrencies.length === 0) return;

    let nextFrom = fromCurrency;
    let nextTo = toCurrency;

    if (!nextFrom || !exchangeCurrencies.includes(nextFrom)) {
      nextFrom = exchangeCurrencies.includes(DEFAULT_FROM_CURRENCY)
        ? DEFAULT_FROM_CURRENCY
        : exchangeCurrencies[0];
    }

    const toOptions = getExchangeToOptions(nextFrom, allCurrencies);
    if (!nextTo || !toOptions.includes(nextTo)) {
      if (
        nextFrom === DEFAULT_FROM_CURRENCY &&
        toOptions.includes(DEFAULT_TO_CURRENCY)
      ) {
        nextTo = DEFAULT_TO_CURRENCY;
      } else {
        nextTo = toOptions[0] ?? "";
      }
    }

    if (nextFrom !== fromCurrency) setFromCurrency(nextFrom);
    if (nextTo !== toCurrency) setToCurrency(nextTo);
  }, [allCurrencies, exchangeCurrencies, fromCurrency, toCurrency]);

  const toOptions = useMemo(
    () => getExchangeToOptions(fromCurrency, allCurrencies),
    [fromCurrency, allCurrencies],
  );

  function handleFromCurrencyChange(value: string) {
    const nextToOptions = getExchangeToOptions(value, allCurrencies);
    setFromCurrency(value);
    if (toCurrency === value || !nextToOptions.includes(toCurrency)) {
      setToCurrency(nextToOptions[0] ?? "");
    }
  }

  function handleToCurrencyChange(value: string) {
    setToCurrency(value);
  }

  function handleSwapCurrencies() {
    if (!isSupportedExchangePair(toCurrency, fromCurrency)) return;
    setFromCurrency(toCurrency);
    setToCurrency(fromCurrency);
  }

  const previewRequired =
    isSupportedExchangePair(fromCurrency, toCurrency) && isPositiveAmount(fromAmount);

  const {
    data: preview,
    isLoading: previewLoading,
    isFetching: previewFetching,
    isSuccess: previewSuccess,
    isError: previewError,
  } = useExchangePreview({
    from_currency: fromCurrency,
    to_currency: toCurrency,
    from_amount: debouncedFromAmount,
  });

  const previewReady =
    previewRequired &&
    fromAmount.trim() === debouncedFromAmount.trim() &&
    previewSuccess &&
    !previewFetching &&
    preview != null &&
    preview.from_amount.trim() === debouncedFromAmount.trim();

  const orderMutation = useCreateExchangeOrder();
  const { data: ordersData, isLoading: ordersLoading } = useExchangeOrders();
  const { data: walletData } = useWalletBalances();
  const orders = ordersData?.orders || [];

  const wallets = walletData?.wallets ?? [];

  const fundingAvailable = useMemo(() => {
    const wallet = wallets.find(
      (item) => item.account_type === "FUNDING" && item.currency === fromCurrency,
    );
    return wallet?.available_balance ?? "0";
  }, [wallets, fromCurrency]);

  function handleConfirmOrder() {
    if (!previewReady) return;
    orderMutation.mutate(
      {
        from_currency: fromCurrency,
        to_currency: toCurrency,
        from_amount: fromAmount,
      },
      {
        onSuccess: () => {
          toast({ title: "下单成功", description: "兑换订单已创建。" });
          setFromAmount("");
          setDebouncedFromAmount("");
        },
        onError: (err) => {
          toast({ variant: "destructive", title: "下单失败", description: err.message });
        },
      }
    );
  }

  return (
    <div className="flex flex-col gap-6 md:h-[calc(100dvh-7.5rem)] md:min-h-0">
      <div className="shrink-0">
        <h1 className="text-2xl font-bold text-slate-900">兑换</h1>
        <p className="mt-1 text-sm text-slate-600">
          从资金账户卖出，买入币种入账到交易账户。
        </p>
      </div>

      <div className="grid min-h-0 flex-1 gap-6 md:grid-cols-2">
        <Card className="flex min-h-0 flex-col overflow-hidden">
          <CardHeader className="shrink-0">
            <CardTitle className="text-base">币种兑换</CardTitle>
          </CardHeader>
          <CardContent className="min-h-0 flex-1 space-y-4 overflow-y-auto">
            <div className="grid grid-cols-[1fr_auto_1fr] items-end gap-2">
              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">卖出</label>
                <SimpleSelect
                  value={fromCurrency}
                  onValueChange={handleFromCurrencyChange}
                  options={exchangeCurrencies}
                />
              </div>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                className="mb-0.5 shrink-0 text-slate-500 hover:text-slate-900"
                onClick={handleSwapCurrencies}
                disabled={!isSupportedExchangePair(toCurrency, fromCurrency)}
                aria-label="交换买入卖出币种"
              >
                <ArrowRightLeft className="h-5 w-5" />
              </Button>
              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">买入</label>
                <SimpleSelect
                  value={toCurrency}
                  onValueChange={handleToCurrencyChange}
                  options={toOptions}
                />
              </div>
            </div>

            <div>
              <div className="mb-1.5 flex items-center justify-between">
                <label className="text-sm font-medium text-slate-700">卖出金额</label>
                <span className="text-xs text-slate-500">
                  资金账户可用 {formatAmount(fundingAvailable, fromCurrency)} {fromCurrency}
                </span>
              </div>
              <div className="flex gap-2">
                <Input
                  type="text"
                  placeholder="0.00"
                  value={fromAmount}
                  onChange={(e) => setFromAmount(e.target.value)}
                  className="flex-1"
                />
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setFromAmount(fundingAvailable)}
                  disabled={parseFloat(fundingAvailable) <= 0}
                >
                  全部
                </Button>
              </div>
            </div>

            {previewRequired && (
              <ExchangePreviewSummary
                preview={preview}
                fromCurrency={fromCurrency}
                toCurrency={toCurrency}
                isLoading={previewLoading}
                isFetching={previewFetching}
                isError={previewError}
              />
            )}

            <Button
              className="w-full bg-blue-700 hover:bg-blue-800 text-white"
              onClick={handleConfirmOrder}
              disabled={orderMutation.isPending || !previewReady}
            >
              {orderMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              确认兑换
            </Button>
          </CardContent>
        </Card>

        <Card className="flex min-h-0 flex-col overflow-hidden">
          <CardHeader className="shrink-0">
            <div className="flex items-center justify-between gap-2">
              <CardTitle className="text-base">兑换记录</CardTitle>
              <Link
                href="/transactions?tab=exchange"
                className="inline-flex items-center text-sm font-medium text-primary hover:underline"
              >
                查看更多
                <ArrowRight className="ml-1 h-3.5 w-3.5" />
              </Link>
            </div>
          </CardHeader>
          <CardContent className="min-h-0 flex-1 overflow-y-auto">
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
                  <ExchangeOrderItem key={order.id} order={order} />
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
