"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { Switch } from "@/components/ui/switch";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useSystemCurrencyConfig, useUpdateSystemCurrencyConfig } from "@/lib/hooks/use-settings";
import { formatChainLabel } from "@/lib/utils/chain";
import { toast } from "@/hooks/use-toast";

function toggleItem(item: string, list: string[], setter: (next: string[]) => void) {
  setter(list.includes(item) ? list.filter((v) => v !== item) : [...list, item]);
}

export default function SettingsPage() {
  const { data, isLoading } = useSystemCurrencyConfig();
  const updateMutation = useUpdateSystemCurrencyConfig();

  const [selectedCrypto, setSelectedCrypto] = useState<string[]>([]);
  const [selectedFiat, setSelectedFiat] = useState<string[]>([]);
  const [selectedChains, setSelectedChains] = useState<Record<string, string[]>>({});
  const [defaultChains, setDefaultChains] = useState<Record<string, string>>({});

  useEffect(() => {
    if (!data) return;
    setSelectedCrypto(data.crypto_currencies ?? []);
    setSelectedFiat(data.fiat_currencies ?? []);
    setSelectedChains(data.crypto_chains ?? {});
    setDefaultChains(data.default_chains ?? {});
  }, [data]);

  const toggleChain = (currency: string, chain: string) => {
    const current = selectedChains[currency] ?? [];
    const next = current.includes(chain)
      ? current.filter((item) => item !== chain)
      : [...current, chain];
    setSelectedChains({ ...selectedChains, [currency]: next });
    if (!next.includes(defaultChains[currency] ?? "") && next.length > 0) {
      setDefaultChains({ ...defaultChains, [currency]: next[0] });
    }
  };

  const handleSave = () => {
    if (selectedCrypto.length === 0 || selectedFiat.length === 0) {
      toast({ title: "请至少选择一种数字货币和一种法币", variant: "destructive" });
      return;
    }
    for (const currency of data?.all_crypto ?? []) {
      if ((selectedChains[currency] ?? []).length === 0) {
        toast({ title: `${currency} 至少需要选择一条支持链`, variant: "destructive" });
        return;
      }
    }
    updateMutation.mutate(
      {
        crypto_currencies: selectedCrypto,
        fiat_currencies: selectedFiat,
        crypto_chains: selectedChains,
        default_chains: defaultChains,
      },
      {
        onSuccess: () => toast({ title: "全局币种与链配置已保存" }),
        onError: (error: Error) =>
          toast({ title: "保存失败", description: error.message, variant: "destructive" }),
      }
    );
  };

  if (isLoading) {
    return (
      <div className="space-y-4 p-6">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  return (
    <div className="space-y-6 p-6">
      <div>
        <h1 className="text-2xl font-semibold text-slate-900">系统设置</h1>
        <p className="mt-1 text-sm text-slate-500">
          币种开关决定商户默认启用项；支持链独立配置，并适用于所有数字货币。
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">币种与链配置</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <div>
            <p className="mb-3 text-sm font-medium text-slate-700">数字货币</p>
            <p className="mb-3 text-xs text-slate-500">开关仅控制商户默认是否启用该币种。</p>
            <div className="grid gap-3 sm:grid-cols-2 md:grid-cols-3">
              {(data?.all_crypto ?? []).map((currency) => (
                <div key={currency} className="flex items-center justify-between rounded-lg border px-3 py-2">
                  <span className="text-sm font-medium">{currency}</span>
                  <Switch
                    checked={selectedCrypto.includes(currency)}
                    onCheckedChange={() => toggleItem(currency, selectedCrypto, setSelectedCrypto)}
                  />
                </div>
              ))}
            </div>
          </div>

          <Separator />

          <div>
            <p className="mb-3 text-sm font-medium text-slate-700">法币</p>
            <div className="grid gap-3 sm:grid-cols-2 md:grid-cols-3">
              {(data?.all_fiat ?? []).map((currency) => (
                <div key={currency} className="flex items-center justify-between rounded-lg border px-3 py-2">
                  <span className="text-sm font-medium">{currency}</span>
                  <Switch
                    checked={selectedFiat.includes(currency)}
                    onCheckedChange={() => toggleItem(currency, selectedFiat, setSelectedFiat)}
                  />
                </div>
              ))}
            </div>
          </div>

          <Separator />

          <div className="space-y-4">
            <p className="text-sm font-medium text-slate-700">数字货币支持链与默认链</p>
            {(data?.all_crypto ?? []).map((currency) => {
              const catalog = data?.catalog_chains?.[currency] ?? [];
              const enabled = selectedChains[currency] ?? [];
              return (
                <div key={currency} className="rounded-lg border p-4 space-y-3">
                  <div className="flex flex-wrap items-center justify-between gap-3">
                    <p className="text-sm font-semibold text-slate-800">{currency}</p>
                    <div className="flex items-center gap-2">
                      <span className="text-xs text-slate-500">默认链</span>
                      <Select
                        value={defaultChains[currency] || enabled[0] || ""}
                        onValueChange={(value) =>
                          setDefaultChains({ ...defaultChains, [currency]: value })
                        }
                      >
                        <SelectTrigger className="w-[200px]">
                          <SelectValue placeholder="选择默认链" />
                        </SelectTrigger>
                        <SelectContent>
                          {enabled.map((chain) => (
                            <SelectItem key={chain} value={chain}>
                              {formatChainLabel(chain)}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                  </div>
                  <div className="grid gap-2 sm:grid-cols-2 md:grid-cols-3">
                    {catalog.map((chain) => (
                      <div key={chain} className="flex items-center justify-between rounded-md border px-3 py-2">
                        <span className="text-sm">{formatChainLabel(chain)}</span>
                        <Switch
                          checked={enabled.includes(chain)}
                          onCheckedChange={() => toggleChain(currency, chain)}
                        />
                      </div>
                    ))}
                  </div>
                </div>
              );
            })}
          </div>

          <div className="flex justify-end">
            <Button onClick={handleSave} disabled={updateMutation.isPending}>
              {updateMutation.isPending ? "保存中..." : "保存全局配置"}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
