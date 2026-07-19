"use client";

import { QRCodeSVG } from "qrcode.react";
import { Loader2, Copy, Check } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { totpCodeSchema, type TotpCodeFormValues } from "@/lib/validations/auth";

interface TotpSetupPanelProps {
  totpUri: string;
  totpSecret: string;
  title?: string;
  description?: string;
  isPending?: boolean;
  errorMessage?: string | null;
  onSubmit: (code: string) => void;
  onBack?: () => void;
}

export function TotpSetupPanel({
  totpUri,
  totpSecret,
  title = "绑定两步验证",
  description = "请使用 Google Authenticator、Microsoft Authenticator 或其他验证器 App 扫描二维码，并输入 6 位验证码完成绑定。",
  isPending,
  errorMessage,
  onSubmit,
  onBack,
}: TotpSetupPanelProps) {
  const [copied, setCopied] = useState(false);
  const form = useForm<TotpCodeFormValues>({
    resolver: zodResolver(totpCodeSchema),
    defaultValues: { code: "" },
  });

  async function handleCopy() {
    try {
      await navigator.clipboard.writeText(totpSecret);
      setCopied(true);
      window.setTimeout(() => setCopied(false), 2000);
    } catch {
      // ignore
    }
  }

  return (
    <div className="space-y-4">
      <div className="text-center">
        <h2 className="text-xl font-semibold">{title}</h2>
        <p className="mt-2 text-sm text-slate-500">{description}</p>
      </div>

      <div className="flex flex-col items-center gap-3">
        <div className="rounded-lg border bg-white p-3">
          <QRCodeSVG value={totpUri} size={180} level="M" includeMargin={false} />
        </div>
        <div className="flex w-full max-w-sm items-center gap-2 rounded-md border bg-slate-50 px-3 py-2">
          <code className="flex-1 break-all text-xs text-slate-700">{totpSecret}</code>
          <button
            type="button"
            onClick={handleCopy}
            className="shrink-0 text-slate-500 hover:text-slate-700"
            aria-label="复制密钥"
          >
            {copied ? <Check className="h-4 w-4 text-green-600" /> : <Copy className="h-4 w-4" />}
          </button>
        </div>
        <p className="text-xs text-slate-400">无法扫码时可手动输入上方密钥</p>
      </div>

      {errorMessage && (
        <div className="rounded-md bg-red-50 p-3 text-sm text-red-600">{errorMessage}</div>
      )}

      <Form {...form}>
        <form
          onSubmit={form.handleSubmit((values) => onSubmit(values.code))}
          className="space-y-4"
        >
          <FormField
            control={form.control}
            name="code"
            render={({ field }) => (
              <FormItem>
                <FormLabel>验证码</FormLabel>
                <FormControl>
                  <Input
                    type="text"
                    inputMode="numeric"
                    autoComplete="one-time-code"
                    placeholder="请输入 6 位验证码"
                    maxLength={6}
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button type="submit" className="w-full" disabled={isPending}>
            {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            确认绑定
          </Button>
          {onBack && (
            <Button type="button" variant="ghost" className="w-full" onClick={onBack}>
              返回
            </Button>
          )}
        </form>
      </Form>
    </div>
  );
}
