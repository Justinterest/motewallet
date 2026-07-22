"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useQueryClient } from "@tanstack/react-query";
import {
  completeAdminAuth,
  useAdminLogin,
  useChangeAdminPassword,
  useConfirmAdmin2FASetup,
  useVerifyAdmin2FA,
} from "@/lib/hooks/use-auth";
import { useAuthStore } from "@/stores/auth-store";
import {
  adminLoginSchema,
  changePasswordSchema,
  type AdminLoginFormValues,
  type ChangePasswordFormValues,
} from "@/lib/validations/auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
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
import { TotpVerifyPanel } from "@/components/auth/totp-verify-panel";
import type { AdminAuthChallenge } from "@/types/auth";

type Step = "credentials" | "verify" | "setup" | "change-password";

export default function LoginPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { setAdmin } = useAuthStore();
  const [step, setStep] = useState<Step>("credentials");
  const [challenge, setChallenge] = useState<AdminAuthChallenge | null>(null);

  const loginMutation = useAdminLogin();
  const verifyMutation = useVerifyAdmin2FA();
  const setupMutation = useConfirmAdmin2FASetup();
  const changePasswordMutation = useChangeAdminPassword();

  const form = useForm<AdminLoginFormValues>({
    resolver: zodResolver(adminLoginSchema),
    defaultValues: {
      username: "",
      password: "",
    },
  });

  const passwordForm = useForm<ChangePasswordFormValues>({
    resolver: zodResolver(changePasswordSchema),
    defaultValues: {
      newPassword: "",
      confirmPassword: "",
    },
  });

  function handleChallenge(next: AdminAuthChallenge) {
    setChallenge(next);
    if (next.status === "REQUIRES_2FA") {
      setStep("verify");
      return;
    }
    if (next.status === "REQUIRES_2FA_SETUP") {
      setStep("setup");
      return;
    }
    if (next.status === "REQUIRES_PASSWORD_CHANGE") {
      setStep("change-password");
      passwordForm.reset();
      return;
    }
    if (next.status === "SUCCESS") {
      completeAdminAuth(next, setAdmin, queryClient, router);
    }
  }

  function onSubmit(values: AdminLoginFormValues) {
    loginMutation.mutate(values, {
      onSuccess: handleChallenge,
    });
  }

  function handleBack() {
    setStep("credentials");
    setChallenge(null);
    verifyMutation.reset();
    setupMutation.reset();
    changePasswordMutation.reset();
    passwordForm.reset();
  }

  const verifyError =
    verifyMutation.error?.message ||
    (verifyMutation.isError ? "验证失败，请重试" : null);
  const setupError =
    setupMutation.error?.message ||
    (setupMutation.isError ? "绑定失败，请重试" : null);
  const changePasswordError =
    changePasswordMutation.error?.message ||
    (changePasswordMutation.isError ? "修改密码失败，请重试" : null);

  return (
    <Card>
      {step === "credentials" && (
        <>
          <CardHeader>
            <CardTitle className="text-center text-xl">管理员登录</CardTitle>
          </CardHeader>
          <CardContent>
            {loginMutation.isError && (
              <div className="mb-4 rounded-md bg-destructive/10 px-4 py-3 text-sm text-destructive">
                {loginMutation.error?.message || "登录失败，请重试"}
              </div>
            )}
            <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                <FormField
                  control={form.control}
                  name="username"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>用户名</FormLabel>
                      <FormControl>
                        <Input placeholder="请输入用户名" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>密码</FormLabel>
                      <FormControl>
                        <Input
                          type="password"
                          placeholder="请输入密码"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <Button
                  type="submit"
                  className="w-full"
                  disabled={loginMutation.isPending}
                >
                  {loginMutation.isPending ? "登录中..." : "登录"}
                </Button>
              </form>
            </Form>
          </CardContent>
        </>
      )}

      {step === "verify" && challenge?.temp_token && (
        <CardContent className="pt-6">
          <TotpVerifyPanel
            isPending={verifyMutation.isPending}
            errorMessage={verifyError}
            onBack={handleBack}
            onSubmit={(code) =>
              verifyMutation.mutate(
                { temp_token: challenge.temp_token!, code },
                { onSuccess: handleChallenge }
              )
            }
          />
        </CardContent>
      )}

      {step === "setup" &&
        challenge?.temp_token &&
        challenge.totp_uri &&
        challenge.totp_secret && (
          <CardContent className="pt-6">
            <TotpSetupPanel
              totpUri={challenge.totp_uri}
              totpSecret={challenge.totp_secret}
              title="开启两步验证"
              description="首次登录或管理员已重置两步验证。请扫描二维码并输入验证码完成绑定后登录。"
              isPending={setupMutation.isPending}
              errorMessage={setupError}
              onBack={handleBack}
              onSubmit={(code) =>
                setupMutation.mutate(
                  { temp_token: challenge.temp_token!, code },
                  { onSuccess: handleChallenge }
                )
              }
            />
          </CardContent>
        )}

      {step === "change-password" && challenge?.temp_token && (
        <CardContent className="pt-6">
          <div className="mb-4 text-center">
            <h2 className="text-xl font-semibold">修改初始密码</h2>
            <p className="mt-2 text-sm text-slate-500">
              首次登录或密码被重置后，必须立即设置新密码才能继续使用。
            </p>
          </div>
          {changePasswordError && (
            <div className="mb-4 rounded-md bg-destructive/10 px-4 py-3 text-sm text-destructive">
              {changePasswordError}
            </div>
          )}
          <Form {...passwordForm}>
            <form
              onSubmit={passwordForm.handleSubmit((values) =>
                changePasswordMutation.mutate(
                  {
                    temp_token: challenge.temp_token!,
                    new_password: values.newPassword,
                  },
                  { onSuccess: handleChallenge }
                )
              )}
              className="space-y-4"
            >
              <FormField
                control={passwordForm.control}
                name="newPassword"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>新密码</FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        placeholder="至少 8 位"
                        autoFocus
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={passwordForm.control}
                name="confirmPassword"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>确认新密码</FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        placeholder="再次输入新密码"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button
                type="submit"
                className="w-full"
                disabled={changePasswordMutation.isPending}
              >
                {changePasswordMutation.isPending
                  ? "提交中..."
                  : "确认修改并登录"}
              </Button>
              <Button
                type="button"
                variant="ghost"
                className="w-full"
                onClick={handleBack}
              >
                返回
              </Button>
            </form>
          </Form>
        </CardContent>
      )}
    </Card>
  );
}
