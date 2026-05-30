"use client";

import { useVerificationOptional } from "./verification-context";

/**
 * Blocks pointer interaction on page content when the merchant is not ACTIVE.
 * Banner links remain clickable (rendered outside this wrapper).
 */
export function VerificationGuard({ children }: { children: React.ReactNode }) {
  const ctx = useVerificationOptional();
  const blocked = ctx && !ctx.canOperate;

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
