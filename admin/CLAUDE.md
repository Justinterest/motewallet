# Admin — 管理后台前端

## 技术栈
- Next.js 15+, React 19+, TypeScript
- Tailwind CSS 4+, shadcn/ui
- TanStack Query（服务端状态）, Zustand（客户端状态）
- React Hook Form + Zod（表单验证）
- Axios（HTTP 请求）

## 目录约定
- `src/app/` — App Router 页面
- `src/components/ui/` — shadcn/ui 基础组件
- `src/components/layout/` — 布局组件（Sidebar、Header）
- `src/lib/api/` — API 请求封装
- `src/lib/hooks/` — 自定义 Hooks
- `src/stores/` — Zustand Store
- `src/types/` — TypeScript 类型

## 界面语言
- 所有 UI 文案使用中文

## 设计参考
- `docs/design/design-admin-pages.md` — 管理后台页面设计
- `docs/design/design-system.md` — 设计系统
- API 端点前缀: `/api/v1/admin/`

## 端口
- 开发环境: localhost:3001 (`pnpm dev -p 3001`)

## 约束
- 只修改 `admin/` 目录下的文件
- JWT token 存储在 httpOnly Cookie (`admin_token`)
- 参考 `docs/api/openapi.yaml` 对接 API
