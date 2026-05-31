"use client";

import { createContext, useContext, useMemo, type ReactNode } from "react";
import { useQuery } from "@tanstack/react-query";
import type { SelectOption } from "@/components/ui/select";
import {
  KYC_COUNTRY_SCENE,
  kycCountriesQueryOptions,
} from "@/lib/kyc/reference-queries";

interface KycReferenceContextValue {
  countryOptions: SelectOption[];
  countriesLoading: boolean;
  countriesError: boolean;
}

const KycReferenceContext = createContext<KycReferenceContextValue | null>(null);

/** Loads country/region list once for the whole KYC wizard. */
export function KycReferenceProvider({ children }: { children: ReactNode }) {
  const { data, isLoading, isError } = useQuery({
    ...kycCountriesQueryOptions(KYC_COUNTRY_SCENE),
    select: (response) =>
      response.items.map((item) => ({
        value: item.country_code,
        label: item.country_name,
      })),
  });

  const value = useMemo(
    () => ({
      countryOptions: data ?? [],
      countriesLoading: isLoading,
      countriesError: isError,
    }),
    [data, isLoading, isError]
  );

  return (
    <KycReferenceContext.Provider value={value}>
      {children}
    </KycReferenceContext.Provider>
  );
}

export function useKycReference() {
  const ctx = useContext(KycReferenceContext);
  if (!ctx) {
    throw new Error("useKycReference must be used within KycReferenceProvider");
  }
  return ctx;
}
