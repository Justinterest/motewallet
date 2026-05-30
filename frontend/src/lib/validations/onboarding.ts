import { z } from "zod";

const enterpriseStepSchema = z.object({
  incorporationCertificatePaths: z.array(z.string()).min(1, "请填写营业执照文件 path"),
  incorporationCertificateNo: z.string().min(1, "请填写注册证书编号"),
  establishTime: z.string().min(1, "请填写成立日期"),
  enterpriseEN: z.string().min(2, "请填写公司英文名称"),
  enterpriseNameCHS: z.string().optional(),
  registerRegion: z.string().min(1, "请选择注册地区"),
  registerAddress: z.string().min(2, "请填写注册地址"),
  businessRegistrationPaths: z.array(z.string()).optional(),
  businessRegistrationNo: z.string().optional(),
  phone: z.string().optional(),
  isChangeEnterpriseNameInFiveYears: z.string().optional(),
  usedEnterpriseName: z.string().optional(),
  enterpriseType: z.string().min(1, "请选择公司类型"),
  enterpriseDomain: z.string().optional(),
  businessRegion: z.string().optional(),
  mainBusinessAddress: z.string().optional(),
  industry: z.string().min(1, "请填写一级行业 code"),
  subIndustry: z.string().min(1, "请填写二级行业 code"),
  initialFundingSource: z.string().min(1, "请选择初始资金来源"),
  wealthSource: z.string().optional(),
  continuousFundingSource: z.string().optional(),
  salesVolumeLastyear: z.string().optional(),
  employeeNum: z.string().optional(),
  openAccountPurpose: z.string().min(1, "请选择开户目的"),
  incumbencyPaths: z.array(z.string()).optional(),
  associationRulesPaths: z.array(z.string()).optional(),
  authenticMaterialsPaths: z.array(z.string()).optional(),
});

const managerStepSchema = z.object({
  managerCountry: z.string().min(1, "请选择管理人国籍"),
  managerAuthType: z.string().min(1, "请填写证件类型 code"),
  managerVerificationType: z.string().min(1, "请选择验证方式"),
  managerIdHoldingPaths: z.array(z.string()).optional(),
  managerCertificateStart: z.string().optional(),
  managerCertificateEnd: z.string().optional(),
  managerSurnameCHS: z.string().optional(),
  managerNameCHS: z.string().optional(),
  managerSurname: z.string().optional(),
  managerNameEN: z.string().min(1, "请填写管理人英文名"),
  managerBirthday: z.string().min(1, "请填写出生日期"),
  managerGender: z.string().min(1, "请选择性别"),
  managerIdCard: z.string().min(1, "请填写证件号码"),
  managerResidenceCountry: z.string().min(1, "请选择居住国家"),
  managerResidenceAddress: z.string().min(2, "请填写居住地址"),
  managerContactsEmail: z.string().email("邮箱格式不正确").optional().or(z.literal("")),
  authorizationLetterPaths: z.array(z.string()).optional(),
  equityStructurePaths: z.array(z.string()).optional(),
  middleTierShareholders: z.string().optional(),
  nnc1Paths: z.array(z.string()).optional(),
});

export const personSchema = z.object({
  gender: z.string().min(1, "请选择性别"),
  idCard: z.string().min(1, "请填写证件号"),
  nameEN: z.string().min(1, "请填写英文姓名"),
  country: z.string().min(1, "请选择国籍"),
  nameCHS: z.string().optional(),
  surname: z.string().min(1, "请填写英文姓"),
  authType: z.string().min(1, "请填写证件类型 code"),
  birthday: z.string().min(1, "请填写出生日期"),
  surnameCHS: z.string().optional(),
  certificateStart: z.string().optional(),
  certificateEnd: z.string().optional(),
  residenceAddress: z.string().min(2, "请填写居住地址"),
  residenceCountry: z.string().min(1, "请选择居住国家"),
  shareholdingRatio: z.string().optional(),
  idHoldingPaths: z.array(z.string()).optional(),
  verificationType: z.string().optional(),
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

    if (data.enterpriseInfo.managerVerificationType === "idHolding") {
      const paths = data.enterpriseInfo.managerIdHoldingPaths ?? [];
      if (paths.filter(Boolean).length < 3) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "手持证件方式须至少填写 3 个文件 path",
          path: ["enterpriseInfo", "managerIdHoldingPaths"],
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
