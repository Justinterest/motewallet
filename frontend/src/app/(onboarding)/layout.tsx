"use client";

import { useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Skeleton } from "@/components/ui/skeleton";
import { useCurrentUser } from "@/lib/hooks/use-auth";

export default function OnboardingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const { isLoading, isError } = useCurrentUser();

  useEffect(() => {
    if (isError) {
      router.push("/login");
    }
  }, [isError, router]);

  if (isLoading) {
    return (
      <div className="auth-shell flex min-h-screen flex-col items-center px-4 pt-16">
        <div className="w-full max-w-2xl">
          <Skeleton className="mx-auto mb-8 h-8 w-48" />
          <Skeleton className="h-80 w-full rounded-lg" />
        </div>
      </div>
    );
  }

  if (isError) {
    return null;
  }

  return (
    <div className="auth-shell flex min-h-screen flex-col items-center px-4 pt-16">
      <div className="w-full max-w-2xl">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold text-primary">Motewallet</h1>
          <p className="mt-1 text-sm text-muted-foreground">签署服务协议</p>
          <Link
            href="/dashboard"
            className="mt-2 inline-block text-sm text-brand-sky hover:underline"
          >
            返回控制台
          </Link>
        </div>
        {children}
      </div>
    </div>
  );
}
