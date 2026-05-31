"use client";

import { useEffect, useRef, useState } from "react";
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
  FormMessage,
} from "@/components/ui/form";
import { FormDatePicker } from "@/components/ui/date-picker";
import { FormSelect } from "@/components/ui/form-select";
import { Input } from "@/components/ui/input";
import { KycAuthTypeSelect } from "@/components/kyc/kyc-auth-type-select";
import { KycCountrySelect } from "@/components/kyc/kyc-country-select";
import { KycReferenceProvider } from "@/components/kyc/kyc-reference-provider";
import { KycFileUpload } from "@/components/kyc/kyc-file-upload";
import { KycFormLabel, kycPlaceholder } from "@/components/kyc/kyc-form-label";
import { CertificateValidityFields } from "@/components/kyc/certificate-validity-fields";
import { CollapsibleSection } from "@/components/kyc/collapsible-section";
import { KycFormRow } from "@/components/kyc/kyc-form-row";
import { PersonFields } from "@/components/kyc/person-fields";
import { cn } from "@/lib/utils";
import { getKycFieldMeta } from "@/lib/kyc/field-meta";
import { KYC_DRAFT_KEY } from "@/lib/kyc/constants";
import { createKycFormDefaults } from "@/lib/kyc/form-defaults";
import { formValuesToSubmitRequest } from "@/lib/kyc/transform";
import {
  EMPLOYEE_NUM_OPTIONS,
  ENTERPRISE_TYPES,
  FUNDING_SOURCES,
  GENDERS,
  INDUSTRIES,
  OPEN_ACCOUNT_PURPOSES,
  SALES_VOLUME_OPTIONS,
  WEALTH_SOURCES,
  YES_NO,
} from "@/lib/kyc/constants";
import { kycFormSchema, type KycFormValues } from "@/lib/validations/onboarding";
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
  const nameChanged = form.watch(
    "enterpriseInfo.isChangeEnterpriseNameInFiveYears"
  );
  const middleTierShareholders = form.watch(
    "enterpriseInfo.middleTierShareholders"
  );
  const managerCountry = form.watch("enterpriseInfo.managerCountry");
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

  useEffect(() => {
    if (middleTierShareholders === "No") {
      form.setValue("enterpriseInfo.equityStructurePaths", []);
    }
  }, [middleTierShareholders, form]);

  const prevManagerCountryRef = useRef(managerCountry);
  useEffect(() => {
    if (
      prevManagerCountryRef.current !== undefined &&
      prevManagerCountryRef.current !== managerCountry
    ) {
      form.setValue("enterpriseInfo.managerAuthType", "");
    }
    prevManagerCountryRef.current = managerCountry;
  }, [managerCountry, form]);

  function goToStep(index: number) {
    setStep(Math.max(0, Math.min(index, STEPS.length - 1)));
  }

  function handleNext() {
    setStep((s) => Math.min(s + 1, STEPS.length - 1));
  }

  function handleBack() {
    setStep((s) => Math.max(s - 1, 0));
  }

  function stepIndexForIssuePath(path: PropertyKey[]): number {
    const parts = path.map(String);
    if (parts[0] === "shareholdersInfo" || parts[0] === "directorInfo") {
      return 2;
    }
    if (parts[0] === "enterpriseInfo" && parts[1]) {
      const key = parts[1];
      if (
        key.startsWith("manager") ||
        key === "authorizationLetterPaths" ||
        key === "equityStructurePaths" ||
        key === "middleTierShareholders" ||
        key === "nnc1Paths"
      ) {
        return 1;
      }
      return 0;
    }
    return 0;
  }

  function focusFirstInvalidStep() {
    const result = kycFormSchema.safeParse(form.getValues());
    if (result.success) return;

    let firstStep = 2;
    result.error.issues.forEach((issue) => {
      firstStep = Math.min(firstStep, stepIndexForIssuePath(issue.path));
      const path = issue.path.join(".") as Parameters<typeof form.setError>[0];
      form.setError(path, { message: issue.message });
    });

    goToStep(firstStep);
    toast({
      variant: "destructive",
      title: "请完善表单",
      description: `请先完成「${STEPS[firstStep]}」中的必填项。`,
    });
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
    <KycReferenceProvider>
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit, focusFirstInvalidStep)}
        className="space-y-6"
      >
        <nav
          className="flex items-center justify-center gap-2"
          aria-label="认证步骤"
        >
          {STEPS.map((label, index) => (
            <div key={label} className="flex items-center">
              <button
                type="button"
                onClick={() => goToStep(index)}
                className="flex items-center rounded-lg px-1 py-1 transition-colors hover:bg-slate-100"
                aria-current={index === step ? "step" : undefined}
              >
                <span
                  className={cn(
                    "flex h-8 w-8 items-center justify-center rounded-full text-sm font-medium",
                    index === step
                      ? "bg-blue-700 text-white"
                      : "bg-slate-200 text-slate-600"
                  )}
                >
                  {index + 1}
                </span>
                <span
                  className={cn(
                    "ml-2 hidden text-sm sm:inline",
                    index === step
                      ? "font-medium text-blue-700"
                      : "text-slate-600"
                  )}
                >
                  {label}
                </span>
              </button>
              {index < STEPS.length - 1 && (
                <div className="mx-3 h-px w-8 bg-slate-200 sm:w-12" />
              )}
            </div>
          ))}
        </nav>

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
                请如实填写企业基本信息并上传证明材料，提交后将进入平台审核。
              </CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4 sm:grid-cols-2">
              <FormField
                control={form.control}
                name="enterpriseInfo.incorporationCertificatePaths"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <FormControl>
                      <KycFileUpload
                        label={getKycFieldMeta("incorporationCertificate").label + " *"}
                        description={
                          getKycFieldMeta("incorporationCertificate").description
                        }
                        exampleImage={getKycFieldMeta("incorporationCertificate").exampleImage}
                        paths={field.value}
                        onChange={field.onChange}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.incorporationCertificateNo"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <KycFormLabel fieldKey="incorporationCertificateNo" />
                    <FormControl>
                      <Input
                        placeholder={kycPlaceholder("incorporationCertificateNo")}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.enterpriseEN"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <KycFormLabel fieldKey="enterpriseEN" />
                    <FormControl>
                      <Input
                        placeholder={kycPlaceholder("enterpriseEN")}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <KycFormRow>
                <FormField
                  control={form.control}
                  name="enterpriseInfo.enterpriseNameCHS"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="enterpriseNameCHS" required={false} />
                      <FormControl>
                        <Input
                          placeholder={kycPlaceholder("enterpriseNameCHS")}
                          {...field}
                        />
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
                      <KycFormLabel fieldKey="registerRegion" />
                      <FormControl>
                        <KycCountrySelect {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </KycFormRow>
              <KycFormRow>
                <FormField
                  control={form.control}
                  name="enterpriseInfo.establishTime"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="establishTime" />
                      <FormControl>
                        <FormDatePicker {...field} />
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
                      <KycFormLabel fieldKey="enterpriseType" />
                      <FormControl>
                        <FormSelect {...field} options={ENTERPRISE_TYPES} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </KycFormRow>
              <FormField
                control={form.control}
                name="enterpriseInfo.registerAddress"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <KycFormLabel fieldKey="registerAddress" />
                    <FormControl>
                      <Input
                        placeholder={kycPlaceholder("registerAddress")}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.mainBusinessAddress"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <KycFormLabel fieldKey="mainBusinessAddress" />
                    <FormControl>
                      <Input
                        placeholder={kycPlaceholder("mainBusinessAddress")}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <KycFormRow>
                <FormField
                  control={form.control}
                  name="enterpriseInfo.phone"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="phone" />
                      <FormControl>
                        <Input {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="enterpriseInfo.isChangeEnterpriseNameInFiveYears"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="isChangeEnterpriseNameInFiveYears" />
                      <FormControl>
                        <FormSelect {...field} options={YES_NO} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </KycFormRow>
              {nameChanged === "Yes" && (
                <FormField
                  control={form.control}
                  name="enterpriseInfo.usedEnterpriseName"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="usedEnterpriseName" required />
                      <FormControl>
                        <Input
                          placeholder={kycPlaceholder("usedEnterpriseName")}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              )}
              <FormField
                control={form.control}
                name="enterpriseInfo.enterpriseDomain"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <KycFormLabel fieldKey="enterpriseDomain" required={false} />
                    <FormControl>
                      <Input
                        placeholder={kycPlaceholder("enterpriseDomain")}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <KycFormRow>
                <FormField
                  control={form.control}
                  name="enterpriseInfo.industry"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="industry" />
                      <FormControl>
                        <FormSelect {...field} options={INDUSTRIES} />
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
                      <KycFormLabel fieldKey="subIndustry" />
                      <FormControl>
                        <Input
                          placeholder={kycPlaceholder("subIndustry")}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </KycFormRow>
              <KycFormRow>
                <FormField
                  control={form.control}
                  name="enterpriseInfo.initialFundingSource"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="initialFundingSource" />
                      <FormControl>
                        <FormSelect {...field} options={FUNDING_SOURCES} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="enterpriseInfo.wealthSource"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="wealthSource" />
                      <FormControl>
                        <FormSelect {...field} options={WEALTH_SOURCES} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </KycFormRow>
              <KycFormRow>
                <FormField
                  control={form.control}
                  name="enterpriseInfo.continuousFundingSource"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="continuousFundingSource" />
                      <FormControl>
                        <FormSelect {...field} options={FUNDING_SOURCES} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="enterpriseInfo.salesVolumeLastyear"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="salesVolumeLastyear" />
                      <FormControl>
                        <FormSelect {...field} options={SALES_VOLUME_OPTIONS} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </KycFormRow>
              <KycFormRow>
                <FormField
                  control={form.control}
                  name="enterpriseInfo.employeeNum"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="employeeNum" required={false} />
                      <FormControl>
                        <FormSelect {...field} options={EMPLOYEE_NUM_OPTIONS} />
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
                      <KycFormLabel fieldKey="openAccountPurpose" />
                      <FormControl>
                        <FormSelect {...field} options={OPEN_ACCOUNT_PURPOSES} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </KycFormRow>
              <FormField
                control={form.control}
                name="enterpriseInfo.associationRulesPaths"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <FormControl>
                      <KycFileUpload
                        label={getKycFieldMeta("associationRules").label + " *"}
                        description={getKycFieldMeta("associationRules").description}
                        exampleImage={getKycFieldMeta("associationRules").exampleImage}
                        paths={field.value ?? [""]}
                        onChange={field.onChange}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.authenticMaterialsPaths"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <FormControl>
                      <KycFileUpload
                        label={getKycFieldMeta("authenticMaterials").label + " *"}
                        description={getKycFieldMeta("authenticMaterials").description}
                        paths={field.value ?? [""]}
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
                        <KycFormLabel fieldKey="businessRegistrationNo" required />
                        <FormControl>
                          <Input
                            placeholder={kycPlaceholder("businessRegistrationNo")}
                            {...field}
                          />
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
                          <KycFileUpload
                            label={getKycFieldMeta("businessRegistration").label + " *"}
                            description={getKycFieldMeta("businessRegistration").description}
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
                填写企业账户管理人的身份与联系方式，须与授权文件一致。
              </CardDescription>
            </CardHeader>
            <CardContent className="grid gap-4 sm:grid-cols-2">
              <KycFormRow>
                <FormField
                  control={form.control}
                  name="enterpriseInfo.managerSurnameCHS"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="managerSurnameCHS" />
                      <FormControl>
                        <Input
                          placeholder={kycPlaceholder("managerSurnameCHS")}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="enterpriseInfo.managerNameCHS"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="managerNameCHS" />
                      <FormControl>
                        <Input {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </KycFormRow>
              <KycFormRow>
                <FormField
                  control={form.control}
                  name="enterpriseInfo.managerSurname"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="managerSurname" />
                      <FormControl>
                        <Input
                          placeholder={kycPlaceholder("managerSurname")}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="enterpriseInfo.managerNameEN"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="managerNameEN" />
                      <FormControl>
                        <Input
                          placeholder={kycPlaceholder("managerNameEN")}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </KycFormRow>
              <KycFormRow>
                <FormField
                  control={form.control}
                  name="enterpriseInfo.managerCountry"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="managerCountry" />
                      <FormControl>
                        <KycCountrySelect {...field} />
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
                      <KycFormLabel fieldKey="managerAuthType" />
                      <FormControl>
                        <KycAuthTypeSelect
                          countryCode={managerCountry}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </KycFormRow>
              <FormField
                control={form.control}
                name="enterpriseInfo.managerIdCard"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <KycFormLabel fieldKey="managerIdCard" />
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <CertificateValidityFields
                control={form.control}
                startName="enterpriseInfo.managerCertificateStart"
                endName="enterpriseInfo.managerCertificateEnd"
                startLabelKey="managerCertificateTerm"
                endLabelKey="managerCertificateEnd"
              />
              <KycFormRow>
                <FormField
                  control={form.control}
                  name="enterpriseInfo.managerBirthday"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="managerBirthday" />
                      <FormControl>
                        <FormDatePicker {...field} />
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
                      <KycFormLabel fieldKey="managerGender" />
                      <FormControl>
                        <FormSelect {...field} options={GENDERS} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </KycFormRow>
              <KycFormRow>
                <FormField
                  control={form.control}
                  name="enterpriseInfo.managerResidenceCountry"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="managerResidenceCountry" />
                      <FormControl>
                        <KycCountrySelect {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="enterpriseInfo.managerContactsEmail"
                  render={({ field }) => (
                    <FormItem>
                      <KycFormLabel fieldKey="managerContactsEmail" required={false} />
                      <FormControl>
                        <Input
                          type="email"
                          placeholder={kycPlaceholder("managerContactsEmail")}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </KycFormRow>
              <FormField
                control={form.control}
                name="enterpriseInfo.managerResidenceAddress"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <KycFormLabel fieldKey="managerResidenceAddress" />
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.managerIdHoldingPaths"
                render={({ field }) => (
                  <FormItem className="sm:col-span-2">
                    <FormControl>
                      <KycFileUpload
                        label={getKycFieldMeta("managerIdHolding").label + " *"}
                        description={getKycFieldMeta("managerIdHolding").description}
                        paths={field.value ?? ["", "", ""]}
                        onChange={field.onChange}
                        minItems={3}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="enterpriseInfo.middleTierShareholders"
                render={({ field }) => (
                  <FormItem>
                    <KycFormLabel fieldKey="middleTierShareholders" />
                    <FormControl>
                      <FormSelect {...field} options={YES_NO} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              {middleTierShareholders === "Yes" && (
                <FormField
                  control={form.control}
                  name="enterpriseInfo.equityStructurePaths"
                  render={({ field }) => (
                    <FormItem className="sm:col-span-2">
                      <FormControl>
                        <KycFileUpload
                          label={getKycFieldMeta("equityStructure").label + " *"}
                          description={
                            getKycFieldMeta("equityStructure").description
                          }
                          paths={field.value ?? [""]}
                          onChange={field.onChange}
                          minItems={1}
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
              <CardContent className="space-y-4">
                {shareholders.fields.map((field, index) => {
                  const nameChs = form.watch(
                    `shareholdersInfo.${index}.nameCHS`
                  );
                  const surnameChs = form.watch(
                    `shareholdersInfo.${index}.surnameCHS`
                  );
                  const titleLabel =
                    [surnameChs, nameChs].filter(Boolean).join("") ||
                    `股东 ${index + 1}`;

                  return (
                    <CollapsibleSection
                      key={field.id}
                      title={titleLabel}
                      defaultOpen={index === 0}
                      actions={
                        shareholders.fields.length > 1 ? (
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            onClick={() => shareholders.remove(index)}
                          >
                            <Trash2 className="h-4 w-4 text-red-500" />
                          </Button>
                        ) : undefined
                      }
                    >
                      <PersonFields
                        control={form.control}
                        namePrefix={`shareholdersInfo.${index}`}
                        showShareholding
                      />
                    </CollapsibleSection>
                  );
                })}
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
              <CardContent className="space-y-4">
                {directors.fields.map((field, index) => {
                  const nameChs = form.watch(`directorInfo.${index}.nameCHS`);
                  const surnameChs = form.watch(
                    `directorInfo.${index}.surnameCHS`
                  );
                  const titleLabel =
                    [surnameChs, nameChs].filter(Boolean).join("") ||
                    `董事 ${index + 1}`;

                  return (
                    <CollapsibleSection
                      key={field.id}
                      title={titleLabel}
                      defaultOpen={index === 0}
                      actions={
                        directors.fields.length > 1 ? (
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            onClick={() => directors.remove(index)}
                          >
                            <Trash2 className="h-4 w-4 text-red-500" />
                          </Button>
                        ) : undefined
                      }
                    >
                      <PersonFields
                        control={form.control}
                        namePrefix={`directorInfo.${index}`}
                      />
                    </CollapsibleSection>
                  );
                })}
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
    </KycReferenceProvider>
  );
}
