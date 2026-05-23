# Frontend — Next.js

## 技术栈
- Next.js 15+ (App Router)
- React 19+, TypeScript
- Tailwind CSS 4+, shadcn/ui
- React Hook Form + Zod（表单验证）
- TanStack Query（服务端状态）
- Zustand（客户端状态）
- Axios（HTTP 请求）

## 目录约定
- `src/app/` — 页面路由（App Router）
- `src/components/ui/` — shadcn/ui 基础组件
- `src/components/business/` — 业务组件
- `src/lib/api/` — API 请求封装
- `src/lib/hooks/` — 自定义 Hooks
- `src/lib/validations/` — Zod Schema
- `src/stores/` — Zustand Store
- `src/types/` — TypeScript 类型

## 命名规范
- 组件：PascalCase（`WithdrawalForm`）
- 文件：kebab-case（`withdrawal-form.tsx`）
- Hook：camelCase 以 `use` 开头（`useWithdrawal`）

## 开发命令
```bash
pnpm dev          # 启动开发服务器
pnpm build        # 构建生产版本
pnpm lint         # ESLint 检查
pnpm test         # 运行单元测试
```

## 注意事项
- 默认使用 Server Components，需要交互时才加 `"use client"`
- API 请求统一通过 `lib/api/` 封装
- 金额显示使用 `formatAmount()` 工具函数，单位换算（分→元）
- JWT Token 存储在 httpOnly Cookie，不用 localStorage
- 后端 API 契约见 `docs/api/openapi.yaml`
