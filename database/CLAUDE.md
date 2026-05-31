# Database — MySQL

## 技术栈
- MySQL 8.0+
- golang-migrate（迁移工具）
- 字符集：utf8mb4_unicode_ci

## 目录约定
- `migrations/` — Schema 迁移（golang-migrate，up/down 成对）
- `seeds/` — 种子数据（纯 SQL，开发/测试用，**不走** migrate）
- `docs/schema.md` — Schema 文档

## 本地初始化流程

```text
创建数据库 → migrate up → 执行 seeds → 启动 backend
```

连接参数与 `backend/.env` 对齐（默认库名 `motewallet`）。

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

## 前置依赖

```bash
# macOS
brew install golang-migrate mysql-client
```

在 `database/` 目录下执行下文命令。

## Schema 迁移（migrations）

```bash
# 建库（首次）
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS motewallet CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 执行全部迁移
migrate -path migrations -database "mysql://root:密码@tcp(localhost:3306)/motewallet" up

# 查看当前版本
migrate -path migrations -database "mysql://root:密码@tcp(localhost:3306)/motewallet" version

# 回滚一步
migrate -path migrations -database "mysql://root:密码@tcp(localhost:3306)/motewallet" down 1

# 回滚到版本 0（会执行所有 down，DROP 表）
migrate -path migrations -database "mysql://root:密码@tcp(localhost:3306)/motewallet" down
```

## 种子数据（seeds）

### 与 migrate 的区别

| | migrations | seeds |
|---|------------|-------|
| 工具 | golang-migrate | `mysql` 客户端 |
| 文件 | `000001_xxx.up.sql` / `.down.sql` | `001_xxx.sql`（无 down） |
| 用途 | 表结构变更 | 开发/测试初始数据 |
| 顺序 | 先执行 | **必须在 migrate up 之后** |

### 命名规范

- 格式：`{三位序号}_{描述}.sql`
- 示例：`001_seed_admin_and_config.sql`

### 执行命令

```bash
# 在 database/ 目录下
mysql -u root -p motewallet < seeds/001_seed_admin_and_config.sql
```

或在 MySQL 交互式客户端中：

```sql
USE motewallet;
SOURCE /path/to/database/seeds/001_seed_admin_and_config.sql;
```

### 当前种子内容

`001_seed_admin_and_config.sql` 插入：

- `admin_users` — 超管 `superadmin` / `admin@motewallet.com`（密码见文件头注释，仅开发用）
- `fee_templates` 及兑换/提现手续费子表 — 默认费率模板
- `system_configs` — KUN、币种、平台名等 KV 配置

### 编写注意

- 种子仅供开发/测试，**不要**用于生产数据修复（生产用独立 migration 或运维脚本）
- 重复执行可能因主键/唯一约束失败；需要重置时先 `migrate down` 再 `up`，再跑 seed
- 新种子文件优先使用幂等写法（如 `INSERT ... ON DUPLICATE KEY UPDATE`），便于本地重复导入

## 注意事项

- 已发布的 migration 禁止修改，只能新增
- 每个 up 必须有对应的 down
- 大表变更考虑在线 DDL
- Schema 变更需同步更新 `docs/schema.md`
