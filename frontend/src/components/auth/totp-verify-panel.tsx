"use client";

import { Loader2 } from "lucide-react";
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

interface TotpVerifyPanelProps {
  title?: string;
  description?: string;
  isPending?: boolean;
  errorMessage?: string | null;
  onSubmit: (code: string) => void;
  onBack?: () => void;
}

export function TotpVerifyPanel({
  title = "两步验证",
  description = "请输入验证器 App 中的 6 位动态验证码。",
  isPending,
  errorMessage,
  onSubmit,
  onBack,
}: TotpVerifyPanelProps) {
  const form = useForm<TotpCodeFormValues>({
    resolver: zodResolver(totpCodeSchema),
    defaultValues: { code: "" },
  });

  return (
    <div className="space-y-4">
      <div className="text-center">
        <h2 className="text-xl font-semibold">{title}</h2>
        <p className="mt-2 text-sm text-slate-500">{description}</p>
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
                    autoFocus
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button type="submit" className="w-full" disabled={isPending}>
            {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            验证并登录
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
