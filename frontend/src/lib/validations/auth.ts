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
