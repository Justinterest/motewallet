"use client";

import { usePathname } from "next/navigation";
import { useVerificationOptional } from "./verification-context";

const VERIFICATION_EXEMPT_PREFIXES = ["/kyc"];

/**
 * Blocks pointer interaction on page content when the merchant is not ACTIVE.
 * Banner links remain clickable (rendered outside this wrapper).
 */
export function VerificationGuard({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const ctx = useVerificationOptional();
  const exempt = VERIFICATION_EXEMPT_PREFIXES.some(
    (prefix) => pathname === prefix || pathname.startsWith(`${prefix}/`)
  );
  const blocked = ctx && !ctx.canOperate && !exempt;

  return (
    <div className="relative">
      {children}
      {blocked && (
        <div
          className="absolute inset-0 z-20 cursor-not-allowed"
          aria-hidden
          title="请先完成企业认证"
        />
      )}
    </div>
  );
}
