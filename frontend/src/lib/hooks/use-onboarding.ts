import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { onboardingApi } from "@/lib/api/onboarding";
import type { SubmitKycRequest } from "@/types/onboarding";

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
