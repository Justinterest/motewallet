"use client";

import { FormSelect } from "@/components/ui/form-select";
import { useKycReference } from "@/components/kyc/kyc-reference-provider";

interface KycCountrySelectProps {
  value?: string;
  onChange: (value: string) => void;
  onBlur?: () => void;
  placeholder?: string;
  disabled?: boolean;
  className?: string;
}

export function KycCountrySelect({
  value,
  onChange,
  onBlur,
  placeholder = "请选择国家/地区",
  disabled,
  className,
}: KycCountrySelectProps) {
  const { countryOptions, countriesLoading, countriesError } = useKycReference();

  return (
    <FormSelect
      value={value}
      onChange={onChange}
      onBlur={onBlur}
      options={countryOptions}
      placeholder={
        countriesLoading
          ? "加载中…"
          : countriesError
            ? "加载失败，请刷新"
            : placeholder
      }
      disabled={
        disabled ||
        countriesLoading ||
        countriesError ||
        countryOptions.length === 0
      }
      className={className}
    />
  );
}
