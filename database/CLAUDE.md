# Database — MySQL

## 技术栈
- MySQL 8.0+
- golang-migrate（迁移工具）
- 字符集：utf8mb4_unicode_ci

## 目录约定
- `migrations/` — 迁移文件（up/down 成对）
- `seeds/` — 种子数据（开发/测试用）
- `docs/schema.md` — Schema 文档

## Migration 命名规范
- 格式：`{六位序号}_{描述}.up.sql` / `.down.sql`
- 示例：`000001_create_users_table.up.sql`

## 设计规范
- 金额字段使用 `DECIMAL(28,8)`（支持法币 2 位和 BTC 8 位小数）
- 时间字段使用 `DATETIME(3)`（毫秒精度）
- 软删除使用 `deleted_at` 字段
- 主键使用 `BIGINT UNSIGNED AUTO_INCREMENT`
- 必须有 `created_at`、`updated_at` 字段
- 敏感字段（手机号、银行卡号）加密存储

## 常用命令
```bash
# 执行迁移
migrate -path migrations -database "mysql://user:pass@tcp(localhost:3306)/motewallet" up

# 回滚一步
migrate -path migrations -database "mysql://user:pass@tcp(localhost:3306)/motewallet" down 1
```

## 注意事项
- 已发布的 migration 禁止修改，只能新增
- 每个 up 必须有对应的 down
- 大表变更考虑在线 DDL
