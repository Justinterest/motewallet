"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { Loader2 } from "lucide-react";
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
import {
  useSubmitCryptoWithdrawal,
  useSubmitFiatWithdrawal,
  useWithdrawalOrders,
} from "@/lib/hooks/use-trading";
import { useBankAccounts } from "@/lib/hooks/use-addresses";
import { useSupportedCurrencies } from "@/lib/hooks/use-supported-currencies";
import { formatAmount } from "@/lib/utils/format";
import { toCurrencyOptions } from "@/lib/utils/currency";
import { TRANSFER_TYPE_LABELS } from "@/lib/utils/bank-account";
import { FIAT_WITHDRAWAL_PURPOSES } from "@/lib/utils/fiat-withdrawal";
import { toast } from "@/hooks/use-toast";
import type { BankAccount } from "@/types/address";

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

type WithdrawTab = "crypto" | "fiat";

export default function WithdrawPage() {
  const [activeTab, setActiveTab] = useState<WithdrawTab>("crypto");
  const { data: supportedCurrencies } = useSupportedCurrencies();
  const cryptoCurrencies = toCurrencyOptions(supportedCurrencies?.crypto_currencies ?? []);
  const fiatCurrencies = toCurrencyOptions(supportedCurrencies?.fiat_currencies ?? []);

  const [currency, setCurrency] = useState("");
  const [chain, setChain] = useState("");
  const [amount, setAmount] = useState("");
  const [toAddress, setToAddress] = useState("");

  const [fiatCurrency, setFiatCurrency] = useState("");
  const [fiatAmount, setFiatAmount] = useState("");
  const [bankAccountId, setBankAccountId] = useState("");
  const [purpose, setPurpose] = useState("OTHER");
  const [postscript, setPostscript] = useState("");

  useEffect(() => {
    const nextCurrency = supportedCurrencies?.crypto_currencies?.[0];
    if (!nextCurrency) return;
    if (!currency || !supportedCurrencies.crypto_currencies.includes(currency)) {
      setCurrency(nextCurrency);
      setChain(chains[nextCurrency]?.[0]?.value || "");
    }
  }, [supportedCurrencies, currency]);

  useEffect(() => {
    const nextCurrency = supportedCurrencies?.fiat_currencies?.[0];
    if (!nextCurrency) return;
    if (!fiatCurrency || !supportedCurrencies.fiat_currencies.includes(fiatCurrency)) {
      setFiatCurrency(nextCurrency);
    }
  }, [supportedCurrencies, fiatCurrency]);

  const { data: bankAccountsData, isLoading: bankAccountsLoading } = useBankAccounts();
  const bankAccounts = bankAccountsData ?? [];

  const accountsForCurrency = useMemo(
    () => bankAccounts.filter((a) => a.currency === fiatCurrency && a.status === "ACTIVE"),
    [bankAccounts, fiatCurrency]
  );

  const bankAccountOptions = useMemo(
    () =>
      accountsForCurrency.map((account) => ({
        value: String(account.id),
        label: formatBankAccountLabel(account),
      })),
    [accountsForCurrency]
  );

  useEffect(() => {
    if (accountsForCurrency.length === 0) {
      setBankAccountId("");
      return;
    }
    if (!bankAccountId || !accountsForCurrency.some((a) => String(a.id) === bankAccountId)) {
      setBankAccountId(String(accountsForCurrency[0].id));
    }
  }, [accountsForCurrency, bankAccountId]);

  const submitCryptoMutation = useSubmitCryptoWithdrawal();
  const submitFiatMutation = useSubmitFiatWithdrawal();
  const { data: ordersData, isLoading: ordersLoading } = useWithdrawalOrders();
  const orders = ordersData?.orders || [];

  function handleCryptoSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!amount || !toAddress) return;

    submitCryptoMutation.mutate(
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

  function handleFiatSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!fiatAmount || !bankAccountId || !postscript) return;

    submitFiatMutation.mutate(
      {
        currency: fiatCurrency,
        amount: fiatAmount,
        bank_account_id: Number(bankAccountId),
        purpose,
        postscript,
      },
      {
        onSuccess: () => {
          toast({ title: "提交成功", description: "法币提现申请已提交，等待审核。" });
          setFiatAmount("");
          setPostscript("");
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

      <div className="flex gap-2 rounded-lg border bg-slate-50 p-1 w-fit">
        <button
          type="button"
          onClick={() => setActiveTab("crypto")}
          className={`rounded-md px-4 py-2 text-sm font-medium transition-colors ${
            activeTab === "crypto" ? "bg-white text-slate-900 shadow-sm" : "text-slate-600 hover:text-slate-900"
          }`}
        >
          加密货币
        </button>
        <button
          type="button"
          onClick={() => setActiveTab("fiat")}
          className={`rounded-md px-4 py-2 text-sm font-medium transition-colors ${
            activeTab === "fiat" ? "bg-white text-slate-900 shadow-sm" : "text-slate-600 hover:text-slate-900"
          }`}
        >
          法币
        </button>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        {activeTab === "crypto" ? (
          <Card>
            <CardHeader>
              <CardTitle className="text-base">加密货币提现</CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleCryptoSubmit} className="space-y-4">
                <div>
                  <label className="mb-1.5 block text-sm font-medium text-slate-700">币种</label>
                  <SimpleSelect
                    value={currency}
                    onValueChange={(value) => {
                      setCurrency(value);
                      setChain(chains[value]?.[0]?.value || "");
                    }}
                    options={cryptoCurrencies}
                  />
                </div>

                <div>
                  <label className="mb-1.5 block text-sm font-medium text-slate-700">链/网络</label>
                  <SimpleSelect value={chain} onValueChange={setChain} options={chains[currency] || []} />
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
                  disabled={submitCryptoMutation.isPending || !amount || !toAddress}
                >
                  {submitCryptoMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  提交提现
                </Button>
              </form>
            </CardContent>
          </Card>
        ) : (
          <Card>
            <CardHeader>
              <CardTitle className="text-base">法币提现</CardTitle>
            </CardHeader>
            <CardContent>
              {fiatCurrencies.length === 0 ? (
                <p className="py-6 text-center text-sm text-slate-400">暂不支持法币提现</p>
              ) : (
                <form onSubmit={handleFiatSubmit} className="space-y-4">
                  <div>
                    <label className="mb-1.5 block text-sm font-medium text-slate-700">币种</label>
                    <SimpleSelect
                      value={fiatCurrency}
                      onValueChange={setFiatCurrency}
                      options={fiatCurrencies}
                    />
                  </div>

                  <div>
                    <div className="mb-1.5 flex items-center justify-between">
                      <label className="text-sm font-medium text-slate-700">收款银行账户</label>
                      <Link href="/bank-accounts" className="text-xs text-blue-700 hover:underline">
                        管理账户
                      </Link>
                    </div>
                    {bankAccountsLoading ? (
                      <Skeleton className="h-9 w-full" />
                    ) : accountsForCurrency.length === 0 ? (
                      <div className="rounded-md border border-dashed px-3 py-4 text-center text-sm text-slate-500">
                        <p>请先绑定 {fiatCurrency} 银行账户</p>
                        <Link
                          href="/bank-accounts"
                          className="mt-2 inline-block text-blue-700 hover:underline"
                        >
                          前往银行账户管理
                        </Link>
                      </div>
                    ) : (
                      <SimpleSelect
                        value={bankAccountId}
                        onValueChange={setBankAccountId}
                        options={bankAccountOptions}
                      />
                    )}
                  </div>

                  <div>
                    <label className="mb-1.5 block text-sm font-medium text-slate-700">提现金额</label>
                    <Input
                      type="text"
                      placeholder="0.00"
                      value={fiatAmount}
                      onChange={(e) => setFiatAmount(e.target.value)}
                    />
                  </div>

                  <div>
                    <label className="mb-1.5 block text-sm font-medium text-slate-700">转账目的</label>
                    <SimpleSelect
                      value={purpose}
                      onValueChange={setPurpose}
                      options={FIAT_WITHDRAWAL_PURPOSES.map((item) => ({
                        value: item.value,
                        label: item.label,
                      }))}
                    />
                  </div>

                  <div>
                    <label className="mb-1.5 block text-sm font-medium text-slate-700">转账附言</label>
                    <Input
                      value={postscript}
                      onChange={(e) => setPostscript(e.target.value)}
                      placeholder="请输入转账附言"
                    />
                  </div>

                  <Button
                    type="submit"
                    className="w-full bg-blue-700 hover:bg-blue-800 text-white"
                    disabled={
                      submitFiatMutation.isPending ||
                      !fiatAmount ||
                      !bankAccountId ||
                      !postscript ||
                      accountsForCurrency.length === 0
                    }
                  >
                    {submitFiatMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                    提交提现
                  </Button>
                </form>
              )}
            </CardContent>
          </Card>
        )}

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
                      <span
                        className={`text-xs font-medium ${
                          order.status === "COMPLETED"
                            ? "text-green-600"
                            : order.status === "FAILED"
                              ? "text-red-600"
                              : "text-amber-600"
                        }`}
                      >
                        {statusLabels[order.status] || order.status}
                      </span>
                    </div>
                    <div className="mt-1 flex items-center justify-between text-xs text-slate-500">
                      <span>
                        {order.type === "FIAT" ? "法币" : order.chain || "加密货币"} ·{" "}
                        {reviewLabels[order.review_status] || order.review_status}
                      </span>
                      <span>{new Date(order.created_at).toLocaleString("zh-CN")}</span>
                    </div>
                    {order.platform_fee && order.platform_fee !== "0" && (
                      <p className="mt-1 text-xs text-slate-400">
                        手续费：{order.platform_fee} {order.currency}
                      </p>
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

function formatBankAccountLabel(account: BankAccount) {
  const transfer = TRANSFER_TYPE_LABELS[account.transfer_type] || account.transfer_type;
  return `${account.bank_name} · ${account.account_no_masked} · ${transfer}`;
}
