"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { TopNav } from "@/components/layout/top-nav";
import { VerificationBanner } from "@/components/layout/verification-banner";
import { VerificationGuard } from "@/components/layout/verification-guard";
import { VerificationProvider } from "@/components/layout/verification-context";
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

  if (isLoading) {
    return (
      <div className="app-shell">
        <div className="fixed top-0 z-50 w-full border-b border-border bg-background/80 backdrop-blur-sm">
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

  return (
    <VerificationProvider user={user}>
      <div className="app-shell">
        <TopNav />
        <main className="pt-14">
          <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
            <VerificationBanner />
            <VerificationGuard>{children}</VerificationGuard>
          </div>
        </main>
      </div>
    </VerificationProvider>
  );
}
