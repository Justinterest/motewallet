"use client";

import { useState } from "react";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2, Eye, EyeOff } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Card,
  CardContent,
  CardFooter,
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

import { loginSchema, type LoginFormValues } from "@/lib/validations/auth";
import {
  useConfirm2FASetup,
  useLogin,
  useVerify2FA,
} from "@/lib/hooks/use-auth";
import type { AuthChallenge } from "@/types/auth";

type Step = "credentials" | "verify" | "setup";

export default function LoginPage() {
  const [showPassword, setShowPassword] = useState(false);
  const [step, setStep] = useState<Step>("credentials");
  const [challenge, setChallenge] = useState<AuthChallenge | null>(null);

  const loginMutation = useLogin();
  const verifyMutation = useVerify2FA();
  const setupMutation = useConfirm2FASetup();

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  function handleChallenge(next: AuthChallenge) {
    setChallenge(next);
    if (next.status === "REQUIRES_2FA") {
      setStep("verify");
      return;
    }
    if (next.status === "REQUIRES_2FA_SETUP") {
      setStep("setup");
    }
  }

  function onSubmit(values: LoginFormValues) {
    loginMutation.mutate(values, {
      onSuccess: handleChallenge,
    });
  }

  function handleBack() {
    setStep("credentials");
    setChallenge(null);
    verifyMutation.reset();
    setupMutation.reset();
  }

  const verifyError =
    verifyMutation.error?.message ||
    (verifyMutation.isError ? "验证失败，请重试" : null);
  const setupError =
    setupMutation.error?.message ||
    (setupMutation.isError ? "绑定失败，请重试" : null);

  return (
    <Card>
      {step === "credentials" && (
        <>
          <CardHeader className="pb-4">
            <CardTitle className="text-center text-xl">登录</CardTitle>
          </CardHeader>
          <CardContent>
            {loginMutation.isError && (
              <div className="mb-4 rounded-md bg-red-50 p-3 text-sm text-red-600">
                {loginMutation.error?.message || "登录失败，请重试"}
              </div>
            )}
            <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>邮箱</FormLabel>
                      <FormControl>
                        <Input
                          type="email"
                          placeholder="请输入邮箱地址"
                          autoComplete="email"
                          {...field}
                        />
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
                        <div className="relative">
                          <Input
                            type={showPassword ? "text" : "password"}
                            placeholder="请输入密码"
                            autoComplete="current-password"
                            {...field}
                          />
                          <button
                            type="button"
                            className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                            onClick={() => setShowPassword(!showPassword)}
                            tabIndex={-1}
                          >
                            {showPassword ? (
                              <EyeOff className="h-4 w-4" />
                            ) : (
                              <Eye className="h-4 w-4" />
                            )}
                          </button>
                        </div>
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <div className="flex items-center space-x-2">
                  <Checkbox id="remember" />
                  <label
                    htmlFor="remember"
                    className="text-sm leading-none text-slate-600"
                  >
                    记住我
                  </label>
                </div>
                <Button
                  type="submit"
                  className="w-full"
                  disabled={loginMutation.isPending}
                >
                  {loginMutation.isPending && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  登录
                </Button>
              </form>
            </Form>
          </CardContent>
          <CardFooter className="justify-center">
            <p className="text-sm text-slate-500">
              还没有账号？
              <Link
                href="/register"
                className="font-medium text-primary hover:underline"
              >
                立即注册
              </Link>
            </p>
          </CardFooter>
        </>
      )}

      {step === "verify" && challenge?.temp_token && (
        <CardContent className="pt-6">
          <TotpVerifyPanel
            isPending={verifyMutation.isPending}
            errorMessage={verifyError}
            onBack={handleBack}
            onSubmit={(code) =>
              verifyMutation.mutate({
                temp_token: challenge.temp_token!,
                code,
              })
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
              description="管理员已重置或尚未绑定两步验证。请扫描二维码并输入验证码完成绑定后登录。"
              isPending={setupMutation.isPending}
              errorMessage={setupError}
              onBack={handleBack}
              onSubmit={(code) =>
                setupMutation.mutate({
                  temp_token: challenge.temp_token!,
                  code,
                })
              }
            />
          </CardContent>
        )}
    </Card>
  );
}
