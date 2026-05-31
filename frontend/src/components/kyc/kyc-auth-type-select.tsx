"use client";

import { FormSelect } from "@/components/ui/form-select";
import { useKycCountryAuthTypes } from "@/lib/hooks/use-onboarding";

interface KycAuthTypeSelectProps {
  countryCode?: string;
  value?: string;
  onChange: (value: string) => void;
  onBlur?: () => void;
  placeholder?: string;
  disabled?: boolean;
  className?: string;
}

export function KycAuthTypeSelect({
  countryCode,
  value,
  onChange,
  onBlur,
  placeholder,
  disabled,
  className,
}: KycAuthTypeSelectProps) {
  const { data: options = [], isLoading, isError } =
    useKycCountryAuthTypes(countryCode);

  const resolvedPlaceholder = !countryCode
    ? "请先选择国家/地区"
    : isLoading
      ? "加载中…"
      : isError
        ? "加载失败，请刷新"
        : options.length === 0
          ? "暂无可用证件类型"
          : (placeholder ?? "请选择证件类型");

  return (
    <FormSelect
      value={value}
      onChange={onChange}
      onBlur={onBlur}
      options={options}
      placeholder={resolvedPlaceholder}
      disabled={
        disabled ||
        !countryCode ||
        isLoading ||
        isError ||
        options.length === 0
      }
      className={className}
    />
  );
}
