"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { AlertCircle, ArrowRight, Info, XCircle } from "lucide-react";
import { cn } from "@/lib/utils";
import { getVerificationBanner } from "@/lib/verification";
import { useVerificationOptional } from "./verification-context";

const variantStyles = {
  warning: {
    container: "border-amber-200 bg-amber-50 text-amber-950",
    icon: "text-amber-600",
    button: "bg-amber-700 hover:bg-amber-800 text-white",
  },
  info: {
    container: "border-blue-200 bg-blue-50 text-blue-950",
    icon: "text-blue-600",
    button: "bg-blue-700 hover:bg-blue-800 text-white",
  },
  error: {
    container: "border-red-200 bg-red-50 text-red-950",
    icon: "text-red-600",
    button: "bg-red-700 hover:bg-red-800 text-white",
  },
};

export function VerificationBanner() {
  const pathname = usePathname();
  const ctx = useVerificationOptional();
  const banner = getVerificationBanner(ctx?.user);

  if (!banner) {
    return null;
  }

  if (
    pathname === banner.ctaHref ||
    pathname.startsWith("/kyc") ||
    pathname.startsWith("/onboarding")
  ) {
    return null;
  }

  const styles = variantStyles[banner.variant];
  const Icon =
    banner.variant === "error"
      ? XCircle
      : banner.variant === "info"
        ? Info
        : AlertCircle;

  return (
    <div
      className={cn(
        "mb-6 flex flex-col gap-3 rounded-lg border px-4 py-3 sm:flex-row sm:items-center sm:justify-between",
        styles.container
      )}
      role="status"
    >
      <div className="flex items-start gap-3">
        <Icon className={cn("mt-0.5 h-5 w-5 shrink-0", styles.icon)} />
        <p className="text-sm leading-relaxed">{banner.message}</p>
      </div>
      <Link
        href={banner.ctaHref}
        className={cn(
          "inline-flex shrink-0 items-center justify-center rounded-md px-4 py-2 text-sm font-medium transition-colors",
          styles.button
        )}
      >
        {banner.ctaLabel}
        <ArrowRight className="ml-1.5 h-4 w-4" />
      </Link>
    </div>
  );
}
