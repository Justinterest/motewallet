# Motewallet 数据库 Schema 文档

## 设计原则

- **金额存储**：使用 `DECIMAL(28,8)` 统一存储所有币种金额，支持法币（2 位小数）和数字货币（BTC 8 位小数）
  - 选择 DECIMAL 而非 BIGINT 的原因：平台同时处理法币和加密货币，精度需求不同（USD 2 位 vs BTC 8 位），使用 DECIMAL 可避免应用层对不同币种做不同的倍率转换
- **资金分层**：
  - `merchant_wallets` — 当前余额快照（balance / frozen_balance）
  - `wallet_ledger` — 钱包变动账本（每次余额/冻结变更一行，含变动前后余额，只追加）
  - `transaction_records` + `*_orders` — 业务流水与业务明细（商户交易列表、手续费、审核等）
- **软删除**：使用 `deleted_at` 字段（DATETIME(3) NULL）
- **时间精度**：统一 `DATETIME(3)` 毫秒级
- **字符集**：`utf8mb4_unicode_ci`
- **主键**：`BIGINT UNSIGNED AUTO_INCREMENT`
- **敏感数据**：银行账号、KUN 密钥等字段 AES 加密存储

## 表清单

核心表在 `000001_init_schema.up.sql` 中创建；后续变更见后续 migration。

| # | 表名 | 描述 |
|---|------|------|
| 1 | fee_templates | 手续费模板 |
| 2 | merchants | 商户 |
| 3 | fee_template_exchange_items | 兑换手续费配置 |
| 4 | fee_template_crypto_withdrawal_items | 数币提现手续费配置 |
| 5 | fee_template_fiat_withdrawal_items | 法币提现手续费配置 |
| 6 | merchant_wallets | 商户钱包余额（当前快照） |
| 7 | crypto_addresses | 数币提现白名单地址 |
| 8 | bank_accounts | 法币提现银行账户 |
| 9 | transaction_records | 平台业务资金流水 |
| 10 | admin_users | 管理员账号 |
| 11 | audit_logs | 操作审计日志 |
| 12 | system_configs | 系统配置（KV） |
| 13 | webhook_logs | KUN Webhook 回调日志 |
| 14 | system_announcements | 系统公告 |
| 15 | deposit_orders | 充值订单 |
| 16 | withdrawal_orders | 提现订单（审核、目标、链上信息） |
| 17 | exchange_orders | 兑换订单 |
| 18 | transfer_orders | 划转订单 |
| 19 | wallet_ledger | 钱包资金变动账本（只追加） |

---

## 表结构详解

### 1. merchants（商户表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK, AUTO_INCREMENT | 商户 ID |
| email | VARCHAR(128) | UNIQUE, NOT NULL | 登录邮箱 |
| password_hash | VARCHAR(255) | NOT NULL | bcrypt 密码哈希 |
| kun_sub_customer_no | VARCHAR(64) | UNIQUE | KUN 子商户号 |
| fee_template_id | BIGINT UNSIGNED | FK → fee_templates | 绑定的手续费模板 |
| supported_crypto_currencies | VARCHAR(128) | NULL | 商户支持的数币 CSV；NULL = 继承系统默认 |
| supported_fiat_currencies | VARCHAR(128) | NULL | 商户支持的法币 CSV；NULL = 继承系统默认 |
| supported_crypto_chains | TEXT | NULL | 商户支持的链 JSON `{USDT:[...]}`；NULL = 继承系统默认 |
| default_crypto_chains | TEXT | NULL | 商户默认链 JSON `{USDT:"TRX_TRC20"}`；NULL = 继承系统默认 |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'PENDING_AGREEMENT' | 商户状态 |
| kyc_auth_id | VARCHAR(128) | | KUN KYC 认证 ID |
| kyc_status | VARCHAR(20) | DEFAULT 'NONE' | KYC 状态 |
| kyc_fail_reason | TEXT | | KYC 失败原因 |
| kyc_submitted_at | DATETIME(3) | | KYC 提交时间 |
| kyc_completed_at | DATETIME(3) | | KYC 完成时间 |
| agreement_signed_at | DATETIME(3) | | 协议签署时间 |
| frozen_at | DATETIME(3) | | 冻结时间 |
| created_at | DATETIME(3) | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) | |
| updated_at | DATETIME(3) | NOT NULL, DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE | |
| deleted_at | DATETIME(3) | NULL | 软删除 |

**status 枚举值：**
- `PENDING_AGREEMENT` — 待签署协议
- `PENDING_KYC` — 待提交 KYC
- `KYC_REVIEWING` — KYC 审核中
- `ACTIVE` — 已激活（KYC 通过）
- `KYC_FAILED` — KYC 失败
- `FROZEN` — 已冻结

**kyc_status 枚举值：**
- `NONE` — 未提交
- `AUTHING` — 审核中
- `AUTH_SUC` — 认证成功
- `AUTH_FAIL` — 认证失败

---

### 2. fee_templates（手续费模板表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | 模板 ID |
| name | VARCHAR(64) | NOT NULL | 模板名称 |
| description | VARCHAR(255) | | 描述 |
| is_default | TINYINT(1) | NOT NULL, DEFAULT 0 | 是否默认模板 |
| created_at | DATETIME(3) | | |
| updated_at | DATETIME(3) | | |
| deleted_at | DATETIME(3) | NULL | |

---

### 3. fee_template_exchange_items（兑换手续费配置）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| fee_template_id | BIGINT UNSIGNED | FK, NOT NULL | 关联模板 |
| from_currency | VARCHAR(10) | NOT NULL | 源币种（USDT/USD/HKD/EUR/USDC/BTC） |
| to_currency | VARCHAR(10) | NOT NULL | 目标币种 |
| fee_rate | DECIMAL(10,6) | NOT NULL, DEFAULT 0 | 费率（如 0.003 = 0.3%） |
| min_fee | DECIMAL(28,8) | NOT NULL, DEFAULT 0 | 最低手续费金额 |
| min_fee_currency | VARCHAR(10) | NOT NULL | 最低手续费币种 |
| created_at | DATETIME(3) | | |
| updated_at | DATETIME(3) | | |

**UNIQUE(fee_template_id, from_currency, to_currency)**

---

### 4. fee_template_crypto_withdrawal_items（数币提现手续费配置）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| fee_template_id | BIGINT UNSIGNED | FK, NOT NULL | 关联模板 |
| currency | VARCHAR(10) | NOT NULL | 币种（USDT/USDC/BTC） |
| chain | VARCHAR(20) | NOT NULL | 链类型（ETH_ERC20/TRX_TRC20/SOL_Solana/BSC_BEP20） |
| fee_rate | DECIMAL(10,6) | NOT NULL, DEFAULT 0 | 费率 |
| fixed_fee | DECIMAL(28,8) | NOT NULL, DEFAULT 0 | 固定费用 |
| created_at | DATETIME(3) | | |
| updated_at | DATETIME(3) | | |

**UNIQUE(fee_template_id, currency, chain)**

---

### 5. fee_template_fiat_withdrawal_items（法币提现手续费配置）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| fee_template_id | BIGINT UNSIGNED | FK, NOT NULL | 关联模板 |
| currency | VARCHAR(10) | NOT NULL | 币种（USD/HKD/EUR） |
| transfer_type | VARCHAR(10) | NOT NULL | 转账类型（LOCAL/CHATS/TT） |
| fee_rate | DECIMAL(10,6) | NOT NULL, DEFAULT 0 | 费率 |
| fixed_fee | DECIMAL(28,8) | NOT NULL, DEFAULT 0 | 固定费用 |
| created_at | DATETIME(3) | | |
| updated_at | DATETIME(3) | | |

**UNIQUE(fee_template_id, currency, transfer_type)**

---

### 6. transaction_records（平台业务资金流水）

业务流水层，统一记录充值/提现/兑换/划转等业务订单摘要。不包含任何业务类型专有字段——具体业务数据在各 `*_orders` 表中。

**注意：** 钱包余额/冻结的每一次变动记录在 `wallet_ledger`，不在本表。一笔业务流水可对应多条钱包账本（如提现：FREEZE → DEDUCT_FROZEN 或 UNFREEZE）。

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| platform_order_id | VARCHAR(64) | UNIQUE, NOT NULL | 平台订单号 |
| merchant_id | BIGINT UNSIGNED | FK, NOT NULL | 商户 ID |
| type | VARCHAR(20) | NOT NULL | 交易类型 |
| sub_type | VARCHAR(30) | | 子类型 |
| amount | DECIMAL(28,8) | NOT NULL | 交易金额 |
| currency | VARCHAR(10) | NOT NULL | 币种 |
| platform_fee | DECIMAL(28,8) | NOT NULL, DEFAULT 0 | 平台手续费 |
| platform_fee_currency | VARCHAR(10) | | 平台手续费币种 |
| actual_amount | DECIMAL(28,8) | | 实际到账金额 |
| remark | VARCHAR(500) | | 备注 |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'PENDING' | 订单状态 |
| completed_at | DATETIME(3) | | 完成时间 |
| failed_reason | VARCHAR(500) | | 失败原因 |
| created_at | DATETIME(3) | | |
| updated_at | DATETIME(3) | | |

**type 枚举值：**
- `DEPOSIT` — 充值
- `WITHDRAWAL` — 提现
- `EXCHANGE` — 兑换
- `TRANSFER` — 划转

**sub_type 枚举值：**
- 充值：`CRYPTO_DEPOSIT`
- 提现：`CRYPTO_WITHDRAWAL`, `FIAT_WITHDRAWAL`
- 兑换：`SPOT_EXCHANGE`, `ONE_TO_ONE`
- 划转：`FUNDING_TO_TRADING`, `TRADING_TO_FUNDING`

**status 枚举值：**
- `PENDING` — 待处理
- `PENDING_REVIEW` — 待审核（提现）
- `PROCESSING` — 处理中
- `SUCCESS` — 成功
- `FAILED` — 失败
- `REJECTED` — 已拒绝（提现审核拒绝）
- `CANCELLED` — 已取消

**关联关系：** 每笔 transaction_record 在对应的业务订单表中有且仅有一条关联记录（1:1）

| type | 关联表 |
|------|--------|
| DEPOSIT | deposit_orders |
| WITHDRAWAL | withdrawal_orders |
| EXCHANGE | exchange_orders |
| TRANSFER | transfer_orders |

---

### 6b. withdrawal_orders（提现订单表）

承载提现专有的业务流程。与 `transaction_records` 一对一关联。

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| transaction_record_id | BIGINT UNSIGNED | FK → transaction_records, UNIQUE, NOT NULL | 关联流水 |
| merchant_id | BIGINT UNSIGNED | FK → merchants, NOT NULL | 商户 ID |
| withdrawal_type | VARCHAR(20) | NOT NULL | CRYPTO / FIAT |
| crypto_address_id | BIGINT UNSIGNED | FK → crypto_addresses | 关联白名单地址 |
| bank_account_id | BIGINT UNSIGNED | FK → bank_accounts | 关联银行账户 |
| to_address | VARCHAR(255) | | 目标地址/银行账号（冗余快照） |
| chain | VARCHAR(20) | | 链类型（数币） |
| tx_id | VARCHAR(128) | | 链上交易哈希（数币） |
| transfer_type | VARCHAR(10) | | 法币转账类型 LOCAL/CHATS/TT |
| purpose | VARCHAR(10) | | 法币提现用途 |
| review_status | VARCHAR(20) | NOT NULL, DEFAULT 'PENDING_REVIEW' | 审核状态 |
| reviewer_id | BIGINT UNSIGNED | | 审核人 ID |
| reviewer_type | VARCHAR(10) | | ADMIN / SYSTEM |
| reviewed_at | DATETIME(3) | | 审核时间 |
| review_remark | VARCHAR(500) | | 审核备注（拒绝原因） |
| kun_order_id | VARCHAR(128) | | KUN 订单号 |
| kun_request_no | VARCHAR(64) | UNIQUE | KUN 幂等键 |
| kun_fee | DECIMAL(28,8) | DEFAULT 0 | KUN 手续费 |
| kun_fee_currency | VARCHAR(10) | | KUN 手续费币种 |
| kun_submitted_at | DATETIME(3) | | 审核通过后提交 KUN 的时间 |
| failed_reason | VARCHAR(500) | | 失败原因 |
| completed_at | DATETIME(3) | | 完成时间 |
| created_at | DATETIME(3) | | |
| updated_at | DATETIME(3) | | |

**review_status 枚举值：**
- `PENDING_REVIEW` — 待审核
- `APPROVED` — 已通过（提交 KUN 中）
- `REJECTED` — 已拒绝

**提现完整状态流转：**
```
商户提交
  → withdrawal_orders.review_status = PENDING_REVIEW
  → transaction_records.status = PENDING_REVIEW

管理员通过（或自动审核通过）
  → withdrawal_orders.review_status = APPROVED, 记录审核人/时间
  → 提交 KUN API
  → transaction_records.status = PROCESSING

管理员拒绝
  → withdrawal_orders.review_status = REJECTED, 记录拒绝原因
  → transaction_records.status = REJECTED
  → 冻结金额解冻

KUN Webhook 回调
  → transaction_records.status = SUCCESS / FAILED
  → withdrawal_orders 记录 tx_id、kun_fee、completed_at
```

**查询模式：**
- 管理员审核队列：`SELECT * FROM withdrawal_orders WHERE review_status = 'PENDING_REVIEW'`
- 商户提现记录：`JOIN transaction_records ON id = transaction_record_id WHERE type = 'WITHDRAWAL'`
- 统一交易列表：直接查 `transaction_records`（所有类型）

---

### 6c. deposit_orders（充值订单表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| transaction_record_id | BIGINT UNSIGNED | FK → transaction_records, UNIQUE, NOT NULL | 关联流水 |
| merchant_id | BIGINT UNSIGNED | FK → merchants, NOT NULL | 商户 ID |
| currency | VARCHAR(10) | NOT NULL | 币种 |
| chain | VARCHAR(20) | NOT NULL | 链类型 |
| to_address | VARCHAR(255) | NOT NULL | 充值目标地址 |
| from_address | VARCHAR(255) | | 发送方地址（Webhook 回调填充） |
| tx_id | VARCHAR(128) | | 链上交易哈希 |
| kun_order_id | VARCHAR(128) | | KUN 订单号 |
| confirmations | INT UNSIGNED | | 链上确认数 |
| completed_at | DATETIME(3) | | |
| created_at | DATETIME(3) | | |
| updated_at | DATETIME(3) | | |

---

### 6d. exchange_orders（兑换订单表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| transaction_record_id | BIGINT UNSIGNED | FK → transaction_records, UNIQUE, NOT NULL | 关联流水 |
| merchant_id | BIGINT UNSIGNED | FK → merchants, NOT NULL | 商户 ID |
| exchange_type | VARCHAR(20) | NOT NULL | SPOT_EXCHANGE / ONE_TO_ONE |
| from_currency | VARCHAR(10) | NOT NULL | 源币种 |
| from_amount | DECIMAL(28,8) | NOT NULL | 源金额 |
| to_currency | VARCHAR(10) | NOT NULL | 目标币种 |
| to_amount | DECIMAL(28,8) | | 成交到账金额（回调后填充） |
| exchange_rate | DECIMAL(20,10) | | 成交汇率 |
| quote_id | VARCHAR(128) | | 询价 ID（报价锁定用） |
| auto_transfer | VARCHAR(3) | DEFAULT 'NO' | 1:1 交易是否自动划转（YES/NO） |
| kun_order_id | VARCHAR(128) | | KUN 订单号 |
| kun_request_no | VARCHAR(64) | UNIQUE | KUN 幂等键 |
| kun_fee | DECIMAL(28,8) | DEFAULT 0 | KUN 手续费 |
| kun_fee_currency | VARCHAR(10) | | KUN 手续费币种 |
| completed_at | DATETIME(3) | | |
| created_at | DATETIME(3) | | |
| updated_at | DATETIME(3) | | |

---

### 6e. transfer_orders（划转订单表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| transaction_record_id | BIGINT UNSIGNED | FK → transaction_records, UNIQUE, NOT NULL | 关联流水 |
| merchant_id | BIGINT UNSIGNED | FK → merchants, NOT NULL | 商户 ID |
| from_account_type | VARCHAR(10) | NOT NULL | 源账户类型（FUNDING / TRADING） |
| to_account_type | VARCHAR(10) | NOT NULL | 目标账户类型（FUNDING / TRADING） |
| kun_order_id | VARCHAR(128) | | KUN 订单号 |
| kun_request_no | VARCHAR(64) | UNIQUE | KUN 幂等键 |
| completed_at | DATETIME(3) | | |
| created_at | DATETIME(3) | | |
| updated_at | DATETIME(3) | | |

---

### 7. crypto_addresses（数币白名单地址表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| merchant_id | BIGINT UNSIGNED | FK, NOT NULL | 商户 ID |
| currency | VARCHAR(10) | NOT NULL | 币种 |
| chain | VARCHAR(20) | NOT NULL | 链类型 |
| address | VARCHAR(255) | NOT NULL | 区块链地址 |
| alias | VARCHAR(64) | NOT NULL | 地址别名 |
| kun_account_id | VARCHAR(128) | | KUN 返回的账户 ID |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'INIT' | 绑定状态 |
| created_at | DATETIME(3) | | |
| updated_at | DATETIME(3) | | |
| deleted_at | DATETIME(3) | NULL | |

**status 枚举值：** `INIT`, `PROCESS`, `PENDING`, `RECHARGE_VERIFY`, `SUCCESS`, `FAILED`, `INVALID`

---

### 8. bank_accounts（法币提现银行账户表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| merchant_id | BIGINT UNSIGNED | FK, NOT NULL | 商户 ID |
| kun_account_id | VARCHAR(128) | | KUN 返回的账户 ID |
| currency_list | VARCHAR(50) | NOT NULL | 支持的币种（逗号分隔：USD,HKD） |
| transfer_type | VARCHAR(10) | NOT NULL | 转账类型（LOCAL/CHATS/TT） |
| account_no | VARCHAR(128) | NOT NULL | 银行账户号（加密存储） |
| account_name | VARCHAR(128) | NOT NULL | 账户名 |
| bank_name | VARCHAR(128) | | 银行名称 |
| bank_code | VARCHAR(20) | | 银行代码（CHATS） |
| swift_code | VARCHAR(20) | | SWIFT Code（TT/CHATS） |
| payee_country_code | VARCHAR(10) | | 收款人国家代码 |
| payee_address | VARCHAR(255) | | 收款人地址 |
| middle_swift_code | VARCHAR(20) | | 中间行 SWIFT Code |
| area | VARCHAR(10) | NOT NULL | 国家代码 |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'ACTIVE' | 状态 |
| created_at | DATETIME(3) | | |
| updated_at | DATETIME(3) | | |
| deleted_at | DATETIME(3) | NULL | |

---

### 9. webhook_logs（Webhook 回调日志表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| event_id | VARCHAR(128) | NOT NULL | KUN 事件 ID |
| event_topic | VARCHAR(64) | NOT NULL | 事件主题 |
| event_type | VARCHAR(20) | NOT NULL | 事件类型（SUCCESS/FAIL/...） |
| customer_no | VARCHAR(64) | | 商户 KUN 号 |
| raw_data | JSON | NOT NULL | 原始回调 JSON 数据 |
| process_status | VARCHAR(20) | NOT NULL, DEFAULT 'PENDING' | 处理状态 |
| process_error | TEXT | | 处理错误信息 |
| processed_at | DATETIME(3) | | 处理完成时间 |
| received_at | DATETIME(3) | NOT NULL | 接收时间 |
| created_at | DATETIME(3) | | |

**process_status 枚举值：** `PENDING`, `SUCCESS`, `FAILED`, `DUPLICATE`

**UNIQUE(event_id, event_topic)** — 幂等去重

---

### 10. admin_users（管理员账号表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| username | VARCHAR(64) | UNIQUE, NOT NULL | 用户名 |
| email | VARCHAR(128) | UNIQUE, NOT NULL | 邮箱 |
| password_hash | VARCHAR(255) | NOT NULL | 密码哈希 |
| role | VARCHAR(20) | NOT NULL, DEFAULT 'OPERATOR' | 角色 |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'ACTIVE' | 状态 |
| last_login_at | DATETIME(3) | | 最后登录时间 |
| created_at | DATETIME(3) | | |
| updated_at | DATETIME(3) | | |
| deleted_at | DATETIME(3) | NULL | |

**role 枚举值：** `SUPER_ADMIN`, `OPERATOR`, `FINANCE`

---

### 11. audit_logs（审计日志表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| operator_id | BIGINT UNSIGNED | NOT NULL | 操作人 ID |
| operator_type | VARCHAR(10) | NOT NULL | 操作人类型（ADMIN/MERCHANT） |
| action | VARCHAR(64) | NOT NULL | 操作类型 |
| target_type | VARCHAR(32) | | 目标对象类型 |
| target_id | VARCHAR(64) | | 目标对象 ID |
| detail | JSON | | 操作详情 |
| ip_address | VARCHAR(45) | | IP 地址 |
| created_at | DATETIME(3) | | |

---

### 12. system_configs（系统配置表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| config_key | VARCHAR(64) | UNIQUE, NOT NULL | 配置键 |
| config_value | TEXT | NOT NULL | 配置值 |
| description | VARCHAR(255) | | 描述 |
| created_at | DATETIME(3) | | |
| updated_at | DATETIME(3) | | |

**币种相关配置键：**
- `supported_crypto_currencies` — 全局默认数币列表（CSV）
- `supported_fiat_currencies` — 全局默认法币列表（CSV）
- `supported_crypto_chains` — 全局默认支持链（JSON：`{"USDT":["ETH_ERC20","TRX_TRC20",...]}`）
- `default_crypto_chains` — 全局默认选中链（JSON：`{"USDT":"TRX_TRC20"}`）

---

### 13. merchant_wallets（商户钱包余额表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| merchant_id | BIGINT UNSIGNED | FK → merchants, NOT NULL | 商户 ID |
| account_type | VARCHAR(10) | NOT NULL | 账户类型 |
| currency | VARCHAR(10) | NOT NULL | 币种 |
| balance | DECIMAL(28,8) | NOT NULL, DEFAULT 0, CHECK ≥ 0 | 余额 |
| frozen_balance | DECIMAL(28,8) | NOT NULL, DEFAULT 0, CHECK ≥ 0 | 冻结金额（提现/兑换处理中） |
| version | BIGINT UNSIGNED | NOT NULL, DEFAULT 0 | 乐观锁版本号 |
| created_at | DATETIME(3) | | |
| updated_at | DATETIME(3) | | |

**account_type 枚举值：**
- `FUNDING` — 资金账户（充值/提现）
- `TRADING` — 交易账户（兑换）

**UNIQUE(merchant_id, account_type, currency)**

**CHECK 约束：** `balance >= frozen_balance`（可用余额 = balance - frozen_balance）

**业务说明：**
- 充值到 FUNDING 账户
- 提现从 FUNDING 账户扣款，处理中金额记入 frozen_balance
- 兑换前需先从 FUNDING 划转到 TRADING
- 余额变动使用乐观锁（version 字段）防并发
- **所有** balance / frozen_balance 变更必须经 `WalletService`，并在同一事务写入 `wallet_ledger`

---

### 13b. wallet_ledger（钱包资金变动账本）

只追加的钱包变动明细。用于对账、审计、按变动前后余额重建钱包状态。Migration：`000006_create_wallet_ledger.up.sql`。

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| merchant_id | BIGINT UNSIGNED | FK → merchants, NOT NULL | 商户 ID |
| wallet_id | BIGINT UNSIGNED | FK → merchant_wallets, NOT NULL | 钱包行 ID |
| account_type | VARCHAR(10) | NOT NULL | FUNDING / TRADING |
| currency | VARCHAR(10) | NOT NULL | 币种 |
| entry_type | VARCHAR(20) | NOT NULL | 变动类型 |
| amount | DECIMAL(28,8) | NOT NULL, CHECK > 0 | 变动金额（正数） |
| balance_before | DECIMAL(28,8) | NOT NULL | 变动前余额 |
| balance_after | DECIMAL(28,8) | NOT NULL | 变动后余额 |
| frozen_before | DECIMAL(28,8) | NOT NULL | 变动前冻结 |
| frozen_after | DECIMAL(28,8) | NOT NULL | 变动后冻结 |
| transaction_record_id | BIGINT UNSIGNED | FK → transaction_records | 关联业务流水 |
| biz_type | VARCHAR(20) | | DEPOSIT / WITHDRAWAL / EXCHANGE / TRANSFER |
| remark | VARCHAR(500) | | 备注（如 rollback） |
| created_at | DATETIME(3) | NOT NULL | 只追加，无 updated_at |

**entry_type 枚举值：**
- `CREDIT` — 入账（balance 增加）
- `FREEZE` — 冻结（frozen_balance 增加）
- `UNFREEZE` — 解冻（frozen_balance 减少）
- `DEDUCT_FROZEN` — 扣减已冻结（balance 与 frozen_balance 同时减少）

**写入约定：**
1. 先创建 `transaction_records`（及对应 `*_orders`），再变更钱包
2. 余额更新与账本插入必须在同一 DB 事务内
3. 禁止绕过 `WalletService` 直接改 `merchant_wallets`

**典型对应关系：**

| 业务动作 | wallet_ledger 行 |
|----------|------------------|
| 充值成功 | CREDIT × 1（FUNDING） |
| 提现提交 | FREEZE × 1 |
| 提现成功 | DEDUCT_FROZEN × 1 |
| 提现失败/拒绝 | UNFREEZE × 1 |
| 划转成功 | DEDUCT_FROZEN（源账户）+ CREDIT（目标账户） |
| 兑换成功 | DEDUCT_FROZEN（付出币种）+ CREDIT（得到币种） |

---

### 14. system_announcements（系统公告表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT UNSIGNED | PK | |
| title | VARCHAR(128) | NOT NULL | 公告标题 |
| content | TEXT | NOT NULL | 公告内容 |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'DRAFT' | 状态 |
| published_at | DATETIME(3) | | 发布时间 |
| created_by | BIGINT UNSIGNED | FK → admin_users, NOT NULL | 创建人 |
| created_at | DATETIME(3) | | |
| updated_at | DATETIME(3) | | |
| deleted_at | DATETIME(3) | NULL | |

**status 枚举值：**
- `DRAFT` — 草稿
- `PUBLISHED` — 已发布
