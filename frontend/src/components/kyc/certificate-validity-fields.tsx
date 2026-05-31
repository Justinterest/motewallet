"use client";

import type { Control, FieldPath } from "react-hook-form";
import {
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { FormDatePicker } from "@/components/ui/date-picker";
import { KycFormLabel } from "@/components/kyc/kyc-form-label";
import { KycFormRow } from "@/components/kyc/kyc-form-row";
import type { KycFieldMetaKey } from "@/lib/kyc/field-meta";
import type { KycFormValues } from "@/lib/validations/onboarding";

interface CertificateValidityFieldsProps {
  control: Control<KycFormValues>;
  startName: FieldPath<KycFormValues>;
  endName: FieldPath<KycFormValues>;
  startLabelKey: KycFieldMetaKey;
  endLabelKey: KycFieldMetaKey;
}

export function CertificateValidityFields({
  control,
  startName,
  endName,
  startLabelKey,
  endLabelKey,
}: CertificateValidityFieldsProps) {
  return (
    <KycFormRow>
      <FormField
        control={control}
        name={startName}
        render={({ field }) => (
          <FormItem>
            <KycFormLabel fieldKey={startLabelKey} />
            <FormControl>
              <FormDatePicker {...field} />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={control}
        name={endName}
        render={({ field }) => (
          <FormItem>
            <KycFormLabel fieldKey={endLabelKey} />
            <FormControl>
              <FormDatePicker {...field} />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
    </KycFormRow>
  );
}
