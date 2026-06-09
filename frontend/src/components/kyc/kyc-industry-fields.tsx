"use client";

import type { Control } from "react-hook-form";
import { useFormContext, useWatch } from "react-hook-form";

import {
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { FormSelect } from "@/components/ui/form-select";
import { KycFormLabel } from "@/components/kyc/kyc-form-label";
import { KycFormRow } from "@/components/kyc/kyc-form-row";
import {
  getSubIndustryOptions,
  INDUSTRY_LEVEL1_OPTIONS,
} from "@/lib/kyc/industries";
import type { KycFormValues } from "@/lib/validations/onboarding";

interface KycIndustryFieldsProps {
  control: Control<KycFormValues>;
}

export function KycIndustryFields({ control }: KycIndustryFieldsProps) {
  const { getValues } = useFormContext<KycFormValues>();
  const watchedIndustry = useWatch({
    control,
    name: "enterpriseInfo.industry",
  });
  const industry =
    watchedIndustry || getValues("enterpriseInfo.industry") || "";
  const subIndustryOptions = getSubIndustryOptions(industry);

  return (
    <KycFormRow>
      <FormField
        control={control}
        name="enterpriseInfo.industry"
        render={({ field }) => (
          <FormItem>
            <KycFormLabel fieldKey="industry" />
            <FormControl>
              <FormSelect
                {...field}
                options={INDUSTRY_LEVEL1_OPTIONS}
                placeholder="请选择一级行业"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={control}
        name="enterpriseInfo.subIndustry"
        render={({ field }) => (
          <FormItem>
            <KycFormLabel fieldKey="subIndustry" />
            <FormControl>
              <FormSelect
                {...field}
                options={subIndustryOptions}
                placeholder={
                  industry ? "请选择二级行业" : "请先选择一级行业"
                }
                disabled={!industry || subIndustryOptions.length === 0}
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
    </KycFormRow>
  );
}
