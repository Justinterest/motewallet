"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useFieldArray, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ChevronLeft, ChevronRight, Loader2, Plus, Trash2 } from "lucide-react";

import { Button } from "@/components/ui/button";
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
import { Input } from "@/components/ui/input";
import { FormSelect } from "@/components/ui/form-select";
import { FilePathList } from "@/components/kyc/file-path-list";
import { PersonFields } from "@/components/kyc/person-fields";
import { cn } from "@/lib/utils";
import { KYC_DRAFT_KEY } from "@/lib/kyc/constants";
import { createKycFormDefaults } from "@/lib/kyc/form-defaults";
import { formValuesToSubmitRequest } from "@/lib/kyc/transform";
import {
  ENTERPRISE_TYPES,
  FUNDING_SOURCES,
  GENDERS,
  OPEN_ACCOUNT_PURPOSES,
  REGISTER_REGIONS,
  VERIFICATION_TYPES,
} from "@/lib/kyc/constants";
import {
  kycFormSchema,
  kycStep0Schema,
  kycStep1Schema,
  kycStep2Schema,
  type KycFormValues,
} from "@/lib/validations/onboarding";
import { useSubmitKyc } from "@/lib/hooks/use-onboarding";
import { toast } from "@/hooks/use-toast";

const STEPS = ["企业信息", "管理人信息", "股东与董事"] as const;

export function KycWizard() {
  const router = useRouter();
  const [step, setStep] = useState(0);
  const submitMutation = useSubmitKyc();

  const form = useForm<KycFormValues>({
    resolver: zodResolver(kycFormSchema),
    defaultValues: createKycFormDefaults(),
    mode: "onBlur",
  });

  const shareholders = useFieldArray({
    control: form.control,
    name: "shareholdersInfo",
  });
  const directors = useFieldArray({
    control: form.control,
    name: "directorInfo",
  });

  const registerRegion = form.watch("enterpriseInfo.registerRegion");
  const managerVerification = form.watch(
    "enterpriseInfo.managerVerificationType"
  );

  useEffect(() => {
    try {
      const raw = localStorage.getItem(KYC_DRAFT_KEY);
      if (raw) {
        const parsed = JSON.parse(raw) as KycFormValues;
        form.reset(parsed);
      }
    } catch {
      /* ignore corrupt draft */
    }
  }, [form]);

  useEffect(() => {
    const sub = form.watch((values) => {
      try {
        localStorage.setItem(KYC_DRAFT_KEY, JSON.stringify(values));
      } catch {
        /* quota exceeded */
      }
    });
    return () => sub.unsubscribe();
  }, [form]);

  async function validateCurrentStep(): Promise<boolean> {
    const values = form.getValues();
    const schemas = [kycStep0Schema, kycStep1Schema, kycStep2Schema] as const;
    const result = schemas[step].safeParse(values);
    if (!result.success) {
      result.error.issues.forEach((issue) => {
        const path = issue.path.join(".") as Parameters<
          typeof form.setError
        >[0];
        form.setError(path, { message: issue.message });
      });
      return false;
    }
    return true;
  }

  async function handleNext() {
    const ok = await validateCurrentStep();
    if (!ok) return;
    setStep((s) => Math.min(s + 1, STEPS.length - 1));
  }

  function handleBack() {
    setStep((s) => Math.max(s - 1, 0));
  }

  async function onSubmit(values: KycFormValues) {
    const payload = formValuesToSubmitRequest(values);
    submitMutation.mutate(payload, {
      onSuccess: () => {
        localStorage.removeItem(KYC_DRAFT_KEY);
        toast({
          title: "提交成功",
          description: "实名认证资料已提交，请等待审核。",
        });
        router.push("/kyc/status");
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
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <div className="flex items-center justify-center gap-2">
          {STEPS.map((label, index) => (
            <div key={label} className="flex items-center">
              <div
                className={cn(
                  "flex h-8 w-8 items-center justify-center rounded-full text-sm font-medium",
                  index <= step
                    ? "bg-blue-700 text-white"
                    : "bg-slate-200 text-slate-500"
                )}
              >
                {index + 1}
              </div>
              <span
                className={cn(
                  "ml-2 hidden text-sm sm:inline",
                  index === step ? "font-medium text-blue-700" : "text-slate-500"
                )}
              >
                {label}
              </span>
              {index < STEPS.length - 1 && (
                <div className="mx-3 h-px w-8 bg-slate-200 sm:w-12" />
              )}
            </div>
          ))}
        </div>

        {submitMutation.isError && (
          <div className="rounded-md bg-red-50 p-3 text-sm text-red-600">
            {submitMutation.error?.message || "提交失败，请重试"}
          </div>
        )}

        {step === 0 && (
          <Card>
            <CardHeader>
              <CardTitle className="text-base">企业信息</CardTitle>
              <CardDescription>
                字段与鲲「子商户入网认证」接口一致；文件须先通过鲲上传接口获取 path。
              </CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4 sm:grid-cols-2">
              <FormField
                control={form.control}
                name="enterpriseInfo.enterpriseEN"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <FormLabel>公司英文名称 *</FormLabel>
                    <FormControl>
                      <Input placeholder="须与注册证书一致" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.enterpriseNameCHS"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>公司中文名称</FormLabel>
                    <FormControl>
                      <Input placeholder="无中文名填「无」" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.registerRegion"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>注册地区 *</FormLabel>
                    <FormControl>
                      <FormSelect {...field} options={REGISTER_REGIONS} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.incorporationCertificateNo"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>注册证书编号 *</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.establishTime"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>成立日期 *</FormLabel>
                    <FormControl>
                      <Input type="date" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.registerAddress"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <FormLabel>注册地址 *</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.enterpriseType"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>公司类型 *</FormLabel>
                    <FormControl>
                      <FormSelect {...field} options={ENTERPRISE_TYPES} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.phone"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>企业电话</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.industry"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>一级行业 code *</FormLabel>
                    <FormControl>
                      <Input placeholder="鲲枚举 industry" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.subIndustry"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>二级行业 code *</FormLabel>
                    <FormControl>
                      <Input placeholder="鲲枚举 subIndustry" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.initialFundingSource"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>初始资金来源 *</FormLabel>
                    <FormControl>
                      <FormSelect {...field} options={FUNDING_SOURCES} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.openAccountPurpose"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>开户目的 *</FormLabel>
                    <FormControl>
                      <FormSelect {...field} options={OPEN_ACCOUNT_PURPOSES} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.incorporationCertificatePaths"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <FormControl>
                      <FilePathList
                        label="营业执照 / 注册证书 *"
                        paths={field.value}
                        onChange={field.onChange}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              {registerRegion === "HK" && (
                <>
                  <FormField
                    control={form.control}
                    name="enterpriseInfo.businessRegistrationNo"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>商业登记证编号 *</FormLabel>
                        <FormControl>
                          <Input {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name="enterpriseInfo.businessRegistrationPaths"
                    render={({ field }) => (
                      <FormItem className="sm:col-span-2">
                        <FormControl>
                          <FilePathList
                            label="商业登记证 *"
                            paths={field.value ?? [""]}
                            onChange={field.onChange}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </>
              )}
            </CardContent>
          </Card>
        )}

        {step === 1 && (
          <Card>
            <CardHeader>
              <CardTitle className="text-base">账户管理人</CardTitle>
              <CardDescription>
                管理人信息包含在 enterpriseInfo 中提交给鲲。
              </CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4 sm:grid-cols-2">
              <FormField
                control={form.control}
                name="enterpriseInfo.managerNameEN"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>管理人英文名 *</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.managerSurnameCHS"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>中文姓</FormLabel>
                    <FormControl>
                      <Input placeholder="仅填姓氏" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.managerCountry"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>国籍 *</FormLabel>
                    <FormControl>
                      <FormSelect {...field} options={REGISTER_REGIONS} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.managerAuthType"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>证件类型 code *</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.managerIdCard"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>证件号码 *</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.managerBirthday"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>出生日期 *</FormLabel>
                    <FormControl>
                      <Input type="date" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.managerGender"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>性别 *</FormLabel>
                    <FormControl>
                      <FormSelect {...field} options={GENDERS} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.managerVerificationType"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>验证方式 *</FormLabel>
                    <FormControl>
                      <FormSelect {...field} options={VERIFICATION_TYPES} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.managerResidenceCountry"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>居住国家 *</FormLabel>
                    <FormControl>
                      <FormSelect {...field} options={REGISTER_REGIONS} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.managerResidenceAddress"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <FormLabel>居住地址 *</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.managerContactsEmail"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <FormLabel>联系邮箱</FormLabel>
                    <FormControl>
                      <Input type="email" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              {managerVerification === "idHolding" && (
                <FormField
                  control={form.control}
                  name="enterpriseInfo.managerIdHoldingPaths"
                  render={({ field }) => (
                    <FormItem className="sm:col-span-2">
                      <FormControl>
                        <FilePathList
                          label="管理人证件照 *"
                          description="至少 3 张：证件照、自拍、手持证件照"
                          paths={field.value ?? ["", "", ""]}
                          onChange={field.onChange}
                          minItems={3}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              )}
            </CardContent>
          </Card>
        )}

        {step === 2 && (
          <div className="space-y-6">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <div>
                  <CardTitle className="text-base">股东信息</CardTitle>
                  <CardDescription>至少一名股东</CardDescription>
                </div>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() =>
                    shareholders.append({
                      ...createKycFormDefaults().shareholdersInfo[0],
                    })
                  }
                >
                  <Plus className="mr-1 h-4 w-4" />
                  添加股东
                </Button>
              </CardHeader>
              <CardContent className="space-y-8">
                {shareholders.fields.map((field, index) => (
                  <div
                    key={field.id}
                    className="rounded-lg border border-slate-200 p-4"
                  >
                    <div className="mb-4 flex items-center justify-between">
                      <h4 className="font-medium text-slate-800">
                        股东 {index + 1}
                      </h4>
                      {shareholders.fields.length > 1 && (
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => shareholders.remove(index)}
                        >
                          <Trash2 className="h-4 w-4 text-red-500" />
                        </Button>
                      )}
                    </div>
                    <PersonFields
                      control={form.control}
                      namePrefix={`shareholdersInfo.${index}`}
                      showShareholding
                    />
                  </div>
                ))}
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <div>
                  <CardTitle className="text-base">董事信息</CardTitle>
                  <CardDescription>至少一名董事</CardDescription>
                </div>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() =>
                    directors.append({
                      ...createKycFormDefaults().directorInfo[0],
                    })
                  }
                >
                  <Plus className="mr-1 h-4 w-4" />
                  添加董事
                </Button>
              </CardHeader>
              <CardContent className="space-y-8">
                {directors.fields.map((field, index) => (
                  <div
                    key={field.id}
                    className="rounded-lg border border-slate-200 p-4"
                  >
                    <div className="mb-4 flex items-center justify-between">
                      <h4 className="font-medium text-slate-800">
                        董事 {index + 1}
                      </h4>
                      {directors.fields.length > 1 && (
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => directors.remove(index)}
                        >
                          <Trash2 className="h-4 w-4 text-red-500" />
                        </Button>
                      )}
                    </div>
                    <PersonFields
                      control={form.control}
                      namePrefix={`directorInfo.${index}`}
                    />
                  </div>
                ))}
              </CardContent>
            </Card>
          </div>
        )}

        <div className="flex justify-between gap-3">
          <Button
            type="button"
            variant="outline"
            onClick={handleBack}
            disabled={step === 0}
          >
            <ChevronLeft className="mr-1 h-4 w-4" />
            上一步
          </Button>
          {step < STEPS.length - 1 ? (
            <Button
              type="button"
              className="bg-blue-700 hover:bg-blue-800 text-white"
              onClick={handleNext}
            >
              下一步
              <ChevronRight className="ml-1 h-4 w-4" />
            </Button>
          ) : (
            <Button
              type="submit"
              className="bg-blue-700 hover:bg-blue-800 text-white"
              disabled={submitMutation.isPending}
            >
              {submitMutation.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              提交认证
            </Button>
          )}
        </div>
      </form>
    </Form>
  );
}
