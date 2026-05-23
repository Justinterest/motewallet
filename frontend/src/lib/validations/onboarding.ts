import { z } from "zod";

export const kycFormSchema = z.object({
  company_name: z.string().min(2, "公司名称至少2个字符"),
  country: z.string().min(2, "请选择国家/地区"),
  registration_number: z.string().min(1, "请输入注册编号"),
  contact_name: z.string().min(2, "联系人姓名至少2个字符"),
  contact_phone: z.string().min(6, "请输入有效的手机号"),
});

export type KycFormValues = z.infer<typeof kycFormSchema>;
