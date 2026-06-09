import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { onboardingApi } from "@/lib/api/onboarding";
import {
  KYC_COUNTRY_SCENE,
  KYC_COUNTRY_SCENE_ADDRESS,
  KYC_COUNTRY_SCENE_NATIONALITY,
  kycAuthTypesQueryOptions,
  kycCountriesQueryOptions,
} from "@/lib/kyc/reference-queries";
import type { SubmitKycRequest } from "@/types/onboarding";
import type { KycCountryScene } from "@/types/kyc-reference";

export { KYC_COUNTRY_SCENE, KYC_COUNTRY_SCENE_ADDRESS, KYC_COUNTRY_SCENE_NATIONALITY };

export function useAgreements() {
  return useQuery({
    queryKey: ["onboarding", "agreements"],
    queryFn: () => onboardingApi.getAgreements(),
  });
}

export function useSignAgreements() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => onboardingApi.signAgreements(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["onboarding"] });
      queryClient.invalidateQueries({ queryKey: ["auth", "me"] });
    },
  });
}

export function useSubmitKyc() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: SubmitKycRequest) => onboardingApi.submitKyc(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["onboarding"] });
      queryClient.invalidateQueries({ queryKey: ["auth", "me"] });
    },
  });
}

export function useKycStatus() {
  return useQuery({
    queryKey: ["onboarding", "kyc-status"],
    queryFn: () => onboardingApi.getKycStatus(),
    refetchInterval: (query) => {
      const status = query.state.data?.kyc_status;
      return status === "AUTHING" ? 15000 : false;
    },
  });
}

/** Shares the same React Query cache as `KycCountrySelect` inside the KYC wizard. */
export function useKycCountries(scene: KycCountryScene = KYC_COUNTRY_SCENE) {
  return useQuery({
    ...kycCountriesQueryOptions(scene),
    select: (data) =>
      data.items.map((item) => ({
        value: item.country_code,
        label: item.country_name,
      })),
  });
}

export function useKycCountryAuthTypes(countryCode: string | undefined) {
  const code = countryCode?.trim().toUpperCase() ?? "";
  return useQuery({
    ...kycAuthTypesQueryOptions(code),
    select: (data) =>
      data.items.map((item) => ({
        value: item.doc_code,
        label: item.doc_name,
      })),
  });
}
