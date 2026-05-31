"use client";

import { FormLabel } from "@/components/ui/form";
import { FieldExampleLink, FieldHint } from "@/components/kyc/field-hint";
import { getKycFieldMeta, type KycFieldMetaKey } from "@/lib/kyc/field-meta";

interface KycFormLabelProps {
  fieldKey: KycFieldMetaKey;
  /** Override doc required flag (e.g. conditional HK fields) */
  required?: boolean;
}

export function KycFormLabel({ fieldKey, required }: KycFormLabelProps) {
  const meta = getKycFieldMeta(fieldKey);
  const isRequired = required ?? meta.required;

  return (
    <div className="space-y-1">
      <div className="flex flex-wrap items-center gap-x-2 gap-y-0.5">
        <FormLabel className="mb-0">
          {meta.label}
          {isRequired ? " *" : ""}
        </FormLabel>
        <FieldExampleLink href={meta.exampleImage} />
      </div>
      <FieldHint meta={{ description: meta.description }} />
    </div>
  );
}

export function kycPlaceholder(fieldKey: KycFieldMetaKey): string | undefined {
  return getKycFieldMeta(fieldKey).placeholder;
}
