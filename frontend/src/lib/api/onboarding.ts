import apiClient from "./client";
import type { AgreementListResponse, KycStatusResponse, SubmitKycRequest } from "@/types/onboarding";

export const onboardingApi = {
  getAgreements: () => apiClient.get<never, AgreementListResponse>("/api/v1/onboarding/agreements"),
  signAgreements: () => apiClient.post<never, void>("/api/v1/onboarding/agreements/sign"),
  submitKyc: (data: SubmitKycRequest) => apiClient.post<never, void>("/api/v1/onboarding/kyc", data),
  getKycStatus: () => apiClient.get<never, KycStatusResponse>("/api/v1/onboarding/kyc/status"),
};
