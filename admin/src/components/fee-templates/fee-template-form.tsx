"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Plus, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  feeTemplateFormSchema,
  type FeeTemplateFormValues,
} from "@/lib/validations/fee-template";
import type {
  FeeTemplateExchangeItem,
  FeeTemplateCryptoWithdrawalItem,
  FeeTemplateFiatWithdrawalItem,
  CreateFeeTemplateRequest,
} from "@/types/fee-template";

const CURRENCY_OPTIONS = ["USDT", "USDC", "BTC", "USD", "HKD", "EUR"];
const CHAIN_OPTIONS = ["ETH", "TRX", "BTC"];
const TRANSFER_TYPE_OPTIONS = ["LOCAL", "CHATS", "TT"];

interface FeeTemplateFormProps {
  defaultValues?: FeeTemplateFormValues;
  defaultExchangeItems?: FeeTemplateExchangeItem[];
  defaultCryptoItems?: FeeTemplateCryptoWithdrawalItem[];
  defaultFiatItems?: FeeTemplateFiatWithdrawalItem[];
  onSubmit: (data: CreateFeeTemplateRequest) => void;
  isPending: boolean;
  submitLabel: string;
}

export function FeeTemplateForm({
  defaultValues,
  defaultExchangeItems,
  defaultCryptoItems,
  defaultFiatItems,
  onSubmit,
  isPending,
  submitLabel,
}: FeeTemplateFormProps) {
  const form = useForm<FeeTemplateFormValues>({
    resolver: zodResolver(feeTemplateFormSchema),
    defaultValues: defaultValues || {
      name: "",
      description: "",
      is_default: false,
    },
  });

  const [exchangeItems, setExchangeItems] = useState<FeeTemplateExchangeItem[]>(
    defaultExchangeItems || []
  );
  const [cryptoItems, setCryptoItems] = useState<FeeTemplateCryptoWithdrawalItem[]>(
    defaultCryptoItems || []
  );
  const [fiatItems, setFiatItems] = useState<FeeTemplateFiatWithdrawalItem[]>(
    defaultFiatItems || []
  );

  const addExchangeItem = () => {
    setExchangeItems([
      ...exchangeItems,
      { from_currency: "USDT", to_currency: "USD", fee_rate: "0", min_fee: "0", min_fee_currency: "USD" },
    ]);
  };

  const removeExchangeItem = (index: number) => {
    setExchangeItems(exchangeItems.filter((_, i) => i !== index));
  };

  const updateExchangeItem = (index: number, field: keyof FeeTemplateExchangeItem, value: string) => {
    setExchangeItems(
      exchangeItems.map((item, i) => (i === index ? { ...item, [field]: value } : item))
    );
  };

  const addCryptoItem = () => {
    setCryptoItems([
      ...cryptoItems,
      { currency: "USDT", chain: "ETH", fee_rate: "0", fixed_fee: "0" },
    ]);
  };

  const removeCryptoItem = (index: number) => {
    setCryptoItems(cryptoItems.filter((_, i) => i !== index));
  };

  const updateCryptoItem = (index: number, field: keyof FeeTemplateCryptoWithdrawalItem, value: string) => {
    setCryptoItems(
      cryptoItems.map((item, i) => (i === index ? { ...item, [field]: value } : item))
    );
  };

  const addFiatItem = () => {
    setFiatItems([
      ...fiatItems,
      { currency: "USD", transfer_type: "LOCAL", fee_rate: "0", fixed_fee: "0" },
    ]);
  };

  const removeFiatItem = (index: number) => {
    setFiatItems(fiatItems.filter((_, i) => i !== index));
  };

  const updateFiatItem = (index: number, field: keyof FeeTemplateFiatWithdrawalItem, value: string) => {
    setFiatItems(
      fiatItems.map((item, i) => (i === index ? { ...item, [field]: value } : item))
    );
  };

  const handleSubmit = (values: FeeTemplateFormValues) => {
    const data: CreateFeeTemplateRequest = {
      name: values.name,
      description: values.description,
      is_default: values.is_default,
      exchange_items: exchangeItems.map((item) => ({ from_currency: item.from_currency, to_currency: item.to_currency, fee_rate: item.fee_rate, min_fee: item.min_fee, min_fee_currency: item.min_fee_currency })),
      crypto_withdrawal_items: cryptoItems.map((item) => ({ currency: item.currency, chain: item.chain, fee_rate: item.fee_rate, fixed_fee: item.fixed_fee })),
      fiat_withdrawal_items: fiatItems.map((item) => ({ currency: item.currency, transfer_type: item.transfer_type, fee_rate: item.fee_rate, fixed_fee: item.fixed_fee })),
    };
    onSubmit(data);
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-6">
        {/* Basic info */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">基本信息</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>模板名称</FormLabel>
                  <FormControl>
                    <Input placeholder="请输入模板名称" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>描述</FormLabel>
                  <FormControl>
                    <Textarea placeholder="请输入模板描述（可选）" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="is_default"
              render={({ field }) => (
                <FormItem className="flex items-center gap-3">
                  <FormLabel className="mt-0">设为默认模板</FormLabel>
                  <FormControl>
                    <Switch
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </CardContent>
        </Card>

        {/* Fee items tabs */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">费率配置</CardTitle>
          </CardHeader>
          <CardContent>
            <Tabs defaultValue="exchange">
              <TabsList className="mb-4">
                <TabsTrigger value="exchange">兑换费率</TabsTrigger>
                <TabsTrigger value="crypto">数币提现费率</TabsTrigger>
                <TabsTrigger value="fiat">法币提现费率</TabsTrigger>
              </TabsList>

              {/* Exchange fee items */}
              <TabsContent value="exchange" className="space-y-4">
                {exchangeItems.map((item, index) => (
                  <div
                    key={index}
                    className="flex items-end gap-3 rounded-lg border border-slate-200 bg-slate-50 p-4"
                  >
                    <div className="flex-1 space-y-1">
                      <Label className="text-xs text-slate-500">源币种</Label>
                      <Select
                        value={item.from_currency}
                        onValueChange={(v) => updateExchangeItem(index, "from_currency", v)}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {CURRENCY_OPTIONS.map((c) => (
                            <SelectItem key={c} value={c}>{c}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="flex-1 space-y-1">
                      <Label className="text-xs text-slate-500">目标币种</Label>
                      <Select
                        value={item.to_currency}
                        onValueChange={(v) => updateExchangeItem(index, "to_currency", v)}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {CURRENCY_OPTIONS.map((c) => (
                            <SelectItem key={c} value={c}>{c}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="flex-1 space-y-1">
                      <Label className="text-xs text-slate-500">费率</Label>
                      <Input
                        value={item.fee_rate}
                        onChange={(e) => updateExchangeItem(index, "fee_rate", e.target.value)}
                        placeholder="0.01"
                      />
                    </div>
                    <div className="flex-1 space-y-1">
                      <Label className="text-xs text-slate-500">最低手续费</Label>
                      <Input
                        value={item.min_fee}
                        onChange={(e) => updateExchangeItem(index, "min_fee", e.target.value)}
                        placeholder="0"
                      />
                    </div>
                    <div className="flex-1 space-y-1">
                      <Label className="text-xs text-slate-500">最低手续费币种</Label>
                      <Select
                        value={item.min_fee_currency}
                        onValueChange={(v) => updateExchangeItem(index, "min_fee_currency", v)}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {CURRENCY_OPTIONS.map((c) => (
                            <SelectItem key={c} value={c}>{c}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      onClick={() => removeExchangeItem(index)}
                    >
                      <Trash2 className="size-4 text-red-500" />
                    </Button>
                  </div>
                ))}
                <Button type="button" variant="outline" onClick={addExchangeItem}>
                  <Plus className="size-4" />
                  添加兑换费率
                </Button>
              </TabsContent>

              {/* Crypto withdrawal fee items */}
              <TabsContent value="crypto" className="space-y-4">
                {cryptoItems.map((item, index) => (
                  <div
                    key={index}
                    className="flex items-end gap-3 rounded-lg border border-slate-200 bg-slate-50 p-4"
                  >
                    <div className="flex-1 space-y-1">
                      <Label className="text-xs text-slate-500">币种</Label>
                      <Select
                        value={item.currency}
                        onValueChange={(v) => updateCryptoItem(index, "currency", v)}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {CURRENCY_OPTIONS.map((c) => (
                            <SelectItem key={c} value={c}>{c}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="flex-1 space-y-1">
                      <Label className="text-xs text-slate-500">链</Label>
                      <Select
                        value={item.chain}
                        onValueChange={(v) => updateCryptoItem(index, "chain", v)}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {CHAIN_OPTIONS.map((c) => (
                            <SelectItem key={c} value={c}>{c}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="flex-1 space-y-1">
                      <Label className="text-xs text-slate-500">费率</Label>
                      <Input
                        value={item.fee_rate}
                        onChange={(e) => updateCryptoItem(index, "fee_rate", e.target.value)}
                        placeholder="0.001"
                      />
                    </div>
                    <div className="flex-1 space-y-1">
                      <Label className="text-xs text-slate-500">固定手续费</Label>
                      <Input
                        value={item.fixed_fee}
                        onChange={(e) => updateCryptoItem(index, "fixed_fee", e.target.value)}
                        placeholder="0"
                      />
                    </div>
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      onClick={() => removeCryptoItem(index)}
                    >
                      <Trash2 className="size-4 text-red-500" />
                    </Button>
                  </div>
                ))}
                <Button type="button" variant="outline" onClick={addCryptoItem}>
                  <Plus className="size-4" />
                  添加数币提现费率
                </Button>
              </TabsContent>

              {/* Fiat withdrawal fee items */}
              <TabsContent value="fiat" className="space-y-4">
                {fiatItems.map((item, index) => (
                  <div
                    key={index}
                    className="flex items-end gap-3 rounded-lg border border-slate-200 bg-slate-50 p-4"
                  >
                    <div className="flex-1 space-y-1">
                      <Label className="text-xs text-slate-500">币种</Label>
                      <Select
                        value={item.currency}
                        onValueChange={(v) => updateFiatItem(index, "currency", v)}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {CURRENCY_OPTIONS.map((c) => (
                            <SelectItem key={c} value={c}>{c}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="flex-1 space-y-1">
                      <Label className="text-xs text-slate-500">转账类型</Label>
                      <Select
                        value={item.transfer_type}
                        onValueChange={(v) => updateFiatItem(index, "transfer_type", v)}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {TRANSFER_TYPE_OPTIONS.map((t) => (
                            <SelectItem key={t} value={t}>{t}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="flex-1 space-y-1">
                      <Label className="text-xs text-slate-500">费率</Label>
                      <Input
                        value={item.fee_rate}
                        onChange={(e) => updateFiatItem(index, "fee_rate", e.target.value)}
                        placeholder="0.01"
                      />
                    </div>
                    <div className="flex-1 space-y-1">
                      <Label className="text-xs text-slate-500">固定手续费</Label>
                      <Input
                        value={item.fixed_fee}
                        onChange={(e) => updateFiatItem(index, "fixed_fee", e.target.value)}
                        placeholder="0"
                      />
                    </div>
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      onClick={() => removeFiatItem(index)}
                    >
                      <Trash2 className="size-4 text-red-500" />
                    </Button>
                  </div>
                ))}
                <Button type="button" variant="outline" onClick={addFiatItem}>
                  <Plus className="size-4" />
                  添加法币提现费率
                </Button>
              </TabsContent>
            </Tabs>
          </CardContent>
        </Card>

        <div className="flex justify-end">
          <Button type="submit" disabled={isPending}>
            {isPending ? "提交中..." : submitLabel}
          </Button>
        </div>
      </form>
    </Form>
  );
}
