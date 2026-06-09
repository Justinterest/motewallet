import type { KycFormValues } from "@/lib/validations/onboarding";

type PersonFormValues = KycFormValues["shareholdersInfo"][number];

/** Map enterprise manager fields to shareholder/director person form values. */
export function managerInfoToPerson(
  enterprise: KycFormValues["enterpriseInfo"]
): PersonFormValues {
  const idHoldingPaths = (enterprise.managerIdHoldingPaths ?? []).filter(Boolean);

  return {
    gender: enterprise.managerGender || "Male",
    idCard: enterprise.managerIdCard || "",
    nameEN: enterprise.managerNameEN || "",
    country: enterprise.managerCountry || "",
    nameCHS: enterprise.managerNameCHS || "",
    surname: enterprise.managerSurname || "",
    authType: enterprise.managerAuthType || "",
    birthday: enterprise.managerBirthday || "",
    surnameCHS: enterprise.managerSurnameCHS || "",
    certificateStart: enterprise.managerCertificateStart || "",
    certificateEnd: enterprise.managerCertificateEnd || "",
    residenceAddress: enterprise.managerResidenceAddress || "",
    residenceCountry: enterprise.managerResidenceCountry || "",
    shareholdingRatio: "",
    idHoldingPaths: idHoldingPaths.length > 0 ? [...idHoldingPaths] : [""],
  };
}
