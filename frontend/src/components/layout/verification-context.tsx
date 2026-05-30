"use client";

import { createContext, useContext, useMemo } from "react";
import type { User } from "@/types/auth";
import { canMerchantOperate } from "@/lib/verification";

interface VerificationContextValue {
  user: User;
  canOperate: boolean;
}

const VerificationContext = createContext<VerificationContextValue | null>(
  null
);

export function VerificationProvider({
  user,
  children,
}: {
  user: User;
  children: React.ReactNode;
}) {
  const value = useMemo(
    () => ({
      user,
      canOperate: canMerchantOperate(user),
    }),
    [user]
  );

  return (
    <VerificationContext.Provider value={value}>
      {children}
    </VerificationContext.Provider>
  );
}

export function useVerification() {
  const ctx = useContext(VerificationContext);
  if (!ctx) {
    throw new Error("useVerification must be used within VerificationProvider");
  }
  return ctx;
}

/** Safe hook for optional provider (e.g. onboarding layout). */
export function useVerificationOptional() {
  return useContext(VerificationContext);
}
