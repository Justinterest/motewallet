"use client";

import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { XIcon } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { FormSearchableSelect } from "@/components/ui/form-searchable-select";
import {
  KYC_COUNTRY_SCENE_ADDRESS,
  kycCountriesQueryOptions,
} from "@/lib/kyc/reference-queries";
import type { KycCountryScene } from "@/types/kyc-reference";

function parseCommaSeparated(value?: string): string[] {
  return value
    ? value
        .split(",")
        .map((s) => s.trim())
        .filter(Boolean)
    : [];
}

interface KycCountryMultiSelectProps {
  scene?: KycCountryScene;
  value?: string;
  onChange: (value: string) => void;
  onBlur?: () => void;
  placeholder?: string;
  disabled?: boolean;
  className?: string;
}

export function KycCountryMultiSelect({
  scene = KYC_COUNTRY_SCENE_ADDRESS,
  value,
  onChange,
  onBlur,
  placeholder = "添加国家/地区",
  disabled,
  className,
}: KycCountryMultiSelectProps) {
  const selected = useMemo(() => parseCommaSeparated(value), [value]);

  const { data: countryOptions = [], isLoading, isError } = useQuery({
    ...kycCountriesQueryOptions(scene),
    select: (response) =>
      response.items.map((item) => ({
        value: item.country_code,
        label: item.country_name,
      })),
  });

  const labelByCode = useMemo(() => {
    const map = new Map<string, string>();
    for (const option of countryOptions) {
      map.set(option.value, option.label);
    }
    return map;
  }, [countryOptions]);

  const availableOptions = countryOptions.filter(
    (option) => !selected.includes(option.value)
  );

  function addCountry(code: string) {
    if (!code || selected.includes(code)) return;
    onChange([...selected, code].join(","));
  }

  function removeCountry(code: string) {
    onChange(selected.filter((item) => item !== code).join(","));
  }

  const selectPlaceholder = isLoading
    ? "加载中…"
    : isError
      ? "加载失败，请刷新"
      : placeholder;

  const selectDisabled =
    disabled ||
    isLoading ||
    isError ||
    countryOptions.length === 0 ||
    availableOptions.length === 0;

  return (
    <div className="space-y-2">
      {selected.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {selected.map((code) => (
            <Badge key={code} variant="secondary" className="gap-1 pr-1">
              {labelByCode.get(code) ?? code}
              <button
                type="button"
                className="hover:bg-muted rounded-sm p-0.5"
                aria-label={`移除 ${labelByCode.get(code) ?? code}`}
                onClick={() => removeCountry(code)}
                disabled={disabled}
              >
                <XIcon className="size-3" />
              </button>
            </Badge>
          ))}
        </div>
      )}
      <FormSearchableSelect
        value=""
        onChange={addCountry}
        onBlur={onBlur}
        options={availableOptions}
        placeholder={selectPlaceholder}
        searchPlaceholder="搜索国家/地区"
        disabled={selectDisabled}
        className={className}
      />
    </div>
  );
}
