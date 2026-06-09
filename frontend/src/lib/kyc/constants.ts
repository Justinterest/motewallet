export const REGISTER_REGIONS = [
  { value: "HK", label: "中国香港" },
  { value: "CN", label: "中国大陆" },
  { value: "SG", label: "新加坡" },
  { value: "US", label: "美国" },
  { value: "GB", label: "英国" },
];

export const ENTERPRISE_TYPES = [
  { value: "Private limited company", label: "私人有限公司" },
  { value: "Public limited company", label: "公众有限公司" },
];

export const FUNDING_SOURCES = [
  { value: "Business income", label: "营业收入 / 商业利润" },
  { value: "Shareholder/investor funds", label: "股东 / 投资者资金" },
  { value: "Investment", label: "投资收益" },
  { value: "Digital currency mining", label: "数字货币挖矿" },
  { value: "Digital currency pledge", label: "数字货币质押" },
  { value: "Others", label: "其他来源" },
];

/** 财富来源（比资金来源多「资产出售」选项） */
export const WEALTH_SOURCES = [
  ...FUNDING_SOURCES.slice(0, 3),
  { value: "Asset sale", label: "资产出售所得" },
  ...FUNDING_SOURCES.slice(3),
];

export const SALES_VOLUME_OPTIONS = [
  { value: "HKD 0-2,500,000", label: "HKD 0–2,500,000" },
  { value: "HKD 2,500,001-5,000,000", label: "HKD 2,500,001–5,000,000" },
  { value: "HKD 5,000,001-10,000,000", label: "HKD 5,000,001–10,000,000" },
  { value: "HKD 10,000,001-30,000,000", label: "HKD 10,000,001–30,000,000" },
  { value: "HKD 30,000,001-50,000,000", label: "HKD 30,000,001–50,000,000" },
  { value: "HKD 50,000,001或以上", label: "HKD 50,000,001 或以上" },
];

export const OPEN_ACCOUNT_PURPOSES = [
  { value: "Business operation", label: "日常企业经营" },
  { value: "Remittance", label: "汇款" },
  { value: "Currency exchange", label: "换汇" },
];

export const EMPLOYEE_NUM_OPTIONS = [
  { value: "<10", label: "少于 10 人" },
  { value: "10-99", label: "10–99 人" },
  { value: "100-500", label: "100–500 人" },
  { value: ">500", label: "500 人以上" },
];

export const GENDERS = [
  { value: "Male", label: "男" },
  { value: "Female", label: "女" },
];

export const YES_NO = [
  { value: "No", label: "否" },
  { value: "Yes", label: "是" },
];

export const KYC_DRAFT_KEY = "motewallet_kyc_draft_v2";
