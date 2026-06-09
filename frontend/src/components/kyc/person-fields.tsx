"use client";

import { useEffect, useRef } from "react";
import type { Control } from "react-hook-form";
import { useFormContext, useWatch } from "react-hook-form";
import { Copy } from "lucide-react";
import { Button } from "@/components/ui/button";
import { toast } from "@/hooks/use-toast";
import { managerInfoToPerson } from "@/lib/kyc/copy-manager";
import {
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { FormDatePicker } from "@/components/ui/date-picker";
import { FormSelect } from "@/components/ui/form-select";
import { Input } from "@/components/ui/input";
import { CertificateValidityFields } from "@/components/kyc/certificate-validity-fields";
import { KycAuthTypeSelect } from "@/components/kyc/kyc-auth-type-select";
import { KycCountrySelect } from "@/components/kyc/kyc-country-select";
import { KYC_COUNTRY_SCENE_NATIONALITY } from "@/lib/kyc/reference-queries";
import { KycFileUpload } from "@/components/kyc/kyc-file-upload";
import { KycFormRow } from "@/components/kyc/kyc-form-row";
import { KycFormLabel, kycPlaceholder } from "@/components/kyc/kyc-form-label";
import { getKycFieldMeta } from "@/lib/kyc/field-meta";
import { GENDERS } from "@/lib/kyc/constants";
import type { KycFormValues } from "@/lib/validations/onboarding";

interface PersonFieldsProps {
  control: Control<KycFormValues>;
  namePrefix: `shareholdersInfo.${number}` | `directorInfo.${number}`;
  showShareholding?: boolean;
  /** Show button to copy account manager info into this person. */
  showCopyFromManager?: boolean;
}

export function PersonFields({
  control,
  namePrefix,
  showShareholding = false,
  showCopyFromManager = false,
}: PersonFieldsProps) {
  const { setValue, getValues } = useFormContext<KycFormValues>();
  const personCountry = useWatch({
    control,
    name: `${namePrefix}.country`,
  });
  const prevCountryRef = useRef(personCountry);

  useEffect(() => {
    if (
      prevCountryRef.current !== undefined &&
      prevCountryRef.current !== personCountry
    ) {
      setValue(`${namePrefix}.authType`, "");
    }
    prevCountryRef.current = personCountry;
  }, [personCountry, namePrefix, setValue]);

  function copyFromManager() {
    const enterprise = getValues("enterpriseInfo");
    const person = managerInfoToPerson(enterprise);

    if (showShareholding) {
      const shareholdingRatio = getValues(`${namePrefix}.shareholdingRatio`);
      setValue(
        namePrefix,
        { ...person, shareholdingRatio: shareholdingRatio ?? "" },
        { shouldDirty: true }
      );
      toast({
        title: "已复制管理人信息",
        description: "请核对持股比例等股东专有字段。",
      });
      return;
    }

    setValue(namePrefix, person, { shouldDirty: true });
    toast({
      title: "已复制管理人信息",
      description: "身份与证件信息已填入，请核对后提交。",
    });
  }

  return (
    <div className="space-y-4">
      {showCopyFromManager ? (
        <div className="flex justify-end">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={copyFromManager}
          >
            <Copy className="mr-1 h-4 w-4" />
            从管理人复制
          </Button>
        </div>
      ) : null}
      <div className="grid gap-4 sm:grid-cols-2">
      <KycFormRow>
        <FormField
          control={control}
          name={`${namePrefix}.surnameCHS`}
          render={({ field }) => (
            <FormItem>
            <KycFormLabel fieldKey="person.surnameCHS" />
            <FormControl>
              <Input placeholder={kycPlaceholder("person.surnameCHS")} {...field} />
            </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={control}
          name={`${namePrefix}.nameCHS`}
          render={({ field }) => (
            <FormItem>
              <KycFormLabel fieldKey="person.nameCHS" />
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
          control={control}
          name={`${namePrefix}.surname`}
          render={({ field }) => (
            <FormItem>
              <KycFormLabel fieldKey="person.surname" />
              <FormControl>
                <Input placeholder="Surname" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={control}
          name={`${namePrefix}.nameEN`}
          render={({ field }) => (
            <FormItem>
              <KycFormLabel fieldKey="person.nameEN" />
              <FormControl>
                <Input placeholder="Given name" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      </KycFormRow>

      <KycFormRow>
        <FormField
          control={control}
          name={`${namePrefix}.gender`}
          render={({ field }) => (
            <FormItem>
              <KycFormLabel fieldKey="person.gender" />
              <FormControl>
                <FormSelect {...field} options={GENDERS} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={control}
          name={`${namePrefix}.birthday`}
          render={({ field }) => (
            <FormItem>
              <KycFormLabel fieldKey="person.birthday" />
              <FormControl>
                <FormDatePicker {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      </KycFormRow>

      <KycFormRow>
        <FormField
          control={control}
          name={`${namePrefix}.country`}
          render={({ field }) => (
            <FormItem>
              <KycFormLabel fieldKey="person.country" />
              <FormControl>
                <KycCountrySelect scene={KYC_COUNTRY_SCENE_NATIONALITY} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={control}
          name={`${namePrefix}.authType`}
          render={({ field }) => (
            <FormItem>
              <KycFormLabel fieldKey="person.authType" />
              <FormControl>
                <KycAuthTypeSelect
                  countryCode={personCountry}
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      </KycFormRow>

      <FormField
        control={control}
        name={`${namePrefix}.idCard`}
        render={({ field }) => (
          <FormItem className="sm:col-span-2">
            <KycFormLabel fieldKey="person.idCard" />
            <FormControl>
              <Input {...field} />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />

      <CertificateValidityFields
        control={control}
        startName={`${namePrefix}.certificateStart`}
        endName={`${namePrefix}.certificateEnd`}
        startLabelKey="person.certificateTerm"
        endLabelKey="person.certificateEnd"
      />

      {showShareholding ? (
        <KycFormRow>
          <FormField
            control={control}
            name={`${namePrefix}.residenceCountry`}
            render={({ field }) => (
              <FormItem>
                <KycFormLabel fieldKey="person.residenceCountry" />
                <FormControl>
                  <KycCountrySelect scene={KYC_COUNTRY_SCENE_NATIONALITY} {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={control}
            name={`${namePrefix}.shareholdingRatio`}
            render={({ field }) => (
              <FormItem>
                <KycFormLabel fieldKey="person.shareholdingRatio" />
                <FormControl>
                  <Input
                    placeholder={kycPlaceholder("person.shareholdingRatio")}
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </KycFormRow>
      ) : (
        <FormField
          control={control}
          name={`${namePrefix}.residenceCountry`}
          render={({ field }) => (
            <FormItem>
              <KycFormLabel fieldKey="person.residenceCountry" />
              <FormControl>
                <KycCountrySelect scene={KYC_COUNTRY_SCENE_NATIONALITY} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      )}

      <FormField
        control={control}
        name={`${namePrefix}.residenceAddress`}
        render={({ field }) => (
          <FormItem className="sm:col-span-2">
            <KycFormLabel fieldKey="person.residenceAddress" />
            <FormControl>
              <Input {...field} />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={control}
        name={`${namePrefix}.idHoldingPaths`}
        render={({ field }) => (
          <FormItem className="sm:col-span-2">
            <FormControl>
              <KycFileUpload
                label={getKycFieldMeta("person.idHolding").label + " *"}
                description={getKycFieldMeta("person.idHolding").description}
                paths={field.value ?? [""]}
                onChange={field.onChange}
                minItems={1}
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
      </div>
    </div>
  );
}
