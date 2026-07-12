# Motewallet 钱包平台 — 产品需求文档（PRD）

## 概述

Motewallet 是一个面向商户的钱包管理平台，对接 KUN（鲲）KUN-SPACE 产品（波兰区域 KUN_PL）作为底层渠道。商户在平台完成注册、协议签署、KYC 认证后，可进行法币/数字货币的充值、兑换和提现操作。平台通过手续费模板实现灵活的费率管理（平台手续费独立于 KUN 手续费），并提供管理后台用于商户管理和系统配置。

## 背景与动机

- 商户需要一个统一的平台管理法币和数字货币的充值、兑换与提现
- KUN 提供底层金融基础设施（KUN-SPACE），我们作为上层平台封装业务逻辑
- **平台手续费与 KUN 手续费独立**：KUN 收取其渠道费用，平台在此基础上加收自己的手续费，因此系统必须自行记录所有资金流水，不能依赖 KUN 的记录
- 需要灵活的手续费机制以支持不同商户的差异化定价

## 目标用户

| 角色 | 描述 |
|------|------|
| 商户（Merchant） | 在平台注册并使用充值/兑换/提现服务的企业用户 |
| 平台管理员（Admin） | 管理商户、配置手续费模板、监控交易的内部人员 |

## 系统架构概览

```
商户（前端）──→ Motewallet 平台（后端）──→ KUN KUN-SPACE API（KUN_PL 波兰区域）
                    │                          │
管理员（后台）──→───┘                          │
                    │                          │
                MySQL 数据库 ←── KUN Webhook 回调
                (平台自有资金记录)
```

### 关键架构决策

1. **区域**：仅使用 KUN_PL（波兰），所有 `regionCode` 参数固定为 `KUN_PL`
2. **无平台审批**：提现直接提交给 KUN 处理，平台不做二次审批
3. **独立资金记录**：平台自行维护所有充值/兑换/提现/划转的交易记录，手续费按平台模板计算，与 KUN 的手续费分开记录
4. **Webhook 驱动**：通过 KUN Webhook 回调接收异步状态通知（充值到账、提现完成、兑换成交、KYC 结果等），辅以主动轮询兜底
5. **双账户体系**：每个商户有**资金账户**和**交易账户**，兑换前必须先将资金从资金账户划转到交易账户

---

## 模块一：商户注册与入网

### 业务流程

```
商户注册（平台账号）
  → 调用 KUN 注册子商户（获取 subCustomerNo）
  → 查询待签署协议
  → 签署协议
  → 提交 KYC 认证资料
  → 等待认证结果（Webhook 回调 CUSTOMER_AUTH 或主动查询）
  → 认证通过 → 账户激活
```

### 用户故事

#### US-101: 商户注册平台账号

**作为**商户，**我想要**在 Motewallet 平台注册账号，**以便**使用平台的充值、兑换和提现服务。

**验收标准：**
- [ ] AC-1: 商户可通过邮箱注册平台账号，邮箱全局唯一
- [ ] AC-2: 注册成功后，系统自动调用 KUN API `/rest/v2.0/customer/register` 创建子商户
- [ ] AC-3: 系统保存 KUN 返回的 `subCustomerNo` 与平台账号关联
- [ ] AC-4: 注册成功后自动绑定默认手续费模板
- [ ] AC-5: 注册成功后跳转到协议签署页面
- [ ] AC-6: 在平台数据库创建商户记录，初始状态为"待签署协议"

**优先级：** Must Have

**KUN API 对接：**
- `POST /rest/v2.0/customer/register` — 请求：`{ email, requestNo }` → 响应：`{ subCustomerNo }`

---

#### US-102: 签署协议

**作为**商户，**我想要**查看并签署平台协议，**以便**合规地使用平台服务。

**验收标准：**
- [ ] AC-1: 系统展示待签署协议列表（标题、版本、内容链接）
- [ ] AC-2: 商户可以查看每份协议的完整内容（支持下载）
- [ ] AC-3: 商户可以一键签署所有协议
- [ ] AC-4: 签署完成后状态变为 `COMPLETED`
- [ ] AC-5: 平台记录协议签署时间和协议版本
- [ ] AC-6: 签署完成后跳转到 KYC 认证页面

**优先级：** Must Have

**KUN API 对接：**
- `POST /rest/v2.0/customer/agreeList` — 查询待签署协议（`signStatus=UNSIGN`, `bizCode=KUN_SPACE_REGIST`）
- `POST /rest/v2.0/customer/agree/auth` — 签署协议（`protocolIds` 逗号分隔）
- `POST /rest/v2.0/customer/agree/download` — 下载协议文件

---

#### US-103: KYC 认证

**作为**商户，**我想要**提交企业认证资料，**以便**通过 KYC 审核后使用资金功能。

**验收标准：**
- [ ] AC-1: 提供分步表单收集企业信息（公司注册信息、股东信息、董事信息）
- [ ] AC-2: 支持上传营业执照、公司章程、身份证件等文件
- [ ] AC-3: 管理人支持人脸识别或手持证件照验证
- [ ] AC-4: 提交后展示"审核中"状态
- [ ] AC-5: 通过 Webhook 回调 `CUSTOMER_AUTH` 自动更新认证结果（SUCCESS/FAIL）
- [ ] AC-6: 辅以主动轮询查询结果（兜底，超过 1 小时未收到回调时触发）
- [ ] AC-7: 认证通过后账户自动激活，可使用充值/提现功能
- [ ] AC-8: 认证失败显示失败原因，支持重新提交

**优先级：** Must Have

**KUN API 对接：**
- `POST /rest/v2.0/customer/subMerchant/register` — 提交 KYC
- `POST /rest/v2.0/customer/merchant/register/query` — 主动查询认证结果
- Webhook: `eventTopic=CUSTOMER_AUTH`, `eventType=SUCCESS/FAIL`

---

## 模块二：充值（Deposit）

### 业务流程

```
数字货币充值：获取充值地址 → 链上转账 → KUN Webhook 回调到账 → 平台记录资金流水
```

> **注意：** 法币充值暂不支持，后续版本再开放。

### 用户故事

#### US-201: 数字货币充值

**作为**商户，**我想要**将数字货币充值到平台账户，**以便**进行兑换或站内转账。

**验收标准：**
- [ ] AC-1: 支持选择币种（USDT, USDC, BTC）
- [ ] AC-2: 支持选择链类型（ETH_ERC20, TRX_TRC20）
- [ ] AC-3: 显示充值地址和二维码，支持复制地址
- [ ] AC-4: 显示充值注意事项（最小金额、确认数、预估到账时间）
- [ ] AC-5: 通过 Webhook 回调 `DIGITAL_RECHARGE` 自动记录充值到账
- [ ] AC-6: 充值到账后更新平台资金账户余额记录
- [ ] AC-7: 平台自行创建充值流水记录（包含 KUN 订单号、链上 txId、金额、币种、链、状态）
- [ ] AC-8: 支持查看充值历史记录

**优先级：** Must Have

**KUN API 对接：**
- `POST /rest/v2.0/customer/crypto/deposit/addresses` — 获取充值地址
- Webhook: `eventTopic=DIGITAL_RECHARGE`, `eventType=SUCCESS`
  - 回调 data 包含：`orderId`, `txId`, `chain`, `orderAmount`, `orderCurrency`, `orderStatus`, `toAddress`, `fromAddress`, `feeAmount`, `feeCurrency`

---

#### ~~US-202: 法币充值~~ （暂不支持，后续版本开放）

---

## 模块三：兑换（Exchange）

### 业务流程

```
确认交易账户余额（必须先划转）
  → 选择兑换币对
  → 获取实时报价
  → 确认兑换（显示平台手续费）
  → 创建订单
  → Webhook 回调成交结果
  → 平台记录兑换流水（含平台手续费）
```

**重要业务规则：**
- 兑换前必须先将资金从**资金账户划转到交易账户**
- 兑换完成后资金留在交易账户，需划转回资金账户才能提现
- 平台手续费在兑换时额外扣除，与 KUN 的 `tradeFee` 分开计算

### 用户故事

#### US-301: 资金划转（资金账户 ↔ 交易账户）

**作为**商户，**我想要**在资金账户和交易账户之间划转资金，**以便**进行兑换交易。

**验收标准：**
- [ ] AC-1: 页面同时显示**资金账户**和**交易账户**的各币种余额
- [ ] AC-2: 支持选择划转方向：
  - 资金账户 → 交易账户（准备兑换）
  - 交易账户 → 资金账户（兑换后转回，以便提现）
- [ ] AC-3: 支持选择币种（法币：USD, HKD, EUR；数字货币：USDT, USDC, BTC）
- [ ] AC-4: 输入划转金额，不能超过源账户余额
- [ ] AC-5: 划转成功后实时更新两个账户余额显示
- [ ] AC-6: 平台自行记录划转流水
- [ ] AC-7: 支持查看划转历史记录

**优先级：** Must Have

**KUN API 对接：**
- `POST /rest/v2.0/user/fund/transfer` — 资金划转
  - `fromAcc` / `toAcc`: 固定为 `KUN_PL`（波兰区域内划转）
  - 注意：此接口是区域间划转，资金账户和交易账户的划转可能使用不同 API 或参数，需确认
- 划转记录查询

---

#### US-302: 实盘兑换（法币 ↔ 数字货币）

**作为**商户，**我想要**将法币兑换为数字货币（或反向兑换），**以便**满足业务需求。

**验收标准：**
- [ ] AC-1: 支持的兑换币对：USDT/USD, USDT/HKD, USDT/USDC, USDT/BTC 及其反向
- [ ] AC-2: 兑换前检查交易账户余额是否充足，不足时提示先划转
- [ ] AC-3: 输入金额后实时显示：
  - KUN 报价（汇率）
  - **平台手续费**（根据商户绑定的手续费模板计算）
  - 预估到账金额（扣除平台手续费后）
- [ ] AC-4: 确认兑换后创建订单
- [ ] AC-5: 通过 Webhook 回调接收成交结果
- [ ] AC-6: 平台自行记录兑换流水：
  - 订单号（平台 + KUN）
  - 兑换金额、成交价、成交金额
  - KUN 手续费（`tradeFee`）
  - **平台手续费**（按模板计算）
  - 实际到账金额
- [ ] AC-7: 支持查看兑换历史记录（含手续费明细）

**优先级：** Must Have

**KUN API 对接：**
- `POST /rest/v2.0/trade/exchange/marketInquiry` — 实时询价
- `POST /rest/v2.0/trade/exchange/order` — 创建兑换订单
- `POST /rest/v2.0/trade/exchange/order/query` — 主动查询（兜底）
- Webhook: `eventTopic` = `LEGAL_EXCHANGE_DIGITAL` / `DIGITAL_EXCHANGE_LEGAL` / `DIGITAL_EXCHANGE_BUY` / `DIGITAL_EXCHANGE_SELL`, `eventType=SUCCESS`

**报价计算规则：**
- 法币→USDT：到账金额 = amount / askPrice
- USDT→法币：到账金额 = amount × askPrice
- 稳定币互换 USDT→USDC：到账金额 = amount / askPrice
- USDC→USDT：到账金额 = amount × askPrice

---

#### US-303: 1:1 交易

**作为**商户，**我想要**进行 1:1 等额交易（如 USD ↔ USDT），**以便**快速完成稳定币和法币的互换。

**验收标准：**
- [ ] AC-1: 支持的交易对：USDT/USD, USD/USDT, USD/USDC, USDC/USD
- [ ] AC-2: 按 1:1 比率兑换，显示**平台手续费**
- [ ] AC-3: 可选是否自动划转到资金账户（autoTransfer=YES）
- [ ] AC-4: 通过 Webhook 回调接收成交结果
- [ ] AC-5: 平台自行记录交易流水（含平台手续费）

**优先级：** Should Have

**KUN API 对接：**
- `POST /rest/v2.0/trade/inner/match/create` — 创建 1:1 订单
- `POST /rest/v2.0/trade/inner/match/query` — 查询订单
- Webhook: `eventTopic` = `MAKER_DIGITAL_EXCHANGE_LEGAL` / `MAKER_LEGAL_EXCHANGE_DIGITAL`

---

## 模块四：提现（Withdrawal）

### 业务流程

```
数字货币提现：选择白名单地址 → 输入金额 → 计算平台手续费 → 确认提交
              → 【待审核】 → 管理员审核（通过/拒绝）
              → 通过 → 提交 KUN → Webhook 回调结果
              → 拒绝 → 资金解冻，通知商户

法币提现：选择银行账户 → 输入金额 → 计算平台手续费 → 确认提交
          → 【待审核】 → 管理员审核（通过/拒绝）
          → 通过 → 提交 KUN → Webhook 回调结果
          → 拒绝 → 资金解冻，通知商户
```

**重要业务规则：**
- **所有提现需平台审核**：商户提交提现后状态为"待审核"，管理员审核通过后才提交给 KUN 处理
- **自动审核**：支持配置自动审核阈值（如 ≤ 1000 USD 自动通过），超过阈值必须人工审核
- 提交提现时立即**冻结对应金额**（frozen_balance），审核拒绝后解冻
- 提现从**资金账户**扣款（交易账户的资金需先划转回资金账户）
- 平台手续费在提现时从提现金额中扣除
- 审核操作记录到审计日志，包含审核人、时间、备注

### 用户故事

#### US-401: 绑定数字货币提现地址（白名单）

**作为**商户，**我想要**绑定常用的数字货币提现地址，**以便**安全快捷地提现。

**验收标准：**
- [ ] AC-1: 支持添加白名单地址（USDT, USDC, BTC）
- [ ] AC-2: 支持选择链类型（ETH_ERC20, TRX_TRC20, SOL_Solana, BSC_BEP20）
- [ ] AC-3: 可设置地址别名，方便识别
- [ ] AC-4: 显示地址绑定状态（INIT → PROCESS → SUCCESS/FAILED）
- [ ] AC-5: 平台同步记录绑定的地址信息
- [ ] AC-6: 支持查看已绑定地址列表
- [ ] AC-7: 支持删除已绑定地址

**优先级：** Must Have

**KUN API 对接：**
- `POST /rest/v2.0/customer/crypto/address/add` — 绑定地址
- 白名单查询 / 删除接口

---

#### US-402: 数字货币提现

**作为**商户，**我想要**将数字货币提现到外部钱包地址，**以便**将资金转出平台。

**验收标准：**
- [ ] AC-1: 只能提现到已绑定的白名单地址
- [ ] AC-2: 支持选择币种和链类型
- [ ] AC-3: 输入金额后显示：
  - **平台手续费**（根据手续费模板：按币种+链类型匹配费率）
  - 实际到账金额 = 提现金额 - 平台手续费 - KUN 链上手续费
- [ ] AC-4: 提现金额不能超过资金账户可用余额
- [ ] AC-5: 确认提现前展示完整摘要（地址、金额、平台手续费、预估到账金额）
- [ ] AC-6: 确认提交后订单状态为 `PENDING_REVIEW`（待审核），提现金额立即冻结
- [ ] AC-7: 低于自动审核阈值的订单自动通过审核，直接提交 KUN
- [ ] AC-8: 超过阈值的订单需管理员在后台审核（通过/拒绝）
- [ ] AC-9: 审核通过 → 提交 KUN → 状态变为 `PROCESSING`
- [ ] AC-10: 审核拒绝 → 状态变为 `REJECTED` → 冻结金额解冻回可用余额
- [ ] AC-11: 平台自行记录提现流水（含平台手续费和 KUN 手续费、审核人、审核时间）
- [ ] AC-12: 通过 Webhook 回调 `DIGITAL_WITHDRAWAL` 更新提现状态
- [ ] AC-13: 商户端显示完整状态流转：待审核 → 处理中 → 成功/失败，或待审核 → 已拒绝
- [ ] AC-14: 支持查看提现历史（含手续费明细和审核信息）

**优先级：** Must Have

**KUN API 对接：**
- `POST /rest/v2.0/trade/crypto/withdrawal` — 发起提现（审核通过后调用，`regionCode=KUN_PL`）
- Webhook: `eventTopic=DIGITAL_WITHDRAWAL`, `eventType=SUCCESS`
  - 回调 data 包含：`orderId`, `txId`, `chain`, `orderAmount`, `orderStatus`, `feeAmount`, `feeCurrency`, `toAddress`

---

#### US-403: 绑定法币提现银行账户

**作为**商户，**我想要**绑定银行账户，**以便**进行法币提现。

**验收标准：**
- [ ] AC-1: 支持绑定银行账户（USD, HKD, EUR）
- [ ] AC-2: 支持选择转账类型（LOCAL 本地转账 / CHATS 港币即时结算 / TT 电汇）
- [ ] AC-3: 填写完整银行信息（账户号、账户名、SWIFT Code 等）
- [ ] AC-4: TT 方式需额外填写 SWIFT Code 和收款人地址
- [ ] AC-5: 平台同步记录绑定的银行账户信息
- [ ] AC-6: 支持查看已绑定账户列表
- [ ] AC-7: 支持解绑账户

**优先级：** Must Have

**KUN API 对接：**
- `POST /rest/v2.0/customer/fiat/address/add` — 绑定账户
- 提现账户列表查询 / 解绑接口

---

#### US-404: 法币提现

**作为**商户，**我想要**将法币提现到银行账户，**以便**将资金转出平台。

**验收标准：**
- [ ] AC-1: 只能提现到已绑定的银行账户
- [ ] AC-2: 支持选择法币币种（USD, HKD, EUR）
- [ ] AC-3: 输入金额后显示：
  - **平台手续费**（根据手续费模板：按币种+转账类型匹配费率）
  - 实际到账金额 = 提现金额 - 平台手续费
- [ ] AC-4: 需填写转账用途（OTHER/TRAD/INVS/GDDS）和备注
- [ ] AC-5: 确认提现前展示完整摘要
- [ ] AC-6: 确认提交后订单状态为 `PENDING_REVIEW`（待审核），提现金额立即冻结
- [ ] AC-7: 低于自动审核阈值的订单自动通过审核，直接提交 KUN
- [ ] AC-8: 超过阈值的订单需管理员在后台审核（通过/拒绝）
- [ ] AC-9: 审核通过 → 提交 KUN → 状态变为 `PROCESSING`
- [ ] AC-10: 审核拒绝 → 状态变为 `REJECTED` → 冻结金额解冻回可用余额
- [ ] AC-11: 平台自行记录提现流水（含平台手续费、审核人、审核时间）
- [ ] AC-12: 通过 Webhook 回调 `LEGAL_WITHDRAWAL` 更新提现状态
- [ ] AC-13: 商户端显示完整状态流转：待审核 → 处理中 → 成功/失败，或待审核 → 已拒绝
- [ ] AC-14: 支持查看提现历史（含手续费明细和审核信息）

**优先级：** Must Have

**KUN API 对接：**
- `POST /rest/v2.0/trade/fiat/withdrawal` — 创建提现订单（审核通过后调用，`regionCode` 相关参数）
- Webhook: `eventTopic=LEGAL_WITHDRAWAL`, `eventType=SUCCESS`
  - 回调 data 包含：`orderId`, `orderAmount`, `orderStatus`, `feeAmount`, `feeCurrency`, `bankAccountNo`, `bankAccountName`

---

#### US-405: 提现审核（管理后台）

**作为**管理员，**我想要**审核商户的提现请求，**以便**防范系统漏洞或异常操作导致的资金损失。

**验收标准：**
- [ ] AC-1: 管理后台提供"提现审核"独立页面，展示所有待审核的提现订单
- [ ] AC-2: 待审核列表显示：订单号、商户邮箱、提现类型（数币/法币）、金额+币种、目标地址/银行账户、平台手续费、提交时间
- [ ] AC-3: 支持按提现类型（数币/法币）、币种、金额范围、商户筛选
- [ ] AC-4: 审核操作包含两个选项：
  - **通过**：订单提交到 KUN 执行，状态 → `PROCESSING`
  - **拒绝**：必须填写拒绝原因，状态 → `REJECTED`，冻结金额立即解冻
- [ ] AC-5: 审核通过时，若 KUN 接口调用失败，状态回退为 `PENDING_REVIEW`，提示管理员稍后重试
- [ ] AC-6: 记录审核人（admin_user_id）、审核时间、审核备注到订单
- [ ] AC-7: 审核操作写入审计日志（audit_logs）
- [ ] AC-8: 支持批量审核（批量通过，不支持批量拒绝——拒绝需逐笔填写原因）
- [ ] AC-9: 待审核订单数量在侧边栏菜单项上显示红色角标（Badge）
- [ ] AC-10: 仅 `SUPER_ADMIN` 和 `OPERATOR` 角色可执行审核操作，`FINANCE` 角色仅可查看

**优先级：** Must Have

#### US-406: 提现自动审核规则配置

**作为**管理员，**我想要**配置自动审核规则，**以便**小额提现无需人工干预，同时大额提现保持人工把关。

**验收标准：**
- [ ] AC-1: 支持在系统配置中设置自动审核阈值（按等值 USD 计算）
- [ ] AC-2: 默认阈值为 0（即所有提现都需人工审核，最安全的初始值）
- [ ] AC-3: 可分别配置数币提现和法币提现的阈值
- [ ] AC-4: 低于阈值的提现自动通过审核，直接提交 KUN
- [ ] AC-5: 自动审核的订单审核人记录为 `SYSTEM`，审核备注记录为"Auto-approved: amount below threshold"
- [ ] AC-6: 阈值修改记录到审计日志

**优先级：** Must Have

---

## 模块五：账户与余额

### 双账户体系

每个认证通过的商户在 KUN 上拥有两类账户：

| 账户类型 | 用途 | 可执行操作 |
|---------|------|-----------|
| **资金账户** | 资金存放、充值入账、提现出金 | 充值、提现、划转到交易账户 |
| **交易账户** | 兑换交易 | 兑换、划转回资金账户 |

```
充值 → 【资金账户】 ←划转→ 【交易账户】 → 兑换
                ↓
              提现
```

### 用户故事

#### US-501: 查看账户余额

**作为**商户，**我想要**查看资金账户和交易账户的各币种余额，**以便**了解资金情况和决定操作。

**验收标准：**
- [ ] AC-1: 分区显示**资金账户**和**交易账户**的余额
- [ ] AC-2: 每个账户显示所有币种余额：
  - 法币：USD, HKD, EUR
  - 数字货币：USDT, USDC, BTC
- [ ] AC-3: 显示可用余额
- [ ] AC-4: 余额数据支持手动刷新
- [ ] AC-5: 提供资产总览（折合 USD 总值，资金账户 + 交易账户合计）
- [ ] AC-6: 在资金账户区域提供快捷入口：数币充值、提现、划转到交易账户
- [ ] AC-7: 在交易账户区域提供快捷入口：去兑换、划转回资金账户

**优先级：** Must Have

**KUN API 对接：**
- `POST /rest/v2.0/trade/account/outAccount/query` — 查询余额（`currency`, `currencyType`）
- `POST /rest/v2.0/account/query/balance` — 查询区域余额（`regionCode=KUN_PL`）

---

## 模块六：平台资金记录

### 核心设计原则

**平台必须自行维护完整的资金流水记录**，原因：
1. 平台手续费与 KUN 手续费不同，需要独立记录
2. 平台需要自己的交易统计和对账能力
3. 商户看到的手续费是平台手续费，而非 KUN 手续费

### 用户故事

#### US-601: 商户交易记录

**作为**商户，**我想要**查看我的所有交易记录，**以便**追踪资金流向。

**验收标准：**
- [ ] AC-1: 统一的交易记录列表，包含所有类型：
  - 充值（数币/法币）
  - 兑换（实盘/1:1）
  - 提现（数币/法币）
  - 划转（资金账户 ↔ 交易账户）
- [ ] AC-2: 支持按类型、状态、时间范围筛选
- [ ] AC-3: 每条记录显示：
  - 交易类型、时间、币种、金额
  - **平台手续费**（非 KUN 手续费）
  - 实际到账/扣款金额
  - 状态（处理中/成功/失败）
  - 平台订单号 + KUN 订单号
- [ ] AC-4: 支持查看交易详情
- [ ] AC-5: 支持分页加载

**优先级：** Must Have

---

#### US-602: 资金流水对账

**作为**平台系统，**需要**自行维护资金流水，**以便**与 KUN 对账和计算平台收益。

**验收标准：**
- [ ] AC-1: 每笔交易在平台数据库创建流水记录，包含：
  - `platform_order_id` — 平台订单号
  - `kun_order_id` — KUN 订单号
  - `merchant_id` — 商户 ID
  - `type` — 交易类型（DEPOSIT/EXCHANGE/WITHDRAWAL/TRANSFER）
  - `amount` — 交易金额
  - `currency` — 币种
  - `platform_fee` — 平台手续费
  - `kun_fee` — KUN 手续费
  - `actual_amount` — 实际到账/扣款金额
  - `status` — 状态
  - `created_at` / `updated_at`
- [ ] AC-2: Webhook 回调时同步更新流水状态
- [ ] AC-3: 平台手续费按商户绑定的模板实时计算
- [ ] AC-4: 支持后台导出对账报表
- [ ] AC-5: 每次钱包余额/冻结变动写入 `wallet_ledger`（含变动前后余额、关联 `transaction_record_id`），与业务流水同事务提交

**优先级：** Must Have

---

## 模块七：手续费模板

### 业务规则

- 平台维护一组**手续费模板**，每个模板定义兑换和提现的**平台手续费**费率
- **平台手续费独立于 KUN 手续费**：KUN 收取渠道费用，平台在此基础上另外收取
- 系统有一个**默认模板**，商户注册时自动绑定
- 管理员可以为不同商户指定不同模板
- 手续费在**兑换**和**提现**两个环节收取（充值不收手续费）
- 手续费计算模式：**max(固定金额, 交易金额 × 费率)**，取两者较大值

### 手续费模板数据结构

```
FeeTemplate（手续费模板）
├── name: 模板名称
├── description: 描述
├── is_default: 是否默认模板
├── exchange_fees[]（兑换手续费配置）
│   ├── from_currency: 源币种
│   ├── to_currency: 目标币种
│   ├── fee_rate: 费率百分比（如 0.5%）
│   └── min_fee: 最低手续费金额
├── crypto_withdrawal_fees[]（数币提现手续费配置）
│   ├── currency: 币种（USDT/USDC/BTC）
│   ├── chain: 链类型（ETH_ERC20/TRX_TRC20/...）
│   ├── fee_rate: 费率百分比
│   └── fixed_fee: 固定手续费
└── fiat_withdrawal_fees[]（法币提现手续费配置）
    ├── currency: 币种（USD/HKD/EUR）
    ├── transfer_type: 转账类型（LOCAL/CHATS/TT）
    ├── fee_rate: 费率百分比
    └── fixed_fee: 固定手续费
```

### 用户故事

#### US-701: 手续费模板管理（管理后台）

**作为**管理员，**我想要**创建和管理手续费模板，**以便**灵活控制平台收费。

**验收标准：**
- [ ] AC-1: 支持创建手续费模板（名称 + 描述）
- [ ] AC-2: 每个模板包含三类费率配置：
  - **兑换手续费**：按币对（如 USDT→USD）设置费率和最低手续费
  - **数币提现手续费**：按币种+链类型设置费率和固定费用
  - **法币提现手续费**：按币种+转账类型设置费率和固定费用
- [ ] AC-3: 手续费计算规则：`max(固定金额, 交易金额 × 费率)`
- [ ] AC-4: 支持设置一个模板为"默认模板"（有且仅有一个）
- [ ] AC-5: 支持编辑模板（修改即时生效）
- [ ] AC-6: 支持删除模板（已绑定商户的模板不可删除，需先解绑）
- [ ] AC-7: 支持查看哪些商户绑定了该模板

**优先级：** Must Have

---

#### US-702: 商户绑定手续费模板（管理后台）

**作为**管理员，**我想要**为商户指定手续费模板，**以便**实现差异化定价。

**验收标准：**
- [ ] AC-1: 商户注册时自动绑定默认模板
- [ ] AC-2: 管理员可修改商户绑定的模板
- [ ] AC-3: 修改即时生效，下一笔交易使用新模板计算手续费
- [ ] AC-4: 支持批量修改多个商户的模板
- [ ] AC-5: 商户详情页显示当前绑定的模板名称

**优先级：** Must Have

---

#### US-703: 手续费展示（商户端）

**作为**商户，**我想要**在兑换和提现时看到平台手续费明细，**以便**了解实际费用。

**验收标准：**
- [ ] AC-1: 兑换页面显示：兑换金额、**平台手续费**、实际到账金额
- [ ] AC-2: 提现页面显示：提现金额、**平台手续费**、实际到账金额
- [ ] AC-3: 手续费计算准确，精度到最小货币单位
- [ ] AC-4: 支持查看当前适用的费率信息

**优先级：** Must Have

---

## 模块八：管理后台

### 用户故事

#### US-801: 商户管理

**作为**管理员，**我想要**管理平台上的所有商户，**以便**监控和维护商户状态。

**验收标准：**
- [ ] AC-1: 商户列表（分页、搜索、筛选）
  - 展示：商户 ID、邮箱、KUN 子商户号、入网状态、手续费模板、注册时间
  - 筛选：入网状态（待签署协议/待认证/审核中/已认证/认证失败）、手续费模板
- [ ] AC-2: 商户详情页
  - 基本信息、入网状态、绑定的手续费模板
  - 资金账户 + 交易账户余额概览
  - 操作记录（充值/兑换/提现/划转历史）
- [ ] AC-3: 支持冻结/解冻商户（冻结后无法执行任何交易）
- [ ] AC-4: 支持修改商户的手续费模板

**优先级：** Must Have

---

#### US-802: 交易监控

**作为**管理员，**我想要**查看平台所有交易记录，**以便**监控资金流动和平台收益。

**验收标准：**
- [ ] AC-1: 交易记录列表（分页、搜索、筛选）
  - 类型：充值/兑换/提现/划转
  - 状态：处理中/成功/失败
  - 时间范围筛选
  - 商户筛选
- [ ] AC-2: 交易详情：
  - 交易金额、**平台手续费**、**KUN 手续费**
  - 平台订单号、KUN 订单号
  - 状态、时间
- [ ] AC-3: 支持导出交易记录（CSV）
- [ ] AC-4: 平台收益统计（按日/周/月汇总平台手续费收入）

**优先级：** Must Have

---

#### US-803: 系统配置

**作为**管理员，**我想要**配置系统参数，**以便**控制平台行为。

**验收标准：**
- [ ] AC-1: 支持配置支持的币种和链类型（开关控制）
- [ ] AC-2: 支持配置 KUN API 连接参数（API Key、密钥、回调地址）
- [ ] AC-3: 支持查看和测试 Webhook 连接状态
- [ ] AC-4: 支持配置默认手续费模板
- [ ] AC-5: 系统公告管理（向商户端推送公告）

**优先级：** Should Have

---

#### US-804: 管理员账户管理

**作为**超级管理员，**我想要**管理后台用户权限，**以便**控制操作权限。

**验收标准：**
- [ ] AC-1: 支持创建管理员账户
- [ ] AC-2: 支持角色划分（超级管理员、运营、财务）
- [ ] AC-3: 不同角色有不同的菜单和操作权限
- [ ] AC-4: 操作日志记录（谁在什么时间做了什么操作）

**优先级：** Should Have

---

## 模块九：KUN Webhook 回调处理

### Webhook 接入规范

#### 回调接收

- 接收方式：HTTP POST，JSON 格式
- 回调格式：
```json
{
  "eventId": "事件唯一 ID",
  "eventType": "SUCCESS / FAIL / PENDING / ACTIVE",
  "eventTopic": "事件主题",
  "data": { ... }
}
```

#### 签名验证

- 算法：SHA256withRSA（RSA2048）
- 签名编码：Base64
- 验签步骤：
  1. 提取 `data` 字段
  2. 添加 `timestamp` 字段
  3. 所有 key 转小写，按 ASCII 升序排序
  4. 拼接为 `key=value` 用 `&` 连接
  5. 使用 KUN 公钥验签
- 时间窗口：timestamp 与服务器时间差 ≤ 5 分钟

#### 响应格式

- 成功：HTTP 200 + `{ "code": "200" }`
- 失败：HTTP 非 200

#### 重试机制

- KUN 每 2 分钟重试一次，最多 10 次
- 平台必须支持幂等处理（同一通知可能重复发送）

### 需要处理的事件

| eventTopic | 事件 | 处理逻辑 |
|------------|------|---------|
| `CUSTOMER_AUTH` | KYC 认证结果 | 更新商户认证状态 |
| `DIGITAL_RECHARGE` | 数币充值到账 | 创建充值流水，更新余额记录 |
| ~~`PAYX.LEGAL_RECHARGE`~~ | ~~法币充值到账~~ | ~~暂不支持~~ |
| `DIGITAL_WITHDRAWAL` | 数币提现完成 | 更新提现流水状态 |
| `LEGAL_WITHDRAWAL` | 法币提现完成 | 更新提现流水状态 |
| `DIGITAL_EXCHANGE_BUY` | 数币兑换（USDT→其他） | 更新兑换流水状态 |
| `DIGITAL_EXCHANGE_SELL` | 数币兑换（其他→USDT） | 更新兑换流水状态 |
| `LEGAL_EXCHANGE_DIGITAL` | 法币→数币兑换 | 更新兑换流水状态 |
| `DIGITAL_EXCHANGE_LEGAL` | 数币→法币兑换 | 更新兑换流水状态 |
| `MAKER_DIGITAL_EXCHANGE_LEGAL` | 1:1 交易（数币→法币） | 更新交易流水状态 |
| `MAKER_LEGAL_EXCHANGE_DIGITAL` | 1:1 交易（法币→数币） | 更新交易流水状态 |

### 用户故事

#### US-901: Webhook 回调处理

**作为**系统，**需要**正确处理 KUN Webhook 回调，**以便**实时更新订单状态。

**验收标准：**
- [ ] AC-1: 提供 Webhook 接收端点 `POST /api/v1/webhook/kun`
- [ ] AC-2: 验证签名（SHA256withRSA + KUN 公钥）
- [ ] AC-3: 验证 timestamp 时间窗口（≤ 5 分钟）
- [ ] AC-4: 幂等处理：同一 `eventId` 重复回调不重复处理
- [ ] AC-5: 根据 `eventTopic` 路由到不同的处理逻辑
- [ ] AC-6: 处理成功返回 HTTP 200 + `{"code":"200"}`
- [ ] AC-7: 记录所有回调日志（含原始数据，用于排查和对账）

**优先级：** Must Have

---

#### US-902: 主动轮询兜底

**作为**系统，**需要**对未收到 Webhook 回调的订单进行主动查询，**以便**确保状态最终一致。

**验收标准：**
- [ ] AC-1: 定时任务扫描状态为"处理中"且超过一定时间的订单
- [ ] AC-2: 调用 KUN 对应的查询接口获取最新状态
- [ ] AC-3: 轮询间隔：首次 5 分钟，之后每 30 分钟一次
- [ ] AC-4: 最长轮询周期：24 小时后标记为异常
- [ ] AC-5: 避免频繁调用，防止触发 KUN 限流

**优先级：** Must Have

---

## 非功能需求

### 性能
- 页面加载时间 < 2 秒
- API 响应时间 < 500ms（不含 KUN 调用）
- Webhook 处理时间 < 1 秒（快速响应，异步处理业务逻辑）
- 支持 100+ 并发商户操作

### 安全
- 所有 API 需 JWT 认证（Webhook 端点除外，使用签名验证）
- 提现操作需二次验证（邮件验证码）
- 敏感数据加密存储（银行账号、KUN 密钥）
- SQL 注入防护（参数化查询）
- HTTPS 强制
- Webhook 签名验证（SHA256withRSA）
- 操作审计日志

### 可靠性
- KUN API 调用失败时重试（指数退避，最多 3 次）
- 关键操作保证幂等性（通过 `requestNo` + 平台订单号）
- Webhook 幂等处理（通过 `eventId` 去重）
- 数据库事务保证资金操作一致性
- 主动轮询兜底确保状态最终一致

### 兼容性
- 浏览器：Chrome 90+、Safari 15+、Firefox 90+、Edge 90+
- 移动端响应式适配

---

## 数据需求

### 核心实体

| 实体 | 描述 |
|------|------|
| Merchant | 商户账号，关联 KUN subCustomerNo，关联手续费模板 |
| MerchantKyc | 商户 KYC 认证信息和状态 |
| FeeTemplate | 手续费模板（名称、描述、是否默认） |
| FeeTemplateExchange | 兑换手续费配置（币对、费率、最低手续费） |
| FeeTemplateCryptoWithdrawal | 数币提现手续费配置（币种、链、费率、固定费用） |
| FeeTemplateFiatWithdrawal | 法币提现手续费配置（币种、转账类型、费率、固定费用） |
| TransactionRecord | **平台业务资金流水**（充值/兑换/提现/划转的统一业务记录，含平台手续费） |
| WalletLedger | **钱包资金变动账本**（每次余额/冻结变更一行，含变动前后余额，关联 TransactionRecord） |
| CryptoAddress | 数币白名单地址（平台记录） |
| BankAccount | 法币提现银行账户（平台记录） |
| WebhookLog | Webhook 回调日志（原始数据存档） |
| AdminUser | 管理员账号 |
| AuditLog | 操作审计日志 |
| SystemConfig | 系统配置 |

---

## 页面清单

### 商户端

| 页面 | 路由 | 描述 |
|------|------|------|
| 登录 | /login | 邮箱+密码登录 |
| 注册 | /register | 邮箱注册 |
| 协议签署 | /onboarding/agreement | 查看并签署协议 |
| KYC 认证 | /onboarding/kyc | 分步提交认证资料 |
| 认证状态 | /onboarding/status | 查看认证进度 |
| 资产总览 | /dashboard | 资金账户+交易账户余额、最近交易 |
| 数币充值 | /deposit/crypto | 获取充值地址 |
| 资金划转 | /transfer | 资金账户 ↔ 交易账户 |
| 兑换 | /exchange | 询价、下单（需先划转到交易账户） |
| 数币提现 | /withdraw/crypto | 选择白名单地址、输入金额 |
| 法币提现 | /withdraw/fiat | 选择银行账户、输入金额 |
| 地址管理 | /settings/addresses | 管理提现白名单地址 |
| 银行账户管理 | /settings/bank-accounts | 管理法币提现账户 |
| 交易记录 | /transactions | 所有交易历史（含手续费明细） |

### 管理后台

| 页面 | 路由 | 描述 |
|------|------|------|
| 管理员登录 | /admin/login | 管理员登录 |
| 仪表盘 | /admin/dashboard | 平台数据概览、手续费收入统计 |
| 商户列表 | /admin/merchants | 商户管理 |
| 商户详情 | /admin/merchants/:id | 商户详细信息 |
| 手续费模板 | /admin/fee-templates | 模板管理 |
| 模板详情 | /admin/fee-templates/:id | 模板编辑 |
| 提现审核 | /admin/withdrawal-reviews | 待审核提现列表、审核操作 |
| 交易记录 | /admin/transactions | 全局交易监控（含双手续费明细） |
| 系统配置 | /admin/settings | 系统参数配置 |
| 管理员管理 | /admin/users | 管理员账户管理 |
| 操作日志 | /admin/audit-logs | 审计日志查看 |
| Webhook 日志 | /admin/webhook-logs | Webhook 回调日志查看 |

---

## API 端点清单（Motewallet 平台自有 API）

### 认证
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/auth/register | 商户注册 |
| POST | /api/v1/auth/login | 商户登录 |
| POST | /api/v1/auth/refresh | 刷新 Token |

### 入网
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/v1/onboarding/agreements | 查询待签署协议 |
| POST | /api/v1/onboarding/agreements/sign | 签署协议 |
| POST | /api/v1/onboarding/kyc | 提交 KYC |
| GET | /api/v1/onboarding/kyc/status | 查询 KYC 状态 |

### 账户
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/v1/account/balances | 查询所有余额（资金账户+交易账户） |
| GET | /api/v1/account/ledger | 查询钱包资金变化（wallet_ledger） |
| POST | /api/v1/account/transfer | 资金划转（资金账户↔交易账户） |
| GET | /api/v1/account/transfers | 划转记录 |

### 充值
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/v1/deposit/crypto/addresses | 获取数币充值地址 |
| GET | /api/v1/deposit/orders | 充值记录 |

### 兑换
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/exchange/quote | 获取报价（含平台手续费计算） |
| POST | /api/v1/exchange/order | 创建兑换订单 |
| GET | /api/v1/exchange/orders | 兑换记录 |
| GET | /api/v1/exchange/orders/:id | 兑换详情 |
| POST | /api/v1/exchange/1to1 | 1:1 交易 |

### 提现
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/withdraw/crypto | 发起数币提现 |
| POST | /api/v1/withdraw/fiat | 发起法币提现 |
| GET | /api/v1/withdraw/orders | 提现记录 |
| GET | /api/v1/withdraw/orders/:id | 提现详情 |

### 地址/账户管理
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/addresses/crypto | 添加数币白名单地址 |
| GET | /api/v1/addresses/crypto | 查询白名单地址 |
| DELETE | /api/v1/addresses/crypto/:id | 删除白名单地址 |
| POST | /api/v1/addresses/bank | 绑定银行账户 |
| GET | /api/v1/addresses/bank | 查询银行账户 |
| DELETE | /api/v1/addresses/bank/:id | 解绑银行账户 |

### 交易记录
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/v1/transactions | 商户交易记录（含平台手续费明细） |
| GET | /api/v1/transactions/:id | 交易详情 |

### Webhook
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/webhook/kun | KUN Webhook 回调接收端点 |

### 管理后台
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/admin/auth/login | 管理员登录 |
| GET | /api/v1/admin/merchants | 商户列表 |
| GET | /api/v1/admin/merchants/:id | 商户详情 |
| PUT | /api/v1/admin/merchants/:id/status | 冻结/解冻商户 |
| PUT | /api/v1/admin/merchants/:id/fee-template | 修改商户模板 |
| POST | /api/v1/admin/fee-templates | 创建模板 |
| GET | /api/v1/admin/fee-templates | 模板列表 |
| GET | /api/v1/admin/fee-templates/:id | 模板详情 |
| PUT | /api/v1/admin/fee-templates/:id | 更新模板 |
| DELETE | /api/v1/admin/fee-templates/:id | 删除模板 |
| GET | /api/v1/admin/withdrawals/pending | 待审核提现列表 |
| POST | /api/v1/admin/withdrawals/:id/approve | 审核通过提现 |
| POST | /api/v1/admin/withdrawals/:id/reject | 审核拒绝提现 |
| POST | /api/v1/admin/withdrawals/batch-approve | 批量审核通过 |
| GET | /api/v1/admin/transactions | 全局交易记录 |
| GET | /api/v1/admin/transactions/stats | 手续费收入统计 |
| GET | /api/v1/admin/audit-logs | 审计日志 |
| GET | /api/v1/admin/webhook-logs | Webhook 日志 |
| GET/PUT | /api/v1/admin/settings | 系统配置 |

---

## 已确认事项（原开放问题）

| # | 问题 | 决策 |
|---|------|------|
| 1 | KUN 回调/Webhook | ✅ KUN 支持 Webhook 回调，平台通过 `POST /api/v1/webhook/kun` 接收，辅以主动轮询兜底 |
| 2 | 区域 | ✅ 仅支持波兰（KUN_PL），所有 regionCode 固定为 KUN_PL |
| 3 | 提现审批 | ✅ 所有提现需平台审核；支持配置自动审核阈值（默认 0，全部人工）；审核通过后提交 KUN |
| 4 | 资金记录 | ✅ 平台自行维护所有资金流水，手续费按平台模板计算，与 KUN 手续费独立 |
| 5 | 双账户 | ✅ 资金账户 + 交易账户，兑换前必须先划转到交易账户 |

## 已全部确认

| # | 问题 | 决策 |
|---|------|------|
| 6 | 资金划转 API | ✅ 使用 KUN 区域划转 API `/rest/v2.0/user/fund/transfer` |
| 7 | 支持的法币 | ✅ 波兰区域支持 USD, HKD, EUR |
| 8 | 邮件通知 | ✅ 暂不需要 |
| 9 | 法币充值 | ✅ 暂不支持，仅支持数字货币充值 |
