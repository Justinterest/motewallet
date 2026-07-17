"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { ArrowRight } from "lucide-react";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { SimpleSelect } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { DepositAddressCard } from "@/components/deposit/deposit-address-card";
import { DepositOrderItem } from "@/components/deposit/deposit-order-item";
import { useDepositAddresses, useDepositOrders } from "@/lib/hooks/use-trading";
import { useSupportedCurrencies } from "@/lib/hooks/use-supported-currencies";
import { toCurrencyOptions } from "@/lib/utils/currency";
import { formatChainLabel, getDepositNetworks, resolveDefaultChain } from "@/lib/utils/network";

export default function DepositPage() {
  const { data: supportedCurrencies } = useSupportedCurrencies();
  const currencies = toCurrencyOptions(supportedCurrencies?.crypto_currencies ?? []);
  const [currency, setCurrency] = useState("");
  const [chain, setChain] = useState("");

  useEffect(() => {
    const nextCurrency = supportedCurrencies?.crypto_currencies?.[0];
    if (!nextCurrency) return;
    if (!currency || !supportedCurrencies.crypto_currencies.includes(currency)) {
      setCurrency(nextCurrency);
      const networks = getDepositNetworks(
        nextCurrency,
        supportedCurrencies.crypto_chains?.[nextCurrency]
      );
      setChain(
        resolveDefaultChain(nextCurrency, networks, supportedCurrencies.default_chains)
      );
    }
  }, [supportedCurrencies, currency]);

  const { data: addressData, isLoading: addressLoading } = useDepositAddresses(currency, chain);
  // 充值记录不按网络过滤，展示全部最近记录
  const { data: ordersData, isLoading: ordersLoading } = useDepositOrders();

  const orders = ordersData?.orders || [];
  const networks = getDepositNetworks(
    currency,
    supportedCurrencies?.crypto_chains?.[currency]
  );
  const selectedNetworkLabel =
    networks.find((item) => item.value === chain)?.label ||
    formatChainLabel(currency, chain || addressData?.network);

  return (
    <div className="flex flex-col gap-6 md:h-[calc(100dvh-7.5rem)] md:min-h-0">
      <div className="shrink-0">
        <h1 className="text-2xl font-bold text-slate-900">充值</h1>
        <p className="mt-1 text-sm text-slate-600">
          向下方地址转入数字货币，确认后自动入账资金账户。
        </p>
      </div>

      <div className="grid min-h-0 flex-1 gap-6 md:grid-cols-2">
        <Card className="flex min-h-0 flex-col overflow-hidden">
          <CardHeader className="shrink-0">
            <CardTitle className="text-base">充值地址</CardTitle>
          </CardHeader>
          <CardContent className="min-h-0 flex-1 space-y-4 overflow-y-auto">
            <div>
              <label className="mb-1.5 block text-sm font-medium text-slate-700">币种</label>
              <SimpleSelect
                value={currency}
                onValueChange={(value) => {
                  setCurrency(value);
                  const nextNetworks = getDepositNetworks(
                    value,
                    supportedCurrencies?.crypto_chains?.[value]
                  );
                  setChain(
                    resolveDefaultChain(value, nextNetworks, supportedCurrencies?.default_chains)
                  );
                }}
                options={currencies}
              />
            </div>

            <div>
              <label className="mb-1.5 block text-sm font-medium text-slate-700">网络</label>
              <SimpleSelect
                value={chain}
                onValueChange={setChain}
                options={networks}
              />
            </div>

            {addressLoading ? (
              <div className="space-y-3 rounded-lg border bg-slate-50 p-4">
                <Skeleton className="mx-auto h-[168px] w-[168px] rounded-xl" />
                <Skeleton className="h-16 w-full" />
                <Skeleton className="h-10 w-full" />
              </div>
            ) : addressData?.address ? (
              <DepositAddressCard
                address={addressData.address}
                currency={currency}
                network={selectedNetworkLabel}
              />
            ) : (
              <p className="text-sm text-slate-400">当前币种/网络暂无可用充值地址</p>
            )}

            <div className="rounded-lg border border-amber-200 bg-amber-50 p-3 text-xs leading-relaxed text-amber-800">
              <p className="font-medium">充值须知</p>
              <ul className="mt-1.5 list-disc space-y-1 pl-4">
                <li>
                  请仅转入 <span className="font-medium">{currency || "所选币种"}</span>
                  ，并确认网络为{" "}
                  <span className="font-medium">{selectedNetworkLabel || "上方所选网络"}</span>
                  。
                </li>
                <li>转入其他币种或错误网络可能导致资产无法找回。</li>
                <li>到账需等待链上确认，确认完成后余额会自动更新，无需手动提交。</li>
              </ul>
            </div>
          </CardContent>
        </Card>

        <Card className="flex min-h-0 flex-col overflow-hidden">
          <CardHeader className="shrink-0">
            <div className="flex items-center justify-between gap-2">
              <CardTitle className="text-base">充值记录</CardTitle>
              <Link
                href="/transactions?tab=deposit"
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
              <p className="py-8 text-center text-sm text-slate-400">暂无充值记录</p>
            ) : (
              <div className="space-y-3">
                {orders.map((order) => (
                  <DepositOrderItem key={order.id} order={order} />
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
