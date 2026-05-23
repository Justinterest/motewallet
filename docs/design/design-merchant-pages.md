# 商户端页面设计文档

## 全局布局

### 导航栏（TopNav）
- 固定顶部，白色背景，底部阴影
- 左侧：Motewallet Logo
- 中间：菜单项（Dashboard / Deposit / Transfer / Exchange / Withdraw / Transactions / Settings）
- 右侧：商户邮箱 + 头像下拉菜单（Settings / Logout）
- 当前页面菜单项高亮（蓝色下划线 + 蓝色文字）

### 页面容器
- 最大宽度：1280px，居中
- 左右 padding：32px
- 背景色：`#F8FAFC`

---

## P-01: 登录页 `/login`

### 布局
- 全屏居中卡片，最大宽度 420px
- 顶部：Motewallet Logo + "Sign in to your account"

### 表单
- Email 输入框
- Password 输入框（带显示/隐藏切换）
- "Remember me" 复选框 + "Forgot password?" 链接
- "Sign In" 主按钮（全宽）
- 底部："Don't have an account? Register"

### 状态
- 默认态：空表单
- 加载态：按钮 loading spinner
- 错误态：输入框红色边框 + 顶部错误提示 banner

---

## P-02: 注册页 `/register`

### 布局
- 全屏居中卡片，最大宽度 420px

### 表单
- Email 输入框
- Password 输入框（含强度指示器）
- Confirm Password 输入框
- "I agree to the Terms of Service" 复选框
- "Create Account" 主按钮
- 底部："Already have an account? Sign In"

---

## P-03: 协议签署页 `/onboarding/agreement`

### 布局
- 步骤指示器：① 签署协议（当前） → ② KYC 认证 → ③ 完成
- 协议列表卡片

### 内容
- 每个协议一行：
  - 协议标题 + 版本号
  - "查看详情" 链接（弹出模态框显示协议内容）
  - 下载按钮（PDF）
  - 签署状态：✓ 已签署（绿色） / ○ 待签署
- 底部："一键签署所有协议" 主按钮
- 签署成功后自动跳转 KYC 页面

---

## P-04: KYC 认证页 `/onboarding/kyc`

### 布局
- 步骤指示器：① 签署协议 ✓ → ② KYC 认证（当前） → ③ 完成
- 分步表单（Step Wizard），3 个子步骤：

### Step 1: 企业信息
- 公司英文名称（必须与证照一致）
- 公司中文名称
- 注册地区（国家选择器）
- 注册地址
- 注册号码
- 营业执照上传（拖拽上传区域）
- 公司类型选择（Private / Public）
- 行业选择（一级 + 二级联动下拉）
- 成立日期
- 业务区域（多选）
- 资金来源（单选）
- 开户目的（单选）
- 公司章程上传
- 其他认证材料上传

### Step 2: 管理人信息
- 姓名（中文 + 英文）
- 国籍
- 出生日期
- 性别
- 证件类型 + 证件号
- 证件有效期
- 居住国家 + 地址
- 验证方式选择（人脸识别 / 手持证件照）
- 证件照上传（至少 3 张：证件正面、自拍、手持证件）

### Step 3: 股东与董事
- 股东信息（动态添加，每个股东一个卡片）：
  - 姓名、国籍、证件号、持股比例、证件照
- 董事信息（动态添加）：
  - 姓名、国籍、证件号、证件照
- "提交认证" 按钮

### 交互
- 支持保存草稿（自动保存到 localStorage）
- 文件上传显示进度条
- 每步表单验证后才能进入下一步

---

## P-05: 认证状态页 `/onboarding/status`

### 三种状态展示

**审核中（AUTHING）：**
- 蓝色圆形进度图标 + "Verification in Progress"
- "Your KYC application is being reviewed. This usually takes 1-3 business days."
- 提交时间显示

**认证成功（AUTH_SUC）：**
- 绿色勾选图标 + "Verification Approved"
- "Your account is now active. You can start using all features."
- "Go to Dashboard" 按钮

**认证失败（AUTH_FAIL）：**
- 红色叉号图标 + "Verification Failed"
- 失败原因展示（来自 KUN `failReason`）
- "Resubmit Application" 按钮

---

## P-06: 资产总览 `/dashboard`

### Stitch 设计稿
- 已生成，Screen ID: `cbac0bfbc50e43cb8f85330182b4271f`

### 布局
1. **资产总览卡片**
   - 总资产（折合 USD）：大号金额
   - 环形图显示资产分布

2. **双账户余额卡片**（并排两列）

   **资金账户卡片：**
   - 标题 + 钱包图标
   - 各币种余额列表（USDT, USDC, BTC, USD, HKD, EUR）
   - 每行：币种图标 + 名称 + 余额（右对齐）
   - 快捷按钮：充值 / 提现 / 划转到交易账户

   **交易账户卡片：**
   - 标题 + 图表图标
   - 各币种余额列表
   - 快捷按钮：去兑换 / 划转回资金账户

3. **最近交易表格**
   - 列：类型（图标+文字）、金额、币种、平台手续费、状态（徽章）、时间
   - "查看全部" 链接

---

## P-07: 数币充值页 `/deposit/crypto`

### 布局
1. 页面标题："Crypto Deposit"

2. **币种和链选择卡片**
   - 币种选择：USDT / USDC / BTC（Tab 切换）
   - 链选择：ETH_ERC20 / TRX_TRC20（下拉或 Tab）

3. **充值地址卡片**
   - 二维码（居中，200x200px）
   - 地址文本（等宽字体，全显示）
   - "复制地址" 按钮

4. **注意事项卡片**（amber 背景）
   - 最小充值金额
   - 链上确认数要求
   - 请勿充值其他币种到此地址

5. **充值记录**
   - 列：金额、币种、链、TxID（截断+链接）、状态、时间

---

## P-08: 资金划转页 `/transfer`

### 布局
1. 页面标题："Fund Transfer"

2. **划转表单卡片**（居中，max-width 480px）
   - 源账户选择器：
     - "Funding Account" ←→ "Trading Account"
     - 中间有交换方向按钮（↔）
   - 源账户余额显示
   - 币种选择下拉（USD, HKD, EUR, USDT, USDC, BTC）
   - 金额输入 + "MAX" 按钮
   - "Confirm Transfer" 主按钮

3. **划转记录**
   - 列：方向（From→To）、金额、币种、状态、时间

---

## P-09: 兑换页 `/exchange`

### Stitch 设计稿
- 已生成，Screen ID: `a8c9d14d27ef443d8562154ed2c80264`

### 布局
1. **交易账户余额不足提醒**（顶部 amber 横幅）
   - "交易账户余额不足，请先划转资金" + 链接跳转 /transfer

2. **兑换表单卡片**（居中，max-width 560px）
   - From 区域：币种下拉 + 金额输入 + 可用余额 + MAX
   - 交换方向按钮（↕）
   - To 区域：币种下拉 + 预估到账金额（灰色只读）

3. **报价信息卡片**
   - 当前汇率：1 USDT = 0.9998 USD
   - 报价有效倒计时：29s

4. **费用摘要卡片**
   - 兑换金额
   - 汇率
   - 预计到账金额
   - **平台手续费**（百分比 + 金额）
   - 实际到账金额（加粗）
   - "确认兑换" 按钮

5. **兑换历史表格**

---

## P-10: 数币提现页 `/withdraw/crypto`

### Stitch 设计稿
- 已生成，Screen ID: `5ec64e2ece27455880247fc0887dfc3c`

### 布局（两列：60/40）
左列：
1. 币种选择 Tab（USDT / USDC / BTC）
2. 链选择下拉
3. 白名单地址选择（别名 + 截断地址）+ "添加新地址" 链接
4. 金额输入 + 可用余额 + MAX
5. 费用摘要（提现金额、平台手续费、预估链上费用、实际到账）
6. 提示信息（amber 背景）："提现提交后需平台审核，审核通过后自动处理"
7. "确认提现" 按钮

右列：
- 最近提现记录列表
  - 状态徽章增加：`Pending Review`（蓝色）、`Rejected`（红色）
  - 拒绝的记录显示拒绝原因（hover 或展开）

---

## P-11: 法币提现页 `/withdraw/fiat`

### 布局（类似数币提现）
左列：
1. 币种选择 Tab（USD / HKD / EUR）
2. 银行账户选择（账户名 + 截断账号）+ "添加银行账户" 链接
3. 金额输入 + 可用余额 + MAX
4. 转账用途选择：OTHER / TRAD / INVS / GDDS
5. 备注输入
6. 费用摘要（提现金额、平台手续费、实际到账）
7. 提示信息（amber 背景）："提现提交后需平台审核，审核通过后自动处理"
8. "确认提现" 按钮

右列：
- 最近法币提现记录
  - 状态徽章增加：`Pending Review`（蓝色）、`Rejected`（红色）
  - 拒绝的记录显示拒绝原因（hover 或展开）

---

## P-12: 地址管理页 `/settings/addresses`

### 布局
1. 页面标题 + "添加新地址" 按钮（右上）

2. **地址列表**（卡片网格或表格）
   - 每个地址卡片：
     - 别名（大字）
     - 完整地址（等宽字体）
     - 币种 + 链类型 徽章
     - 状态徽章（SUCCESS / PROCESS / FAILED）
     - 操作：复制地址 / 删除

3. **添加地址模态框**
   - 币种选择
   - 链类型选择
   - 地址输入
   - 别名输入
   - "确认添加" 按钮

---

## P-13: 银行账户管理页 `/settings/bank-accounts`

### 布局
1. 页面标题 + "添加银行账户" 按钮

2. **账户列表**
   - 每个账户卡片：
     - 银行名称 + 图标
     - 账户名
     - 账户号（截断显示）
     - 币种 + 转账类型 徽章
     - 操作：解绑

3. **添加账户模态框**
   - 币种选择（USD / HKD / EUR）
   - 转账类型（LOCAL / CHATS / TT）
   - 银行账户号
   - 账户名
   - SWIFT Code（TT/CHATS 时必填）
   - 银行名称
   - 收款人地址（TT/CHATS 时必填）
   - 中间行 SWIFT Code（可选）

---

## P-14: 交易记录页 `/transactions`

### 布局
1. 页面标题："Transaction History"

2. **筛选栏**
   - 交易类型下拉（All / Deposit / Exchange / Withdrawal / Transfer）
   - 状态下拉（All / Processing / Success / Failed）
   - 日期范围选择器
   - 搜索框（订单号）

3. **交易记录表格**
   - 列：
     - 类型（图标 + 文字）
     - 金额（绿色入账 / 红色出账）
     - 币种
     - 平台手续费
     - 状态（徽章）
     - 时间
     - 操作（查看详情）
   - 分页器

4. **交易详情侧边抽屉**（点击查看详情触发）
   - 平台订单号
   - KUN 订单号
   - 交易类型
   - 金额、币种
   - 平台手续费
   - KUN 手续费（如有）
   - 实际到账/扣款金额
   - 状态和时间线
   - 链上信息（TxID + explorer 链接，仅数币）

---

## 页面状态规范

所有页面需考虑以下状态：

| 状态 | 展示 |
|------|------|
| 加载中 | 骨架屏（Skeleton） |
| 空数据 | 插画 + "暂无数据" 文案 + 引导操作按钮 |
| 错误 | 红色 Alert + 错误信息 + "重试" 按钮 |
| 网络异常 | 全局 Toast 提示 |

---

## 响应式断点

| 断点 | 宽度 | 布局调整 |
|------|------|---------|
| Desktop | ≥1280px | 标准布局 |
| Tablet | 768-1279px | 双列变单列，侧边栏收起 |
| Mobile | <768px | 单列布局，底部 Tab 导航 |
