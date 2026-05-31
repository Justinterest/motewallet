import type { FileRef, SubmitKycRequest } from "@/types/onboarding";
import type { KycFormValues } from "@/lib/validations/onboarding";

/** Staged S3 object keys; backend resolves to KUN URLs on submit. */
function pathsToFileRefs(paths?: string[]): FileRef[] {
  return (paths ?? []).filter(Boolean).map((path) => ({ path }));
}

function personToApi(
  p: KycFormValues["shareholdersInfo"][number]
): SubmitKycRequest["shareholdersInfo"][number] {
  const term =
    p.certificateStart && p.certificateEnd
      ? [p.certificateStart, p.certificateEnd]
      : undefined;

  return {
    gender: p.gender,
    idCard: p.idCard,
    nameEN: p.nameEN,
    country: p.country,
    nameCHS: p.nameCHS || undefined,
    surname: p.surname,
    authType: p.authType,
    birthday: p.birthday,
    surnameCHS: p.surnameCHS || undefined,
    certificateTerm: term,
    residenceAddress: p.residenceAddress,
    residenceCountry: p.residenceCountry,
    shareholdingRatio: p.shareholdingRatio || undefined,
    idHolding: pathsToFileRefs(p.idHoldingPaths),
    verificationType: "idHolding",
  };
}

/** Map validated form values to onboarding API request body. */
export function formValuesToSubmitRequest(
  values: KycFormValues
): SubmitKycRequest {
  const e = values.enterpriseInfo;

  return {
    enterpriseInfo: {
      incorporationCertificate: pathsToFileRefs(e.incorporationCertificatePaths),
      incorporationCertificateNo: e.incorporationCertificateNo,
      establishTime: e.establishTime,
      enterpriseEN: e.enterpriseEN,
      enterpriseNameCHS: e.enterpriseNameCHS || "无",
      registerRegion: e.registerRegion,
      registerAddress: e.registerAddress,
      businessRegistration: pathsToFileRefs(e.businessRegistrationPaths),
      businessRegistrationNo: e.businessRegistrationNo || undefined,
      phone: e.phone || undefined,
      isChangeEnterpriseNameInFiveYears: e.isChangeEnterpriseNameInFiveYears,
      usedEnterpriseName: e.usedEnterpriseName || undefined,
      enterpriseType: e.enterpriseType,
      enterpriseDomain: e.enterpriseDomain || undefined,
      businessRegion: e.businessRegion
        ? e.businessRegion.split(",").map((s) => s.trim()).filter(Boolean)
        : undefined,
      mainBusinessAddress: e.mainBusinessAddress || undefined,
      industry: e.industry,
      subIndustry: e.subIndustry,
      initialFundingSource: e.initialFundingSource,
      wealthSource: e.wealthSource || undefined,
      continuousFundingSource: e.continuousFundingSource || undefined,
      salesVolumeLastyear: e.salesVolumeLastyear || undefined,
      employeeNum: e.employeeNum || undefined,
      openAccountPurpose: e.openAccountPurpose,
      incumbency: pathsToFileRefs(e.incumbencyPaths),
      associationRules: pathsToFileRefs(e.associationRulesPaths),
      authenticMaterials: pathsToFileRefs(e.authenticMaterialsPaths),
      managerCountry: e.managerCountry,
      managerAuthType: e.managerAuthType,
      managerVerificationType: "idHolding",
      managerIdHolding: pathsToFileRefs(e.managerIdHoldingPaths),
      managerCertificateTerm:
        e.managerCertificateStart && e.managerCertificateEnd
          ? [e.managerCertificateStart, e.managerCertificateEnd]
          : undefined,
      managerSurnameCHS: e.managerSurnameCHS || undefined,
      managerNameCHS: e.managerNameCHS || undefined,
      managerSurname: e.managerSurname || undefined,
      managerNameEN: e.managerNameEN,
      managerBirthday: e.managerBirthday,
      managerGender: e.managerGender,
      managerIdCard: e.managerIdCard,
      managerResidenceCountry: e.managerResidenceCountry,
      managerResidenceAddress: e.managerResidenceAddress,
      managerContactsEmail: e.managerContactsEmail || undefined,
      authorizationLetter: pathsToFileRefs(e.authorizationLetterPaths),
      equityStructure: pathsToFileRefs(e.equityStructurePaths),
      middleTierShareholders: e.middleTierShareholders || undefined,
      nnc1: pathsToFileRefs(e.nnc1Paths),
    },
    shareholdersInfo: values.shareholdersInfo.map(personToApi),
    directorInfo: values.directorInfo.map((p) => {
      const { shareholdingRatio: _ratio, ...rest } = personToApi({
        ...p,
        shareholdingRatio: p.shareholdingRatio ?? "",
      });
      return rest;
    }),
  };
}
