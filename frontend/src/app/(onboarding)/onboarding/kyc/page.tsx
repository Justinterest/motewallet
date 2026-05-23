"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2 } from "lucide-react";

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
import { NativeSelect } from "@/components/ui/select";

import { kycFormSchema, type KycFormValues } from "@/lib/validations/onboarding";
import { useSubmitKyc } from "@/lib/hooks/use-onboarding";
import { toast } from "@/hooks/use-toast";

const countries = [
  { value: "CN", label: "中国大陆" },
  { value: "HK", label: "中国香港" },
  { value: "SG", label: "新加坡" },
  { value: "US", label: "美国" },
  { value: "GB", label: "英国" },
  { value: "JP", label: "日本" },
];

export default function KycPage() {
  const router = useRouter();
  const submitKycMutation = useSubmitKyc();

  const form = useForm<KycFormValues>({
    resolver: zodResolver(kycFormSchema),
    defaultValues: {
      company_name: "",
      country: "",
      registration_number: "",
      contact_name: "",
      contact_phone: "",
    },
  });

  function onSubmit(values: KycFormValues) {
    submitKycMutation.mutate(values, {
      onSuccess: () => {
        toast({
          title: "提交成功",
          description: "实名认证资料已提交，请等待审核。",
        });
        router.push("/onboarding/status");
      },
      onError: (error) => {
        toast({
          variant: "destructive",
          title: "提交失败",
          description: error.message || "请检查信息后重试",
        });
      },
    });
  }

  return (
    <div className="space-y-4">
      <div className="text-center">
        <h2 className="text-lg font-semibold text-slate-900">实名认证</h2>
        <p className="mt-1 text-sm text-slate-500">
          请填写企业信息完成实名认证
        </p>
      </div>

      <Card>
        <CardHeader className="pb-4">
          <CardTitle className="text-base">企业信息</CardTitle>
        </CardHeader>
        <CardContent>
          {submitKycMutation.isError && (
            <div className="mb-4 rounded-md bg-red-50 p-3 text-sm text-red-600">
              {submitKycMutation.error?.message || "提交失败，请重试"}
            </div>
          )}

          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="company_name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>公司名称</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="请输入公司全称"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="country"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>国家/地区</FormLabel>
                    <FormControl>
                      <NativeSelect
                        value={field.value}
                        onChange={field.onChange}
                        onBlur={field.onBlur}
                      >
                        <option value="" disabled>
                          请选择国家/地区
                        </option>
                        {countries.map((c) => (
                          <option key={c.value} value={c.value}>
                            {c.label}
                          </option>
                        ))}
                      </NativeSelect>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="registration_number"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>注册编号</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="统一社会信用代码 / 注册编号"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="contact_name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>联系人姓名</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="请输入联系人姓名"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="contact_phone"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>联系电话</FormLabel>
                    <FormControl>
                      <Input
                        type="tel"
                        placeholder="请输入联系电话"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button
                type="submit"
                className="w-full bg-blue-700 hover:bg-blue-800 text-white"
                disabled={submitKycMutation.isPending}
                size="lg"
              >
                {submitKycMutation.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                提交认证
              </Button>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
