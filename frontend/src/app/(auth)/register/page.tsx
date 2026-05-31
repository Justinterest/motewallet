"use client";

import { useEffect, useState } from "react";
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

import {
  registerSchema,
  type RegisterFormValues,
} from "@/lib/validations/auth";
import { useRegister, useSendVerificationCode } from "@/lib/hooks/use-auth";

const RESEND_COOLDOWN_SECONDS = 60;

export default function RegisterPage() {
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [countdown, setCountdown] = useState(0);
  const [sendMessage, setSendMessage] = useState<string | null>(null);
  const registerMutation = useRegister();
  const sendCodeMutation = useSendVerificationCode();

  const form = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      email: "",
      verificationCode: "",
      password: "",
      confirmPassword: "",
      agreeTerms: false as unknown as true,
    },
  });

  useEffect(() => {
    if (countdown <= 0) {
      return;
    }

    const timer = window.setTimeout(() => {
      setCountdown((value) => value - 1);
    }, 1000);

    return () => window.clearTimeout(timer);
  }, [countdown]);

  async function handleSendCode() {
    setSendMessage(null);

    const emailValid = await form.trigger("email");
    if (!emailValid) {
      return;
    }

    const email = form.getValues("email");
    sendCodeMutation.mutate(email, {
      onSuccess: () => {
        setSendMessage("验证码已发送，请查收邮件");
        setCountdown(RESEND_COOLDOWN_SECONDS);
      },
      onError: (error) => {
        setSendMessage(error.message || "验证码发送失败，请稍后重试");
      },
    });
  }

  function onSubmit(values: RegisterFormValues) {
    registerMutation.mutate({
      email: values.email,
      password: values.password,
      verification_code: values.verificationCode,
    });
  }

  const canSendCode =
    countdown === 0 && !sendCodeMutation.isPending;

  return (
    <Card>
      <CardHeader className="pb-4">
        <CardTitle className="text-center text-xl">创建账号</CardTitle>
      </CardHeader>
      <CardContent>
        {registerMutation.isError && (
          <div className="mb-4 rounded-md bg-red-50 p-3 text-sm text-red-600">
            {registerMutation.error?.message || "注册失败，请重试"}
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
              name="verificationCode"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>邮箱验证码</FormLabel>
                  <div className="flex gap-2">
                    <FormControl>
                      <Input
                        type="text"
                        inputMode="numeric"
                        autoComplete="one-time-code"
                        placeholder="请输入 6 位验证码"
                        maxLength={6}
                        className="flex-1"
                        {...field}
                      />
                    </FormControl>
                    <Button
                      type="button"
                      variant="outline"
                      className="shrink-0"
                      disabled={!canSendCode}
                      onClick={handleSendCode}
                    >
                      {sendCodeMutation.isPending ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : countdown > 0 ? (
                        `${countdown}s`
                      ) : (
                        "获取验证码"
                      )}
                    </Button>
                  </div>
                  {sendMessage && (
                    <p
                      className={
                        sendCodeMutation.isError
                          ? "text-sm text-red-600"
                          : "text-sm text-slate-500"
                      }
                    >
                      {sendMessage}
                    </p>
                  )}
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
                        placeholder="至少 8 位字符"
                        autoComplete="new-password"
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
            <FormField
              control={form.control}
              name="confirmPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>确认密码</FormLabel>
                  <FormControl>
                    <div className="relative">
                      <Input
                        type={showConfirmPassword ? "text" : "password"}
                        placeholder="请再次输入密码"
                        autoComplete="new-password"
                        {...field}
                      />
                      <button
                        type="button"
                        className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                        onClick={() =>
                          setShowConfirmPassword(!showConfirmPassword)
                        }
                        tabIndex={-1}
                      >
                        {showConfirmPassword ? (
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
            <FormField
              control={form.control}
              name="agreeTerms"
              render={({ field }) => (
                <FormItem>
                  <div className="flex items-center space-x-2">
                    <FormControl>
                      <Checkbox
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                    <FormLabel className="text-sm font-normal leading-none text-slate-600">
                      我同意服务条款
                    </FormLabel>
                  </div>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button
              type="submit"
              className="w-full"
              disabled={registerMutation.isPending}
            >
              {registerMutation.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              注册
            </Button>
          </form>
        </Form>
      </CardContent>
      <CardFooter className="justify-center">
        <p className="text-sm text-slate-500">
          已有账号？
          <Link href="/login" className="font-medium text-primary hover:underline">
            立即登录
          </Link>
        </p>
      </CardFooter>
    </Card>
  );
}
