"use client";

import { useEffect, useState } from "react";
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
import { KycCountrySelect } from "@/components/kyc/kyc-country-select";
import {
  useAddBankAccount,
  useBankAccounts,
  useDeleteBankAccount,
} from "@/lib/hooks/use-addresses";
import { useSupportedCurrencies } from "@/lib/hooks/use-supported-currencies";
import { toCurrencyOptions } from "@/lib/utils/currency";
import { getTransferTypesForCurrency, TRANSFER_TYPE_LABELS } from "@/lib/utils/bank-account";
import { toast } from "@/hooks/use-toast";

const emptyBankForm = {
  bank_name: "",
  bank_country: "",
  swift_code: "",
  account_name: "",
  account_no: "",
  transfer_type: "LOCAL",
  account_type: "",
  payee_address: "",
  payee_address_second: "",
  bank_code: "",
  bank_address: "",
  middle_swift_code: "",
};

function createEmptyBankForm(currency: string) {
  return {
    ...emptyBankForm,
    transfer_type: getTransferTypesForCurrency(currency)[0]?.value || "LOCAL",
  };
}

export default function BankAccountsPage() {
  const { data: supportedCurrencies } = useSupportedCurrencies();
  const fiatCurrencies = toCurrencyOptions(supportedCurrencies?.fiat_currencies ?? []);

  const [currency, setCurrency] = useState("");
  const [bindDialogOpen, setBindDialogOpen] = useState(false);
  const [bankForm, setBankForm] = useState(emptyBankForm);

  useEffect(() => {
    const nextCurrency = supportedCurrencies?.fiat_currencies?.[0];
    if (!nextCurrency) return;
    if (!currency || !supportedCurrencies.fiat_currencies.includes(currency)) {
      setCurrency(nextCurrency);
    }
  }, [supportedCurrencies, currency]);

  const { data: bankAccountsData, isLoading } = useBankAccounts();
  const bankAccounts = bankAccountsData ?? [];
  const addBankMutation = useAddBankAccount();
  const deleteBankMutation = useDeleteBankAccount();

  const transferTypeOptions = getTransferTypesForCurrency(currency);
  const needsSwiftCode = bankForm.transfer_type === "TT" || bankForm.transfer_type === "CHATS";
  const needsWireFields = needsSwiftCode;
  const needsBankCode = bankForm.transfer_type === "CHATS";
  const needsAccountType = bankForm.transfer_type === "TT";

  const canSubmitBind =
    bankForm.account_no &&
    bankForm.account_name &&
    bankForm.bank_name &&
    bankForm.bank_country &&
    (!needsSwiftCode || bankForm.swift_code) &&
    (!needsWireFields || (bankForm.payee_address && bankForm.payee_address_second)) &&
    (!needsBankCode || bankForm.bank_code) &&
    (!needsAccountType || bankForm.account_type);

  function openBindDialog() {
    setBankForm(createEmptyBankForm(currency));
    setBindDialogOpen(true);
  }

  function closeBindDialog() {
    setBindDialogOpen(false);
    setBankForm(createEmptyBankForm(currency));
  }

  function handleBindBankAccount(e: React.FormEvent) {
    e.preventDefault();
    if (!canSubmitBind) return;

    addBankMutation.mutate(
      { currency, ...bankForm, payee_country_code: bankForm.bank_country },
      {
        onSuccess: () => {
          toast({ title: "绑定成功", description: "银行账户已添加。" });
          closeBindDialog();
        },
        onError: (err) => {
          toast({ variant: "destructive", title: "绑定失败", description: err.message });
        },
      }
    );
  }

  function handleDeleteBankAccount(id: number) {
    deleteBankMutation.mutate(id, {
      onSuccess: () => {
        toast({ title: "已解绑", description: "银行账户已移除。" });
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
          <h1 className="text-2xl font-bold text-slate-900">银行账户</h1>
          <p className="mt-1 text-sm text-slate-500">管理法币提现收款账户，法币提现前需先绑定对应币种的银行账户。</p>
        </div>
        {fiatCurrencies.length > 0 && (
          <Button className="bg-blue-700 hover:bg-blue-800 text-white" onClick={openBindDialog}>
            <Plus className="mr-2 h-4 w-4" />
            添加银行账户
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
        <DialogContent className="flex max-h-[90vh] flex-col gap-0 overflow-hidden p-0 sm:max-w-2xl">
          <DialogHeader className="border-b px-6 py-4">
            <DialogTitle>绑定银行账户</DialogTitle>
            <DialogDescription>填写银行信息并提交绑定，绑定后可于法币提现页选择该账户。</DialogDescription>
          </DialogHeader>

          <form onSubmit={handleBindBankAccount} className="flex min-h-0 flex-1 flex-col">
            <div className="grid max-h-[calc(90vh-9rem)] gap-4 overflow-y-auto px-6 py-4 sm:grid-cols-2">
              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">币种</label>
                <SimpleSelect
                  value={currency}
                  onValueChange={(value) => {
                    setCurrency(value);
                    setBankForm((prev) => ({
                      ...prev,
                      transfer_type: getTransferTypesForCurrency(value)[0]?.value || "LOCAL",
                      bank_country: "",
                    }));
                  }}
                  options={fiatCurrencies}
                />
              </div>

              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">转账类型</label>
                <SimpleSelect
                  value={bankForm.transfer_type}
                  onValueChange={(value) => setBankForm((prev) => ({ ...prev, transfer_type: value }))}
                  options={transferTypeOptions}
                />
              </div>

              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">银行名称</label>
                <Input
                  value={bankForm.bank_name}
                  onChange={(e) => setBankForm((prev) => ({ ...prev, bank_name: e.target.value }))}
                  placeholder="请输入银行名称"
                />
              </div>

              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">银行所在国家/地区</label>
                <KycCountrySelect
                  scene="WITHDRAWAL"
                  currency={currency}
                  value={bankForm.bank_country}
                  onChange={(value) => setBankForm((prev) => ({ ...prev, bank_country: value }))}
                />
              </div>

              {needsSwiftCode && (
                <div>
                  <label className="mb-1.5 block text-sm font-medium text-slate-700">SWIFT Code</label>
                  <Input
                    value={bankForm.swift_code}
                    onChange={(e) => setBankForm((prev) => ({ ...prev, swift_code: e.target.value }))}
                    placeholder="请输入 SWIFT Code"
                  />
                </div>
              )}

              {needsAccountType && (
                <div>
                  <label className="mb-1.5 block text-sm font-medium text-slate-700">账户类型</label>
                  <Input
                    value={bankForm.account_type}
                    onChange={(e) => setBankForm((prev) => ({ ...prev, account_type: e.target.value }))}
                    placeholder="TT 电汇必填"
                  />
                </div>
              )}

              {needsBankCode && (
                <div>
                  <label className="mb-1.5 block text-sm font-medium text-slate-700">银行代码</label>
                  <Input
                    value={bankForm.bank_code}
                    onChange={(e) => setBankForm((prev) => ({ ...prev, bank_code: e.target.value }))}
                    placeholder="CHATS 必填"
                  />
                </div>
              )}

              {needsWireFields && (
                <>
                  <div>
                    <label className="mb-1.5 block text-sm font-medium text-slate-700">收款方地址 1</label>
                    <Input
                      value={bankForm.payee_address}
                      onChange={(e) => setBankForm((prev) => ({ ...prev, payee_address: e.target.value }))}
                      placeholder="不支持中文"
                    />
                  </div>
                  <div>
                    <label className="mb-1.5 block text-sm font-medium text-slate-700">收款方地址 2</label>
                    <Input
                      value={bankForm.payee_address_second}
                      onChange={(e) => setBankForm((prev) => ({ ...prev, payee_address_second: e.target.value }))}
                      placeholder="不支持中文"
                    />
                  </div>
                  <div>
                    <label className="mb-1.5 block text-sm font-medium text-slate-700">银行地址</label>
                    <Input
                      value={bankForm.bank_address}
                      onChange={(e) => setBankForm((prev) => ({ ...prev, bank_address: e.target.value }))}
                      placeholder="选填"
                    />
                  </div>
                  <div>
                    <label className="mb-1.5 block text-sm font-medium text-slate-700">中间行 SWIFT Code</label>
                    <Input
                      value={bankForm.middle_swift_code}
                      onChange={(e) => setBankForm((prev) => ({ ...prev, middle_swift_code: e.target.value }))}
                      placeholder="选填"
                    />
                  </div>
                </>
              )}

              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">账户名</label>
                <Input
                  value={bankForm.account_name}
                  onChange={(e) => setBankForm((prev) => ({ ...prev, account_name: e.target.value }))}
                  placeholder="请输入账户名"
                />
              </div>

              <div>
                <label className="mb-1.5 block text-sm font-medium text-slate-700">银行账号</label>
                <Input
                  value={bankForm.account_no}
                  onChange={(e) => setBankForm((prev) => ({ ...prev, account_no: e.target.value }))}
                  placeholder="请输入银行账号"
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
                disabled={addBankMutation.isPending || !canSubmitBind}
              >
                {addBankMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                确认绑定
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">已绑定账户</CardTitle>
        </CardHeader>
        <CardContent>
          {fiatCurrencies.length === 0 ? (
            <p className="py-8 text-center text-sm text-slate-400">暂不支持法币账户绑定</p>
          ) : isLoading ? (
            <div className="space-y-3">
              <Skeleton className="h-16 w-full" />
              <Skeleton className="h-16 w-full" />
            </div>
          ) : bankAccounts.length === 0 ? (
            <div className="py-8 text-center">
              <p className="text-sm text-slate-400">暂无绑定的银行账户</p>
              <Button variant="link" className="mt-2 text-blue-700" onClick={openBindDialog}>
                立即添加
              </Button>
            </div>
          ) : (
            <div className="space-y-3">
              {bankAccounts.map((account) => (
                <div key={account.id} className="flex items-start justify-between rounded-lg border p-4">
                  <div className="min-w-0 space-y-1">
                    <p className="text-sm font-medium text-slate-900">
                      {account.bank_name} · {account.currency}
                    </p>
                    <p className="text-sm text-slate-600">{account.account_name}</p>
                    <p className="text-xs text-slate-500">{account.account_no_masked}</p>
                    <p className="text-xs text-slate-400">
                      {TRANSFER_TYPE_LABELS[account.transfer_type] || account.transfer_type}
                      {account.swift_code ? ` · ${account.swift_code}` : ""}
                    </p>
                  </div>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="shrink-0 text-slate-400 hover:text-red-600"
                    disabled={deleteBankMutation.isPending}
                    onClick={() => handleDeleteBankAccount(account.id)}
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
        页面发起法币提现。
      </p>
    </div>
  );
}
