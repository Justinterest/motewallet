# Motewallet Withdrawal — 钱包提现系统

## 项目概述

钱包提现系统，支持用户发起提现请求、审批流程、资金结算和状态追踪。

## 技术栈

| 层级 | 技术 | 版本要求 |
|------|------|---------|
| 前端 | Next.js + React + TypeScript | Next.js 15+, React 19+ |
| UI | Tailwind CSS + shadcn/ui | Tailwind 4+ |
| 后端 | Golang + Gin | Go 1.22+ |
| ORM | GORM | v2 |
| 数据库 | MySQL | 8.0+ |
| 认证 | JWT | - |
| API 规范 | OpenAPI 3.0 | - |

## 项目结构

```
├── frontend/          # 商户端 Next.js 前端（localhost:3000）
├── admin/             # 管理后台 Next.js 前端（localhost:3001）
├── backend/           # Golang 后端 API（localhost:8080）
├── database/          # MySQL 迁移和种子数据
├── tests/             # E2E 和集成测试
├── docs/
│   ├── prd/           # 产品需求文档
│   ├── design/        # 设计文档
│   └── api/           # API 契约（OpenAPI）
└── .claude/agents/    # Agent 团队定义
```

## Agent 团队

本项目使用 7 个专业化 Agent 覆盖完整开发生命周期：

| Agent | 调用方式 | 职责 |
|-------|---------|------|
| Product Manager | `@product-manager` | 需求分析、PRD、用户故事、验收标准 |
| UI/UX Designer | `@ui-ux-designer` | 设计系统、页面设计、组件规范 |
| Frontend Developer | `@frontend-developer` | Next.js 页面和组件开发 |
| Backend Developer | `@backend-developer` | Golang API 和业务逻辑 |
| Database Engineer | `@database-engineer` | MySQL Schema、Migration、查询优化 |
| Test Engineer | `@test-engineer` | 单元测试、集成测试、E2E 测试 |
| QA / Acceptance | `@qa-acceptance` | 功能验收、浏览器验证、验收报告 |

## 协作流程

```
需求 → @product-manager 输出 PRD
     → @ui-ux-designer 输出设计稿
     → @frontend-developer + @backend-developer + @database-engineer 并行开发
     → @test-engineer 编写测试
     → @qa-acceptance 验收
```

## 关键约定

### API 契约优先
- API 定义在 `docs/api/openapi.yaml`
- 前后端通过 OpenAPI 契约解耦，可并行开发
- 任何 API 变更必须先更新契约文档

### 编码规范
- **前端**: ESLint + Prettier, 组件使用 PascalCase, 文件使用 kebab-case
- **后端**: gofmt + golint, 遵循 Go 标准项目布局
- **数据库**: Migration 文件按时间戳命名，禁止直接修改已发布的 migration
- **Git**: Conventional Commits（feat/fix/docs/test/refactor）

### 分支策略
- `main` — 生产分支
- `develop` — 开发分支
- `feature/*` — 功能分支
- `hotfix/*` — 紧急修复

### 安全要求
- 所有 API 需要 JWT 认证（登录接口除外）
- 提现操作需要二次验证
- 敏感数据（密钥、密码）禁止硬编码，使用环境变量
- SQL 查询使用参数化，防止注入
