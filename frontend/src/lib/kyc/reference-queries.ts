import { onboardingApi } from "@/lib/api/onboarding";
import type { KycCountryScene } from "@/types/kyc-reference";

export const KYC_COUNTRY_SCENE_NATIONALITY: KycCountryScene = "REGISTER";
export const KYC_COUNTRY_SCENE_ADDRESS: KycCountryScene = "REGISTER_ADDRESS";
/** Default scene for address-related country fields. */
export const KYC_COUNTRY_SCENE: KycCountryScene = KYC_COUNTRY_SCENE_ADDRESS;

/** Shared React Query keys for KYC reference data (countries / auth types). */
export const kycReferenceKeys = {
  countries: (scene: KycCountryScene = KYC_COUNTRY_SCENE) =>
    ["onboarding", "countries", scene] as const,
  authTypes: (countryCode: string) =>
    ["onboarding", "auth-types", countryCode] as const,
};

const REFERENCE_STALE_MS = 24 * 60 * 60 * 1000;

export function kycCountriesQueryOptions(scene: KycCountryScene = KYC_COUNTRY_SCENE) {
  return {
    queryKey: kycReferenceKeys.countries(scene),
    queryFn: () => onboardingApi.getCountries({ scene, language: "ZH_CN" }),
    staleTime: REFERENCE_STALE_MS,
    gcTime: REFERENCE_STALE_MS,
    refetchOnMount: false,
    refetchOnWindowFocus: false,
  } as const;
}

export function kycAuthTypesQueryOptions(countryCode: string) {
  const code = countryCode.trim().toUpperCase();
  return {
    queryKey: kycReferenceKeys.authTypes(code),
    queryFn: () => onboardingApi.getCountryAuthTypes(code),
    staleTime: REFERENCE_STALE_MS,
    gcTime: REFERENCE_STALE_MS,
    refetchOnMount: false,
    refetchOnWindowFocus: false,
    enabled: Boolean(code),
  } as const;
}
