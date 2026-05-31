"use client";

import type { Control } from "react-hook-form";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { FormSelect } from "@/components/ui/form-select";
import { FilePathList } from "@/components/kyc/file-path-list";
import { GENDERS, REGISTER_REGIONS, VERIFICATION_TYPES } from "@/lib/kyc/constants";
import type { KycFormValues } from "@/lib/validations/onboarding";

interface PersonFieldsProps {
  control: Control<KycFormValues>;
  namePrefix: `shareholdersInfo.${number}` | `directorInfo.${number}`;
  showShareholding?: boolean;
}

export function PersonFields({
  control,
  namePrefix,
  showShareholding = false,
}: PersonFieldsProps) {
  return (
    <div className="grid gap-4 sm:grid-cols-2">
      <FormField
        control={control}
        name={`${namePrefix}.nameCHS`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>中文姓名</FormLabel>
            <FormControl>
              <Input placeholder="中文姓名" {...field} />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={control}
        name={`${namePrefix}.surname`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>英文姓 *</FormLabel>
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
            <FormLabel>英文名 *</FormLabel>
            <FormControl>
              <Input placeholder="Given name" {...field} />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={control}
        name={`${namePrefix}.gender`}
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
        control={control}
        name={`${namePrefix}.birthday`}
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
        control={control}
        name={`${namePrefix}.country`}
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
        control={control}
        name={`${namePrefix}.authType`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>证件类型 code *</FormLabel>
            <FormControl>
              <Input placeholder="鲲「国家认证类型」接口 docCode" {...field} />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={control}
        name={`${namePrefix}.idCard`}
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
        control={control}
        name={`${namePrefix}.certificateStart`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>证件有效期起</FormLabel>
            <FormControl>
              <Input type="date" {...field} />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={control}
        name={`${namePrefix}.certificateEnd`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>证件有效期止</FormLabel>
            <FormControl>
              <Input type="date" {...field} />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={control}
        name={`${namePrefix}.residenceCountry`}
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
        control={control}
        name={`${namePrefix}.residenceAddress`}
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
      {showShareholding && (
        <FormField
          control={control}
          name={`${namePrefix}.shareholdingRatio`}
          render={({ field }) => (
            <FormItem>
              <FormLabel>持股比例 *</FormLabel>
              <FormControl>
                <Input placeholder="例如 25" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      )}
      <FormField
        control={control}
        name={`${namePrefix}.verificationType`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>验证方式</FormLabel>
            <FormControl>
              <FormSelect {...field} options={VERIFICATION_TYPES} />
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
              <FilePathList
                label="手持证件 / 证件照文件 path"
                description="idHolding 方式至少 3 个 path（证件照、自拍、手持照）"
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
  );
}
