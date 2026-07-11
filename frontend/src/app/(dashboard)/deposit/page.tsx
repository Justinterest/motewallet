"use client";

import { useEffect, useState } from "react";
import { Copy, Check } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { SimpleSelect } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { useDepositAddresses, useDepositOrders } from "@/lib/hooks/use-trading";
import { useSupportedCurrencies } from "@/lib/hooks/use-supported-currencies";
import { formatAmount } from "@/lib/utils/format";
import { toCurrencyOptions } from "@/lib/utils/currency";
import { getDepositNetworks, resolveDefaultChain, formatDepositStatus } from "@/lib/utils/network";

export default function DepositPage() {
  const { data: supportedCurrencies } = useSupportedCurrencies();
  const currencies = toCurrencyOptions(supportedCurrencies?.crypto_currencies ?? []);
  const [currency, setCurrency] = useState("");
  const [chain, setChain] = useState("");
  const [copied, setCopied] = useState(false);

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
  const { data: ordersData, isLoading: ordersLoading } = useDepositOrders(currency, chain);

  const orders = ordersData?.orders || [];
  const networks = getDepositNetworks(
    currency,
    supportedCurrencies?.crypto_chains?.[currency]
  );

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
            <CardTitle className="text-base">充值地址</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
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
              <Skeleton className="h-20 w-full" />
            ) : addressData?.address ? (
              <div className="rounded-lg border bg-slate-50 p-4">
                <p className="mb-1 text-xs text-slate-500">收款地址</p>
                <p className="break-all font-mono text-sm text-slate-900">{addressData.address}</p>
                {addressData.network && (
                  <p className="mt-1 text-xs text-slate-500">网络：{addressData.network}</p>
                )}
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

            <div className="rounded-lg border border-amber-200 bg-amber-50 p-3 text-xs text-amber-800">
              <p>请仅向该地址转入所选币种，并确认网络与上方选择一致。</p>
              <p className="mt-1">到账后余额将自动更新，一般需要等待区块确认。</p>
            </div>
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
                      <p className="text-xs text-slate-500">
                        {order.network} · {new Date(order.created_at).toLocaleString("zh-CN")}
                      </p>
                      {order.tx_hash && (
                        <p className="mt-0.5 truncate text-xs text-slate-400" title={order.tx_hash}>
                          交易哈希：{order.tx_hash}
                        </p>
                      )}
                    </div>
                    <span
                      className={`text-xs font-medium ${
                        order.status === "COMPLETED" || order.status === "SUCCESS"
                          ? "text-green-600"
                          : order.status === "FAILED"
                            ? "text-red-600"
                            : "text-amber-600"
                      }`}
                    >
                      {formatDepositStatus(order.status)}
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
