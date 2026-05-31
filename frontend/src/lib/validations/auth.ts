import { z } from "zod";

export const loginSchema = z.object({
  email: z
    .string()
    .min(1, { message: "请输入邮箱地址" })
    .email({ message: "请输入有效的邮箱地址" }),
  password: z.string().min(1, { message: "请输入密码" }),
});

export type LoginFormValues = z.infer<typeof loginSchema>;

export const registerSchema = z
  .object({
    email: z
      .string()
      .min(1, { message: "请输入邮箱地址" })
      .email({ message: "请输入有效的邮箱地址" }),
    verificationCode: z
      .string()
      .length(6, { message: "请输入 6 位验证码" })
      .regex(/^\d{6}$/, { message: "验证码必须为 6 位数字" }),
    password: z
      .string()
      .min(8, { message: "密码长度至少为 8 位" }),
    confirmPassword: z
      .string()
      .min(1, { message: "请确认密码" }),
    agreeTerms: z.literal(true, {
      errorMap: () => ({ message: "请同意服务条款" }),
    }),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "两次输入的密码不一致",
    path: ["confirmPassword"],
  });

export type RegisterFormValues = z.infer<typeof registerSchema>;
