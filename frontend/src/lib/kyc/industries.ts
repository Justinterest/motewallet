/**
 * Industry enumeration from Kun docs:
 * https://opendocs.kun.global/docs/additional-details/enumeration-reference
 *
 * API fields `industry` and `subIndustry` expect level-1 / level-2 codes, not display names.
 */

export type IndustryLevel1 = {
  code: string;
  name: string;
};

export type IndustryLevel2 = {
  code: string;
  level1Code: string;
  level1Name: string;
  name: string;
};

export const INDUSTRY_LEVEL1: IndustryLevel1[] = [
  { code: "011", name: "Web3 服务" },
  { code: "001", name: "一般贸易" },
  { code: "002", name: "电商贸易" },
  { code: "003", name: "大宗贸易" },
  { code: "004", name: "物流与仓储" },
  { code: "005", name: "软件服务" },
  { code: "006", name: "广告与泛娱乐" },
  { code: "007", name: "金融与保险" },
  { code: "008", name: "教育服务" },
  { code: "009", name: "旅行与住宿" },
  { code: "010", name: "其他" },
];

export const INDUSTRY_LEVEL2: IndustryLevel2[] = [
  { code: "011001", level1Code: "011", level1Name: "Web3 服务", name: "加密货币交易平台 / 经纪商" },
  { code: "011002", level1Code: "011", level1Name: "Web3 服务", name: "其他数字资产服务提供商（VASP）" },
  { code: "011003", level1Code: "011", level1Name: "Web3 服务", name: "其他加密货币 / Web3 服务（非 VASP）" },
  { code: "001001", level1Code: "001", level1Name: "一般贸易", name: "烟草制品" },
  { code: "001002", level1Code: "001", level1Name: "一般贸易", name: "医疗器械" },
  { code: "001003", level1Code: "001", level1Name: "一般贸易", name: "化学制品" },
  { code: "001004", level1Code: "001", level1Name: "一般贸易", name: "计算机和 3C 数码" },
  { code: "001005", level1Code: "001", level1Name: "一般贸易", name: "芯片" },
  { code: "001006", level1Code: "001", level1Name: "一般贸易", name: "服饰箱包 / 皮革 / 鞋类 / 家居百货 / 日用品" },
  { code: "001007", level1Code: "001", level1Name: "一般贸易", name: "户外用品" },
  { code: "001008", level1Code: "001", level1Name: "一般贸易", name: "电器 / 体育用品 / 乐器 / 书籍 / 综合商品" },
  { code: "001009", level1Code: "001", level1Name: "一般贸易", name: "个人护理 / 母婴产品" },
  { code: "001010", level1Code: "001", level1Name: "一般贸易", name: "宠物园艺" },
  { code: "001011", level1Code: "001", level1Name: "一般贸易", name: "机械 / 电气设备及组件" },
  { code: "001012", level1Code: "001", level1Name: "一般贸易", name: "家具家装" },
  { code: "001013", level1Code: "001", level1Name: "一般贸易", name: "汽配用品" },
  { code: "001014", level1Code: "001", level1Name: "一般贸易", name: "食品酒水" },
  { code: "001015", level1Code: "001", level1Name: "一般贸易", name: "塑料和橡胶制品" },
  { code: "001016", level1Code: "001", level1Name: "一般贸易", name: "其他" },
  { code: "002001", level1Code: "002", level1Name: "电商贸易", name: "烟草制品" },
  { code: "002002", level1Code: "002", level1Name: "电商贸易", name: "户外用品" },
  { code: "002003", level1Code: "002", level1Name: "电商贸易", name: "汽配用品" },
  { code: "002004", level1Code: "002", level1Name: "电商贸易", name: "宠物园艺" },
  { code: "002005", level1Code: "002", level1Name: "电商贸易", name: "个人护理 / 母婴产品" },
  { code: "002006", level1Code: "002", level1Name: "电商贸易", name: "服饰箱包 / 皮革 / 鞋类 / 家居百货 / 日用品" },
  { code: "002007", level1Code: "002", level1Name: "电商贸易", name: "电器 / 体育用品 / 乐器 / 书籍 / 综合商品" },
  { code: "002008", level1Code: "002", level1Name: "电商贸易", name: "计算机和 3C 数码" },
  { code: "002009", level1Code: "002", level1Name: "电商贸易", name: "其他" },
  { code: "003001", level1Code: "003", level1Name: "大宗贸易", name: "能源化工" },
  { code: "003002", level1Code: "003", level1Name: "大宗贸易", name: "金属及其原材料" },
  { code: "003003", level1Code: "003", level1Name: "大宗贸易", name: "农产品" },
  { code: "003004", level1Code: "003", level1Name: "大宗贸易", name: "其他大宗贸易" },
  { code: "004001", level1Code: "004", level1Name: "物流与仓储", name: "物流运输" },
  { code: "004002", level1Code: "004", level1Name: "物流与仓储", name: "仓储" },
  { code: "005001", level1Code: "005", level1Name: "软件服务", name: "内容提供商 / 开发者" },
  { code: "005002", level1Code: "005", level1Name: "软件服务", name: "其他信息和技术服务" },
  { code: "006001", level1Code: "006", level1Name: "广告与泛娱乐", name: "广告" },
  { code: "006002", level1Code: "006", level1Name: "广告与泛娱乐", name: "游戏" },
  { code: "006003", level1Code: "006", level1Name: "广告与泛娱乐", name: "直播 / 语聊" },
  { code: "006004", level1Code: "006", level1Name: "广告与泛娱乐", name: "其他" },
  { code: "007001", level1Code: "007", level1Name: "金融与保险", name: "支付机构" },
  { code: "007002", level1Code: "007", level1Name: "金融与保险", name: "信用中介和相关活动" },
  { code: "007003", level1Code: "007", level1Name: "金融与保险", name: "证券 / 商品合约 / 其他金融投资" },
  { code: "007004", level1Code: "007", level1Name: "金融与保险", name: "保险相关活动" },
  { code: "007005", level1Code: "007", level1Name: "金融与保险", name: "基金 / 信托 / 其他金融工具" },
  { code: "008001", level1Code: "008", level1Name: "教育服务", name: "教育服务" },
  { code: "009001", level1Code: "009", level1Name: "旅行与住宿", name: "汽车租赁" },
  { code: "009002", level1Code: "009", level1Name: "旅行与住宿", name: "机票 / 酒店" },
  { code: "009003", level1Code: "009", level1Name: "旅行与住宿", name: "其他" },
  { code: "010001", level1Code: "010", level1Name: "其他", name: "其他" },
];

export const INDUSTRY_LEVEL1_OPTIONS = INDUSTRY_LEVEL1.map((item) => ({
  value: item.code,
  label: item.name,
}));

export function getSubIndustryOptions(level1Code: string) {
  const code = level1Code.trim();
  if (!code) return [];
  return INDUSTRY_LEVEL2.filter((item) => item.level1Code === code).map(
    (item) => ({
      value: item.code,
      label: item.name,
    })
  );
}

export function isSubIndustryForLevel1(
  subIndustryCode: string,
  level1Code: string
) {
  if (!subIndustryCode || !level1Code) return false;
  return INDUSTRY_LEVEL2.some(
    (item) => item.code === subIndustryCode && item.level1Code === level1Code
  );
}

/** Map legacy display-name drafts to Kun API industry codes. */
export function normalizeIndustryDraftFields<
  T extends { industry: string; subIndustry: string },
>(enterprise: T): T {
  let { industry, subIndustry } = enterprise;

  if (industry && !INDUSTRY_LEVEL1.some((item) => item.code === industry)) {
    const match = INDUSTRY_LEVEL1.find((item) => item.name === industry);
    if (match) industry = match.code;
  }

  if (subIndustry && !INDUSTRY_LEVEL2.some((item) => item.code === subIndustry)) {
    const byName = INDUSTRY_LEVEL2.find((item) => item.name === subIndustry);
    if (byName) {
      subIndustry = byName.code;
      if (!industry) industry = byName.level1Code;
    }
  }

  return { ...enterprise, industry, subIndustry };
}
