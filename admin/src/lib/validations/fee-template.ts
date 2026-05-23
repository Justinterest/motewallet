import { z } from "zod";

export const feeTemplateFormSchema = z.object({
  name: z.string().min(1, "模板名称不能为空"),
  description: z.string().optional(),
  is_default: z.boolean(),
});

export type FeeTemplateFormValues = z.infer<typeof feeTemplateFormSchema>;
