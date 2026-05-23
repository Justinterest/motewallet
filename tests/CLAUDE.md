# Tests — 测试

## 测试分层
- **单元测试** — 在各自项目目录内（`frontend/`, `backend/`）
- **E2E 测试** — 在 `tests/e2e/` 目录
- **集成测试** — 在 `backend/` 内，使用 testcontainers

## 目录约定
- `e2e/` — Playwright / Shiplight E2E 测试
- E2E 测试文件命名：`*.test.yaml`（Shiplight）或 `*.spec.ts`（Playwright）

## 技术栈
| 层级 | 工具 |
|------|------|
| 前端单元 | Vitest + React Testing Library |
| 后端单元 | Go testing + testify |
| 后端集成 | testcontainers-go (MySQL) |
| E2E | Playwright + Shiplight YAML |

## 运行命令
```bash
# 前端单元测试
cd frontend && pnpm test

# 后端单元测试
cd backend && go test ./...

# E2E 测试
cd tests && npx playwright test
# 或
cd tests && npx shiplight test
```

## 注意事项
- 测试必须幂等，可重复运行
- 测试间不能有执行顺序依赖
- 金额断言使用精确值，不用近似
- 集成测试使用 testcontainers（真实 MySQL），不 mock 数据库
- E2E 测试前需要初始化测试数据
