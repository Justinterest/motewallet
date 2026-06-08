export interface Agreement {
  id: string;
  protocol_id: string;
  title: string;
  version?: string;
  url: string;
  required: boolean;
}

export interface AgreementListResponse {
  agreements: Agreement[];
  signed: boolean;
}

export interface KycStatusResponse {
  status: string;
  kyc_status: string;
  kyc_fail_reason?: string;
  kyc_submitted_at?: string;
  kyc_completed_at?: string;
  agreement_signed_at?: string;
}

export interface FileRef {
  path: string;
}

export interface EnterpriseInfo {
  incorporationCertificate?: FileRef[];
  incorporationCertificateNo: string;
  establishTime: string;
  enterpriseEN: string;
  enterpriseNameCHS: string;
  registerRegion: string;
  registerAddress: string;
  businessRegistration?: FileRef[];
  businessRegistrationNo?: string;
  phone?: string;
  isChangeEnterpriseNameInFiveYears?: string;
  usedEnterpriseName?: string;
  enterpriseType: string;
  enterpriseDomain?: string;
  businessRegion?: string[];
  mainBusinessAddress?: string;
  industry: string;
  subIndustry: string;
  initialFundingSource: string;
  wealthSource?: string;
  continuousFundingSource?: string;
  salesVolumeLastyear?: string;
  employeeNum?: string;
  openAccountPurpose: string;
  incumbency?: FileRef[];
  associationRules?: FileRef[];
  authenticMaterials?: FileRef[];
  managerCountry: string;
  managerAuthType: string;
  managerVerificationType: string;
  managerIdHolding?: FileRef[];
  managerCertificateTerm?: string[];
  managerSurnameCHS?: string;
  managerNameCHS?: string;
  managerSurname?: string;
  managerNameEN: string;
  managerBirthday: string;
  managerGender: string;
  managerIdCard: string;
  managerResidenceCountry: string;
  managerResidenceAddress: string;
  managerContactsEmail?: string;
  authorizationLetter?: FileRef[];
  equityStructure?: FileRef[];
  middleTierShareholders?: string;
  nnc1?: FileRef[];
}

export interface PersonInfo {
  gender: string;
  idCard: string;
  nameEN: string;
  country: string;
  nameCHS?: string;
  surname: string;
  authType: string;
  birthday: string;
  surnameCHS?: string;
  certificateTerm?: string[];
  residenceAddress: string;
  residenceCountry: string;
  shareholdingRatio?: string;
  idHolding?: FileRef[];
  verificationType?: string;
}

export interface SubmitKycRequest {
  requestNo?: string;
  enterpriseInfo: EnterpriseInfo;
  shareholdersInfo: PersonInfo[];
  directorInfo: PersonInfo[];
}
