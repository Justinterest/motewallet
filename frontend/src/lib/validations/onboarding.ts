import { z } from "zod";

const enterpriseStepSchema = z.object({
  incorporationCertificatePaths: z
    .array(z.string())
    .refine((arr) => arr.some(Boolean), "请上传营业执照 / 注册证书"),
  incorporationCertificateNo: z.string().min(1, "请填写注册证书编号"),
  establishTime: z.string().min(1, "请填写成立日期"),
  enterpriseEN: z.string().min(2, "请填写公司英文名称"),
  enterpriseNameCHS: z.string().optional(),
  registerRegion: z.string().min(1, "请选择注册地区"),
  registerAddress: z.string().min(2, "请填写注册地址"),
  businessRegistrationPaths: z.array(z.string()).optional(),
  businessRegistrationNo: z.string().optional(),
  phone: z.string().min(1, "请填写企业电话"),
  isChangeEnterpriseNameInFiveYears: z.string().min(1, "请选择是否更改过企业名称"),
  usedEnterpriseName: z.string().optional(),
  enterpriseType: z.string().min(1, "请选择公司类型"),
  enterpriseDomain: z.string().optional(),
  businessRegion: z.string().optional(),
  mainBusinessAddress: z.string().min(2, "请填写营业地址"),
  industry: z.string().min(1, "请选择一级行业"),
  subIndustry: z.string().min(1, "请填写二级行业"),
  initialFundingSource: z.string().min(1, "请选择原始资金来源"),
  wealthSource: z.string().min(1, "请选择财富来源"),
  continuousFundingSource: z.string().min(1, "请选择持续资金来源"),
  salesVolumeLastyear: z.string().min(1, "请选择去年销售额"),
  employeeNum: z.string().optional(),
  openAccountPurpose: z.string().min(1, "请选择开户目的"),
  incumbencyPaths: z.array(z.string()).optional(),
  associationRulesPaths: z
    .array(z.string())
    .refine((arr) => arr.some(Boolean), "请上传公司章程"),
  authenticMaterialsPaths: z
    .array(z.string())
    .refine((arr) => arr.some(Boolean), "请上传真实性证明材料"),
});

const managerStepSchema = z.object({
  managerCountry: z.string().min(1, "请选择管理人国籍"),
  managerAuthType: z.string().min(1, "请选择证件类型"),
  managerIdHoldingPaths: z
    .array(z.string())
    .refine((arr) => arr.filter(Boolean).length >= 3, "请上传 3 张手持证件照"),
  managerCertificateStart: z.string().optional(),
  managerCertificateEnd: z.string().optional(),
  managerSurnameCHS: z.string().min(1, "请填写管理人中文姓氏"),
  managerNameCHS: z.string().min(1, "请填写管理人中文名"),
  managerSurname: z.string().min(1, "请填写管理人英文姓氏"),
  managerNameEN: z.string().min(1, "请填写管理人英文名"),
  managerBirthday: z.string().min(1, "请填写出生日期"),
  managerGender: z.string().min(1, "请选择性别"),
  managerIdCard: z.string().min(1, "请填写证件号码"),
  managerResidenceCountry: z.string().min(1, "请选择居住国家"),
  managerResidenceAddress: z.string().min(2, "请填写居住地址"),
  managerContactsEmail: z.string().email("邮箱格式不正确").optional().or(z.literal("")),
  authorizationLetterPaths: z.array(z.string()).optional(),
  equityStructurePaths: z.array(z.string()).optional(),
  middleTierShareholders: z.string().min(1, "请选择是否有中间层股东"),
  nnc1Paths: z.array(z.string()).optional(),
});

export const personSchema = z.object({
  gender: z.string().min(1, "请选择性别"),
  idCard: z.string().min(1, "请填写证件号"),
  nameEN: z.string().min(1, "请填写英文姓名"),
  country: z.string().min(1, "请选择国籍"),
  nameCHS: z.string().min(1, "请填写中文名"),
  surname: z.string().min(1, "请填写英文姓"),
  authType: z.string().min(1, "请选择证件类型"),
  birthday: z.string().min(1, "请填写出生日期"),
  surnameCHS: z.string().min(1, "请填写中文姓氏"),
  certificateStart: z.string().optional(),
  certificateEnd: z.string().optional(),
  residenceAddress: z.string().min(2, "请填写居住地址"),
  residenceCountry: z.string().min(1, "请选择居住国家"),
  shareholdingRatio: z.string().optional(),
  idHoldingPaths: z
    .array(z.string())
    .refine((arr) => arr.some(Boolean), "请上传手持证件照"),
});

const shareholderSchema = personSchema.extend({
  shareholdingRatio: z.string().min(1, "请填写持股比例"),
});

const kycFormBaseSchema = z.object({
  enterpriseInfo: enterpriseStepSchema.merge(managerStepSchema),
  shareholdersInfo: z.array(shareholderSchema).min(1, "至少添加一名股东"),
  directorInfo: z.array(personSchema).min(1, "至少添加一名董事"),
});

export const kycFormSchema = kycFormBaseSchema
  .superRefine((data, ctx) => {
    const region = data.enterpriseInfo.registerRegion;
    if (region === "HK") {
      if (!data.enterpriseInfo.businessRegistrationNo) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "香港企业须填写商业登记证编号",
          path: ["enterpriseInfo", "businessRegistrationNo"],
        });
      }
      const brPaths = data.enterpriseInfo.businessRegistrationPaths;
      if (!brPaths || brPaths.length === 0 || !brPaths.some(Boolean)) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "香港企业须上传商业登记证",
          path: ["enterpriseInfo", "businessRegistrationPaths"],
        });
      }
    }

    if (data.enterpriseInfo.isChangeEnterpriseNameInFiveYears === "Yes") {
      if (!data.enterpriseInfo.usedEnterpriseName?.trim()) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "更改过企业名称时须填写曾用名",
          path: ["enterpriseInfo", "usedEnterpriseName"],
        });
      }
    }

    if (data.enterpriseInfo.middleTierShareholders === "Yes") {
      const equityPaths = data.enterpriseInfo.equityStructurePaths;
      if (!equityPaths || equityPaths.length === 0 || !equityPaths.some(Boolean)) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "存在中间层股东时须上传股权结构图",
          path: ["enterpriseInfo", "equityStructurePaths"],
        });
      }
    }

  });

export type KycFormValues = z.infer<typeof kycFormSchema>;

export const kycStep0Schema = z.object({ enterpriseInfo: enterpriseStepSchema });
export const kycStep1Schema = z.object({ enterpriseInfo: managerStepSchema });

export const kycStep2Schema = z.object({
  shareholdersInfo: kycFormBaseSchema.shape.shareholdersInfo,
  directorInfo: kycFormBaseSchema.shape.directorInfo,
});
