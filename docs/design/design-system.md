# Motewallet 设计系统

## Stitch 设计项目

- **项目 ID**: `7183872683002954167`
- **设计系统 ID**: `assets/16511664270239647761`
- **项目链接**: Stitch 平台内查看

## 品牌调性

Motewallet 是面向 B2B 商户的金融平台，设计应传达：**专业、安全、清晰**。

## 色彩系统

| 色彩 | Hex | 用途 |
|------|-----|------|
| Primary Blue | `#1E40AF` | 主操作、导航高亮、主按钮 |
| Secondary Sky | `#0EA5E9` | 次要操作、链接、信息态 |
| Success Green | `#10B981` | 成功态、正数金额、充值指示 |
| Error Red | `#EF4444` | 错误态、负数金额、提现警告 |
| Warning Amber | `#F59E0B` | 待处理态、警告指示 |
| Neutral Slate | `#64748B` | 正文、边框、禁用态 |
| Background | `#F8FAFC` | 页面背景 |
| Card Background | `#FFFFFF` | 卡片背景 |
| Border | `#E2E8F0` | 卡片和分割线边框 |

## 字体

| 层级 | 字号/字重 | 用途 |
|------|----------|------|
| Display | 32px / 700 | 页面标题 |
| Heading | 24px / 600 | 区块标题 |
| Subheading | 18px / 600 | 卡片标题 |
| Body | 14px / 400 | 正文 |
| Caption | 12px / 400 | 标签、时间戳 |
| Amount | 24px / 700 tabular-nums | 金额显示 |

字体家族：**Inter**（全局统一）

## 间距

- 基础单位：4px
- 组件内边距：16px
- 卡片间距：24px
- 区块间距：32px

## 圆角

- 按钮：8px
- 卡片：8px
- 输入框：8px
- 徽章/Tag：9999px（full round）
- 模态框：12px

## 阴影

- 卡片：`0 1px 3px rgba(0,0,0,0.1), 0 1px 2px rgba(0,0,0,0.06)`
- 下拉菜单：`0 4px 6px rgba(0,0,0,0.1), 0 2px 4px rgba(0,0,0,0.06)`
- 模态框：`0 20px 25px rgba(0,0,0,0.1), 0 10px 10px rgba(0,0,0,0.04)`

## 组件规范

### 按钮

| 变体 | 样式 | 用途 |
|------|------|------|
| Primary | 蓝色填充 `bg-blue-800 text-white` | 主操作（确认兑换、确认提现） |
| Secondary | 蓝色描边 `border-blue-800 text-blue-800` | 次要操作（取消、返回） |
| Danger | 红色填充 `bg-red-500 text-white` | 危险操作（删除、冻结） |
| Ghost | 透明 `text-blue-800` | 文字链接式按钮 |

尺寸：
- Large: h-12 px-6 text-base（表单主按钮）
- Default: h-10 px-4 text-sm
- Small: h-8 px-3 text-xs（表格内操作）

### 状态徽章

| 状态 | 样式 |
|------|------|
| Success / Completed | `bg-green-100 text-green-800` |
| Processing | `bg-amber-100 text-amber-800` |
| Pending | `bg-blue-100 text-blue-800` |
| Pending Review | `bg-purple-100 text-purple-800` |
| Failed | `bg-red-100 text-red-800` |
| Rejected | `bg-red-100 text-red-800` |
| Cancelled | `bg-slate-100 text-slate-600` |

### 金额显示规则

- 使用 `tabular-nums` 确保数字等宽对齐
- 正数金额（入账）：绿色 `text-green-600` 带 `+` 前缀
- 负数金额（出账）：红色 `text-red-600` 带 `-` 前缀
- 始终显示币种代码：`1,000.00 USDT`、`$8,000.00 USD`
- 右对齐
- 法币精度：2 位小数
- BTC 精度：8 位小数
- USDT/USDC 精度：2 位小数

### 卡片

```
白色背景 + subtle 阴影 + 8px 圆角 + 1px border-slate-200
内边距 24px
标题区域 + 分割线 + 内容区域
```

### 表格

- 表头：`bg-slate-50 text-slate-600 text-xs uppercase tracking-wider`
- 行间距：交替背景 `even:bg-slate-50`
- 单元格内边距：`px-4 py-3`
- 操作列靠右

### 表单输入

- 高度：h-10
- 边框：`border-slate-300 focus:border-blue-500 focus:ring-blue-500`
- 标签：`text-sm font-medium text-slate-700` 在输入框上方
- 错误态：`border-red-500 text-red-500` + 错误文本在下方

## 金融 UI 专项规则

1. **金额旁始终显示币种代码**
2. **绿色表示入账，红色表示出账**
3. **确认操作前展示费用摘要卡片**（含平台手续费明细）
4. **提现和兑换需要二步确认**（输入 → 摘要确认）
5. **敏感信息（地址、银行卡号）截断显示**：`0x1234...abcd`
