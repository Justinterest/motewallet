"use client";

import { type ReactNode } from "react";
import { useQueries } from "@tanstack/react-query";
import {
  KYC_COUNTRY_SCENE_ADDRESS,
  KYC_COUNTRY_SCENE_NATIONALITY,
  kycCountriesQueryOptions,
} from "@/lib/kyc/reference-queries";

/** Prefetches nationality and address country lists for the KYC wizard. */
export function KycReferenceProvider({ children }: { children: ReactNode }) {
  useQueries({
    queries: [
      kycCountriesQueryOptions(KYC_COUNTRY_SCENE_NATIONALITY),
      kycCountriesQueryOptions(KYC_COUNTRY_SCENE_ADDRESS),
    ],
  });

  return children;
}
