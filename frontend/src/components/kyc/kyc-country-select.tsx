"use client";

import { useQuery } from "@tanstack/react-query";
import { FormSearchableSelect } from "@/components/ui/form-searchable-select";
import {
  KYC_COUNTRY_SCENE_ADDRESS,
  kycCountriesQueryOptions,
} from "@/lib/kyc/reference-queries";
import type { KycCountryScene } from "@/types/kyc-reference";

interface KycCountrySelectProps {
  scene?: KycCountryScene;
  currency?: string;
  value?: string;
  onChange: (value: string) => void;
  onBlur?: () => void;
  placeholder?: string;
  disabled?: boolean;
  className?: string;
}

export function KycCountrySelect({
  scene = KYC_COUNTRY_SCENE_ADDRESS,
  currency,
  value,
  onChange,
  onBlur,
  placeholder = "请选择国家/地区",
  disabled,
  className,
}: KycCountrySelectProps) {
  const { data: countryOptions = [], isLoading, isError } = useQuery({
    ...kycCountriesQueryOptions(scene, currency),
    select: (response) =>
      response.items.map((item) => ({
        value: item.country_code,
        label: item.country_name,
      })),
  });

  return (
    <FormSearchableSelect
      value={value}
      onChange={onChange}
      onBlur={onBlur}
      options={countryOptions}
      placeholder={
        isLoading
          ? "加载中…"
          : isError
            ? "加载失败，请刷新"
            : placeholder
      }
      searchPlaceholder="搜索国家/地区"
      disabled={
        disabled ||
        isLoading ||
        isError ||
        countryOptions.length === 0
      }
      className={className}
    />
  );
}
