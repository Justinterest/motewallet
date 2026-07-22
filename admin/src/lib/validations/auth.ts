import { z } from "zod";

export const adminLoginSchema = z.object({
  username: z.string().min(1, "请输入用户名"),
  password: z.string().min(1, "请输入密码"),
});

export type AdminLoginFormValues = z.infer<typeof adminLoginSchema>;

export const totpCodeSchema = z.object({
  code: z
    .string()
    .length(6, { message: "请输入 6 位验证码" })
    .regex(/^\d{6}$/, { message: "验证码必须为 6 位数字" }),
});

export type TotpCodeFormValues = z.infer<typeof totpCodeSchema>;

export const changePasswordSchema = z
  .object({
    newPassword: z.string().min(8, "新密码至少 8 位"),
    confirmPassword: z.string().min(1, "请确认新密码"),
  })
  .refine((data) => data.newPassword === data.confirmPassword, {
    message: "两次输入的密码不一致",
    path: ["confirmPassword"],
  });

export type ChangePasswordFormValues = z.infer<typeof changePasswordSchema>;

export const createEmployeeSchema = z.object({
  username: z
    .string()
    .min(3, "用户名至少 3 个字符")
    .max(64, "用户名最多 64 个字符"),
  email: z
    .string()
    .min(1, "请输入邮箱")
    .email("请输入有效的邮箱地址"),
  role: z.enum(["SUPER_ADMIN", "OPERATOR", "FINANCE"], {
    required_error: "请选择角色",
  }),
  password: z
    .string()
    .optional()
    .refine((v) => !v || v.length >= 8, {
      message: "密码至少 8 位",
    }),
});

export type CreateEmployeeFormValues = z.infer<typeof createEmployeeSchema>;
