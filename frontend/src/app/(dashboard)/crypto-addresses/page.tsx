"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { Loader2, Plus, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { SimpleSelect } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import {
  useAddCryptoAddress,
  useCryptoAddresses,
  useDeleteCryptoAddress,
} from "@/lib/hooks/use-addresses";
import { useSupportedCurrencies } from "@/lib/hooks/use-supported-currencies";
import { toCurrencyOptions } from "@/lib/utils/currency";
import { formatChainLabel, getWithdrawalNetworks } from "@/lib/utils/network";
import { toast } from "@/hooks/use-toast";

const emptyForm = {
  address: "",
  alias: "",
};

function maskAddress(address: string) {
  if (address.length <= 12) return address;
  return `${address.slice(0, 6)}...${address.slice(-6)}`;
}

export default function CryptoAddressesPage() {
  const { data: supportedCurrencies } = useSupportedCurrencies();
  const cryptoCurrencies = toCurrencyOptions(supportedCurrencies?.crypto_currencies ?? []);

  const [currency, setCurrency] = useState("");
  const [chain, setChain] = useState("");
  const [bindDialogOpen, setBindDialogOpen] = useState(false);
  const [form, setForm] = useState(emptyForm);

  useEffect(() => {
    const nextCurrency = supportedCurrencies?.crypto_currencies?.[0];
    if (!nextCurrency) return;
    if (!currency || !supportedCurrencies.crypto_currencies.includes(currency)) {
      setCurrency(nextCurrency);
      setChain(getWithdrawalNetworks(nextCurrency)[0]?.value || "");
    }
  }, [supportedCurrencies, currency]);

  const networks = getWithdrawalNetworks(currency);
  const { data: cryptoAddressesData, isLoading } = useCryptoAddresses();
  const cryptoAddresses = cryptoAddressesData ?? [];
  const addMutation = useAddCryptoAddress();
  const deleteMutation = useDeleteCryptoAddress();

  const canSubmit = form.address.trim() && form.alias.trim() && currency && chain;

  function openBindDialog() {
    setForm(emptyForm);
    setBindDialogOpen(true);
  }

  function closeBindDialog() {
    setBindDialogOpen(false);
    setForm(emptyForm);
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!canSubmit) return;

    addMutation.mutate(
      { currency, chain, address: form.address.trim(), alias: form.alias.trim() },
      {
        onSuccess: () => {
          toast({ title: "绑定成功", description: "提现白名单地址已添加。" });
          closeBindDialog();
        },
        onError: (err) => {
          toast({ variant: "destructive", title: "绑定失败", description: err.message });
        },
      }
    );
  }

  function handleDelete(id: number) {
    deleteMutation.mutate(id, {
      onSuccess: () => {
        toast({ title: "已解绑", description: "白名单地址已移除。" });
      },
      onError: (err) => {
        toast({ variant: "destructive", title: "解绑失败", description: err.message });
      },
    });
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">提现地址</h1>
          <p className="mt-1 text-sm text-slate-500">
            管理数币提现白名单地址，仅可向已绑定地址发起提现。
          </p>
        </div>
        {cryptoCurrencies.length > 0 && (
          <Button className="bg-blue-700 hover:bg-blue-800 text-white" onClick={openBindDialog}>
            <Plus className="mr-2 h-4 w-4" />
            添加白名单地址
          </Button>
        )}
      </div>

      <Dialog
        open={bindDialogOpen}
        onOpenChange={(open) => {
          if (open) setBindDialogOpen(true);
          else closeBindDialog();
        }}
      >
        <DialogContent className="flex max-h-[90vh] flex-col gap-0 overflow-visible p-0 sm:max-w-lg">
          <DialogHeader className="border-b px-6 py-4">
            <DialogTitle>绑定白名单地址</DialogTitle>
            <DialogDescription>
              地址绑定后需通过平台审核方可用于提现，请确保地址与网络匹配。
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={handleSubmit} className="flex min-h-0 flex-1 flex-col">
            <div className="grid gap-4 overflow-y-auto px-6 py-4">
              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">币种</label>
                <SimpleSelect
                  value={currency}
                  onValueChange={(value) => {
                    setCurrency(value);
                    setChain(getWithdrawalNetworks(value)[0]?.value || "");
                  }}
                  options={cryptoCurrencies}
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

              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">地址别名</label>
                <Input
                  value={form.alias}
                  onChange={(e) => setForm((prev) => ({ ...prev, alias: e.target.value }))}
                  placeholder="例如：主钱包"
                />
              </div>

              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">钱包地址</label>
                <Input
                  value={form.address}
                  onChange={(e) => setForm((prev) => ({ ...prev, address: e.target.value }))}
                  placeholder="请输入完整链上地址"
                />
              </div>
            </div>

            <DialogFooter className="border-t px-6 py-4">
              <Button type="button" variant="outline" onClick={closeBindDialog}>
                取消
              </Button>
              <Button
                type="submit"
                className="bg-blue-700 hover:bg-blue-800 text-white"
                disabled={addMutation.isPending || !canSubmit}
              >
                {addMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                确认绑定
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">已绑定地址</CardTitle>
        </CardHeader>
        <CardContent>
          {cryptoCurrencies.length === 0 ? (
            <p className="py-8 text-center text-sm text-slate-400">暂不支持数币提现</p>
          ) : isLoading ? (
            <div className="space-y-3">
              <Skeleton className="h-16 w-full" />
              <Skeleton className="h-16 w-full" />
            </div>
          ) : cryptoAddresses.length === 0 ? (
            <div className="py-8 text-center">
              <p className="text-sm text-slate-400">暂无白名单地址</p>
              <Button variant="link" className="mt-2 text-blue-700" onClick={openBindDialog}>
                立即添加
              </Button>
            </div>
          ) : (
            <div className="space-y-3">
              {cryptoAddresses.map((account) => (
                <div key={account.id} className="flex items-start justify-between rounded-lg border p-4">
                  <div className="min-w-0 space-y-1">
                    <p className="text-sm font-medium text-slate-900">
                      {account.alias} · {account.currency}
                    </p>
                    <p className="font-mono text-xs text-slate-600">{maskAddress(account.address)}</p>
                    <p className="text-xs text-slate-400">
                      {formatChainLabel(account.currency, account.chain)}
                    </p>
                  </div>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="shrink-0 text-slate-400 hover:text-red-600"
                    disabled={deleteMutation.isPending}
                    onClick={() => handleDelete(account.id)}
                  >
                    <Trash2 className="mr-1 h-4 w-4" />
                    解绑
                  </Button>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      <p className="text-sm text-slate-500">
        绑定完成后，可前往
        <Link href="/withdraw" className="mx-1 text-blue-700 hover:underline">
          提现
        </Link>
        页面发起数币提现。
      </p>
    </div>
  );
}
