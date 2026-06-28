import apiClient from "./client";
import type { AgreementListResponse, KycStatusResponse, SubmitKycRequest } from "@/types/onboarding";
import type {
  CountryAuthTypesResponse,
  CountryListResponse,
  KycCountryScene,
} from "@/types/kyc-reference";

export const onboardingApi = {
  getAgreements: () => apiClient.get<never, AgreementListResponse>("/api/v1/onboarding/agreements"),
  signAgreements: () => apiClient.post<never, void>("/api/v1/onboarding/agreements/sign"),
  submitKyc: (data: SubmitKycRequest) => apiClient.post<never, void>("/api/v1/onboarding/kyc", data),
  getKycStatus: () => apiClient.get<never, KycStatusResponse>("/api/v1/onboarding/kyc/status"),
  getCountries: (params?: { scene?: KycCountryScene; language?: string; currency?: string }) =>
    apiClient.get<never, CountryListResponse>("/api/v1/onboarding/countries", { params }),
  getCountryAuthTypes: (countryCode: string) =>
    apiClient.get<never, CountryAuthTypesResponse>(
      `/api/v1/onboarding/countries/${encodeURIComponent(countryCode)}/auth-types`
    ),
};
