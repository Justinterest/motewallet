"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2, ShieldCheck } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { TotpSetupPanel } from "@/components/auth/totp-setup-panel";
import {
  useConfirmTotpRebind,
  usePrepareTotpRebind,
  useTotpStatus,
} from "@/lib/hooks/use-auth";
import { totpCodeSchema, type TotpCodeFormValues } from "@/lib/validations/auth";
import type { TotpSetup } from "@/types/auth";

type Step = "status" | "verify-current" | "setup-new";

export default function SettingsPage() {
  const [step, setStep] = useState<Step>("status");
  const [setup, setSetup] = useState<TotpSetup | null>(null);

  const { data: status, isLoading } = useTotpStatus();
  const prepareMutation = usePrepareTotpRebind();
  const confirmMutation = useConfirmTotpRebind();

  const currentForm = useForm<TotpCodeFormValues>({
    resolver: zodResolver(totpCodeSchema),
    defaultValues: { code: "" },
  });

  function handlePrepare(values: TotpCodeFormValues) {
    prepareMutation.mutate(values.code, {
      onSuccess: (data) => {
        setSetup(data);
        setStep("setup-new");
      },
    });
  }

  function handleConfirm(code: string) {
    confirmMutation.mutate(code, {
      onSuccess: () => {
        setSetup(null);
        setStep("status");
        currentForm.reset();
      },
    });
  }

  function handleCancel() {
    setStep("status");
    setSetup(null);
    prepareMutation.reset();
    confirmMutation.reset();
    currentForm.reset();
  }

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">设置</h1>
        <p className="mt-1 text-sm text-slate-500">管理账户安全与两步验证</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <ShieldCheck className="h-5 w-5" />
            两步验证
          </CardTitle>
          <CardDescription>
            使用验证器 App 保护账户安全。重新绑定后，旧的验证器条目将立即失效。
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="flex items-center gap-2 text-sm text-slate-500">
              <Loader2 className="h-4 w-4 animate-spin" />
              加载中…
            </div>
          ) : step === "status" ? (
            <div className="flex items-center justify-between gap-4">
              <div>
                <div className="flex items-center gap-2">
                  <span className="font-medium">当前状态</span>
                  {status?.enabled ? (
                    <Badge className="bg-green-100 text-green-700 border-green-200">
                      已开启
                    </Badge>
                  ) : (
                    <Badge variant="outline">未开启</Badge>
                  )}
                </div>
                <p className="mt-1 text-sm text-slate-500">
                  {status?.enabled
                    ? "下次登录需输入验证器中的动态验证码"
                    : "登录时将强制要求绑定两步验证"}
                </p>
              </div>
              {status?.enabled && (
                <Button onClick={() => setStep("verify-current")}>
                  重新绑定
                </Button>
              )}
            </div>
          ) : step === "verify-current" ? (
            <div className="space-y-4">
              <p className="text-sm text-slate-500">
                请先输入当前验证器中的 6 位验证码，以确认身份。
              </p>
              {prepareMutation.isError && (
                <div className="rounded-md bg-red-50 p-3 text-sm text-red-600">
                  {prepareMutation.error?.message || "验证失败，请重试"}
                </div>
              )}
              <Form {...currentForm}>
                <form
                  onSubmit={currentForm.handleSubmit(handlePrepare)}
                  className="space-y-4"
                >
                  <FormField
                    control={currentForm.control}
                    name="code"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>当前验证码</FormLabel>
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
                  <div className="flex gap-2">
                    <Button
                      type="submit"
                      disabled={prepareMutation.isPending}
                    >
                      {prepareMutation.isPending && (
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      )}
                      下一步
                    </Button>
                    <Button type="button" variant="ghost" onClick={handleCancel}>
                      取消
                    </Button>
                  </div>
                </form>
              </Form>
            </div>
          ) : (
            setup && (
              <TotpSetupPanel
                totpUri={setup.totp_uri}
                totpSecret={setup.totp_secret}
                title="绑定新的验证器"
                description="请用新的验证器条目扫描二维码，并输入新验证码完成重新绑定。"
                isPending={confirmMutation.isPending}
                errorMessage={
                  confirmMutation.error?.message ||
                  (confirmMutation.isError ? "绑定失败，请重试" : null)
                }
                onBack={handleCancel}
                onSubmit={handleConfirm}
              />
            )
          )}
        </CardContent>
      </Card>
    </div>
  );
}
