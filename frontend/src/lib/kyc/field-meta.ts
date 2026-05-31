/** KYC 表单字段标签与 placeholder（面向商户用户） */

export type KycFieldMetaKey =
  | "requestNo"
  | "incorporationCertificate"
  | "incorporationCertificateNo"
  | "establishTime"
  | "enterpriseEN"
  | "enterpriseNameCHS"
  | "registerRegion"
  | "registerAddress"
  | "businessRegistration"
  | "businessRegistrationNo"
  | "phone"
  | "isChangeEnterpriseNameInFiveYears"
  | "usedEnterpriseName"
  | "enterpriseType"
  | "enterpriseDomain"
  | "businessRegion"
  | "mainBusinessAddress"
  | "industry"
  | "subIndustry"
  | "initialFundingSource"
  | "wealthSource"
  | "continuousFundingSource"
  | "salesVolumeLastyear"
  | "employeeNum"
  | "openAccountPurpose"
  | "incumbency"
  | "associationRules"
  | "authenticMaterials"
  | "managerCountry"
  | "managerAuthType"
  | "managerIdHolding"
  | "managerCertificateTerm"
  | "managerCertificateEnd"
  | "managerSurnameCHS"
  | "managerNameCHS"
  | "managerSurname"
  | "managerNameEN"
  | "managerBirthday"
  | "managerGender"
  | "managerIdCard"
  | "managerResidenceCountry"
  | "managerResidenceAddress"
  | "managerContactsEmail"
  | "authorizationLetter"
  | "equityStructure"
  | "middleTierShareholders"
  | "nnc1"
  | "person.gender"
  | "person.idCard"
  | "person.nameEN"
  | "person.country"
  | "person.nameCHS"
  | "person.surname"
  | "person.authType"
  | "person.birthday"
  | "person.surnameCHS"
  | "person.certificateTerm"
  | "person.certificateEnd"
  | "person.residenceAddress"
  | "person.residenceCountry"
  | "person.shareholdingRatio"
  | "person.idHolding";

export interface KycFieldMeta {
  label: string;
  description?: string;
  placeholder?: string;
  required?: boolean;
  exampleImage?: string;
}

const EXAMPLES = {
  incorporationCertificate: "/kyc/examples/incorporation-certificate.jpeg",
  associationRules: "/kyc/examples/association-rules.png",
} as const;

export const KYC_FIELD_META: Record<KycFieldMetaKey, KycFieldMeta> = {
  requestNo: {
    label: "申请流水号",
  },
  incorporationCertificate: {
    label: "公司注册证书",
    required: true,
    description:
      "请上传公司注册证书扫描件或照片。香港公司一般为公司注册证书（CI）。",
    exampleImage: EXAMPLES.incorporationCertificate,
  },
  incorporationCertificateNo: {
    label: "公司注册证书编号",
    required: true,
    placeholder: "与注册证书一致",
  },
  establishTime: {
    label: "成立日期",
    required: true,
  },
  enterpriseEN: {
    label: "企业英文名称",
    required: true,
    placeholder: "与注册证书英文名称一致",
  },
  enterpriseNameCHS: {
    label: "企业中文名称",
    placeholder: "无中文名填「无」",
  },
  registerRegion: {
    label: "公司注册地",
    required: true,
  },
  registerAddress: {
    label: "注册地址",
    required: true,
    placeholder: "英文地址，香港须与注册文件一致",
  },
  businessRegistration: {
    label: "商业登记证",
    description: "注册地为香港时，须上传商业登记证（BR）扫描件。",
  },
  businessRegistrationNo: {
    label: "商业登记证编号",
    placeholder: "香港企业必填",
  },
  phone: {
    label: "企业联系电话",
    required: true,
  },
  isChangeEnterpriseNameInFiveYears: {
    label: "过去 5 年是否更改过企业名称",
    required: true,
  },
  usedEnterpriseName: {
    label: "企业曾用名",
    placeholder: "更改过名称时填写",
  },
  enterpriseType: {
    label: "企业类型",
    required: true,
  },
  enterpriseDomain: {
    label: "公司网址",
    placeholder: "https://",
  },
  businessRegion: {
    label: "业务所在国家/地区",
    placeholder: "例如 HK, SG（多个用逗号分隔）",
  },
  mainBusinessAddress: {
    label: "营业地址",
    required: true,
    placeholder: "建议使用英文",
  },
  industry: {
    label: "一级行业",
    required: true,
  },
  subIndustry: {
    label: "二级行业",
    required: true,
    placeholder: "须属于已选一级行业，如 B2B Trade",
  },
  initialFundingSource: {
    label: "原始资金来源",
    required: true,
  },
  wealthSource: {
    label: "财富来源",
    required: true,
  },
  continuousFundingSource: {
    label: "持续资金来源",
    required: true,
  },
  salesVolumeLastyear: {
    label: "去年销售额（港币）",
    required: true,
  },
  employeeNum: {
    label: "员工人数",
  },
  openAccountPurpose: {
    label: "开户目的",
    required: true,
  },
  incumbency: {
    label: "董事在职证明",
  },
  associationRules: {
    label: "公司章程",
    required: true,
    description: "请上传公司章程或组织大纲扫描件。",
    exampleImage: EXAMPLES.associationRules,
  },
  authenticMaterials: {
    label: "真实性证明材料",
    required: true,
    description:
      "用于佐证所提交资料真实性的支持文件，按审核要求提供。",
  },
  managerCountry: {
    label: "账户管理人国籍",
    required: true,
  },
  managerAuthType: {
    label: "管理人证件类型",
    required: true,
    placeholder: "请先选择国籍",
  },
  managerIdHolding: {
    label: "管理人手持证件照",
    required: true,
    description:
      "须至少上传 3 张：证件照、本人自拍、手持证件照。",
  },
  managerCertificateTerm: {
    label: "证件有效期起",
    required: true,
  },
  managerCertificateEnd: {
    label: "证件有效期止",
    required: true,
  },
  managerSurnameCHS: {
    label: "管理人中文姓氏",
    required: true,
    placeholder: "仅填姓氏",
  },
  managerNameCHS: {
    label: "管理人中文名",
    required: true,
  },
  managerSurname: {
    label: "管理人英文姓氏",
    required: true,
    placeholder: "Surname",
  },
  managerNameEN: {
    label: "管理人英文名",
    required: true,
    placeholder: "Given name",
  },
  managerBirthday: {
    label: "管理人出生日期",
    required: true,
  },
  managerGender: {
    label: "管理人性别",
    required: true,
  },
  managerIdCard: {
    label: "管理人证件号码",
    required: true,
  },
  managerResidenceCountry: {
    label: "管理人居住国家/地区",
    required: true,
  },
  managerResidenceAddress: {
    label: "管理人居住地址",
    required: true,
  },
  managerContactsEmail: {
    label: "管理人联系邮箱",
    placeholder: "name@company.com",
  },
  authorizationLetter: {
    label: "授权书",
  },
  equityStructure: {
    label: "股权结构图",
    description: "体现股东层级与持股比例的结构说明文件。",
  },
  middleTierShareholders: {
    label: "是否有中间层股东",
    required: true,
    description: "是否存在中间控股公司或代持层级。",
  },
  nnc1: {
    label: "NNC1 表格",
    description: "注册地为香港时须上传公司成立表格 NNC1 扫描件。",
  },
  "person.gender": {
    label: "性别",
    required: true,
  },
  "person.idCard": {
    label: "证件号码",
    required: true,
  },
  "person.nameEN": {
    label: "英文名",
    required: true,
    placeholder: "Given name",
  },
  "person.country": {
    label: "国籍",
    required: true,
  },
  "person.nameCHS": {
    label: "中文名",
    required: true,
  },
  "person.surname": {
    label: "英文姓氏",
    required: true,
    placeholder: "Surname",
  },
  "person.authType": {
    label: "证件类型",
    required: true,
    placeholder: "请先选择国籍",
  },
  "person.birthday": {
    label: "出生日期",
    required: true,
  },
  "person.surnameCHS": {
    label: "中文姓氏",
    required: true,
    placeholder: "仅填姓氏",
  },
  "person.certificateTerm": {
    label: "证件有效期起",
    required: true,
  },
  "person.certificateEnd": {
    label: "证件有效期止",
    required: true,
  },
  "person.residenceAddress": {
    label: "居住地址",
    required: true,
  },
  "person.residenceCountry": {
    label: "居住国家/地区",
    required: true,
  },
  "person.shareholdingRatio": {
    label: "持股比例 (%)",
    required: true,
    placeholder: "如 25 表示 25%",
  },
  "person.idHolding": {
    label: "手持证件照",
    required: true,
    description: "请上传证件照、本人自拍、手持证件照等相关照片。",
  },
};

export function getKycFieldMeta(key: KycFieldMetaKey): KycFieldMeta {
  return KYC_FIELD_META[key];
}
