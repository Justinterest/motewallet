"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { TopNav } from "@/components/layout/top-nav";
import { Skeleton } from "@/components/ui/skeleton";
import { useCurrentUser } from "@/lib/hooks/use-auth";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const { data: user, isLoading, isError } = useCurrentUser();

  useEffect(() => {
    if (isError) {
      router.push("/login");
    }
  }, [isError, router]);

  useEffect(() => {
    if (!user) return;

    if (user.status === "PENDING_AGREEMENT") {
      router.replace("/onboarding/agreement");
    } else if (user.status === "PENDING_KYC") {
      router.replace("/onboarding/kyc");
    } else if (user.status === "PROCESSING") {
      router.replace("/onboarding/status");
    }
  }, [user, router]);

  if (isLoading) {
    return (
      <div className="min-h-screen bg-slate-50">
        <div className="fixed top-0 z-50 w-full border-b bg-white">
          <div className="mx-auto flex h-14 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
            <Skeleton className="h-6 w-28" />
            <Skeleton className="h-8 w-8 rounded-full" />
          </div>
        </div>
        <main className="pt-14">
          <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
            <Skeleton className="mb-6 h-8 w-32" />
            <div className="grid gap-6 md:grid-cols-2">
              <Skeleton className="h-40 rounded-lg" />
              <Skeleton className="h-40 rounded-lg" />
            </div>
          </div>
        </main>
      </div>
    );
  }

  if (isError || !user) {
    return null;
  }

  // If user is not ACTIVE, the redirect useEffect will handle it
  if (user.status !== "ACTIVE") {
    return null;
  }

  return (
    <div className="min-h-screen bg-slate-50">
      <TopNav />
      <main className="pt-14">
        <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
          {children}
        </div>
      </main>
    </div>
  );
}
