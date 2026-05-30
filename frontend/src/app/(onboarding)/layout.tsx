"use client";

import { useEffect } from "react";
import Link from "next/link";
import { useRouter, usePathname } from "next/navigation";
import { CheckCircle2 } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { useCurrentUser } from "@/lib/hooks/use-auth";
import { cn } from "@/lib/utils";

const steps = [
  { label: "签署协议", path: "/onboarding/agreement" },
  { label: "实名认证", path: "/onboarding/kyc" },
  { label: "认证状态", path: "/onboarding/status" },
];

export default function OnboardingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const { data: user, isLoading, isError } = useCurrentUser();

  useEffect(() => {
    if (isError) {
      router.push("/login");
    }
  }, [isError, router]);

  if (isLoading) {
    return (
      <div className="flex min-h-screen flex-col items-center bg-slate-50 px-4 pt-16">
        <div className="w-full max-w-2xl">
          <Skeleton className="mx-auto mb-8 h-8 w-48" />
          <Skeleton className="mx-auto mb-12 h-12 w-full max-w-md" />
          <Skeleton className="h-80 w-full rounded-lg" />
        </div>
      </div>
    );
  }

  if (isError || !user) {
    return null;
  }

  const currentStepIndex = steps.findIndex((s) => pathname.startsWith(s.path));

  return (
    <div className="flex min-h-screen flex-col items-center bg-slate-50 px-4 pt-16">
      <div className="w-full max-w-2xl">
        {/* Header */}
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold text-blue-800">Motewallet</h1>
          <p className="mt-1 text-sm text-slate-500">商户入驻流程</p>
          <Link
            href="/dashboard"
            className="mt-2 inline-block text-sm text-blue-600 hover:underline"
          >
            返回控制台（可先浏览各功能模块）
          </Link>
        </div>

        {/* Step indicator */}
        <div className="mb-10 flex items-center justify-center gap-0">
          {steps.map((step, index) => {
            const isCompleted = index < currentStepIndex;
            const isCurrent = index === currentStepIndex;

            return (
              <div key={step.path} className="flex items-center">
                <div className="flex flex-col items-center">
                  <div
                    className={cn(
                      "flex h-8 w-8 items-center justify-center rounded-full text-sm font-medium transition-colors",
                      isCompleted
                        ? "bg-green-100 text-green-700"
                        : isCurrent
                          ? "bg-blue-700 text-white"
                          : "bg-slate-200 text-slate-500"
                    )}
                  >
                    {isCompleted ? (
                      <CheckCircle2 className="h-5 w-5" />
                    ) : (
                      index + 1
                    )}
                  </div>
                  <span
                    className={cn(
                      "mt-2 text-xs font-medium",
                      isCurrent ? "text-blue-700" : "text-slate-500"
                    )}
                  >
                    {step.label}
                  </span>
                </div>
                {index < steps.length - 1 && (
                  <div
                    className={cn(
                      "mx-4 mb-5 h-0.5 w-16 sm:w-24",
                      index < currentStepIndex ? "bg-green-400" : "bg-slate-200"
                    )}
                  />
                )}
              </div>
            );
          })}
        </div>

        {/* Content */}
        {children}
      </div>
    </div>
  );
}
