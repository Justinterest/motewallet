"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { CheckCircle2, XCircle, Clock } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useKycStatus } from "@/lib/hooks/use-onboarding";

export default function KycStatusPage() {
  const router = useRouter();
  const { data, isLoading, isError } = useKycStatus();

  useEffect(() => {
    if (data && data.kyc_status === "NONE" && data.status === "PENDING_KYC") {
      router.replace("/kyc");
    }
  }, [data, router]);

  if (isLoading) {
    return (
      <div className="mx-auto max-w-2xl">
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-16">
            <Skeleton className="mb-4 h-16 w-16 rounded-full" />
            <Skeleton className="h-6 w-40" />
            <Skeleton className="mt-2 h-4 w-60" />
          </CardContent>
        </Card>
      </div>
    );
  }

  if (isError) {
    return (
      <div className="mx-auto max-w-2xl">
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-16">
            <XCircle className="mb-4 h-16 w-16 text-red-400" />
            <p className="text-sm text-red-600">加载认证状态失败，请刷新页面重试。</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  const kycStatus = data?.kyc_status;

  if (kycStatus === "AUTH_SUC") {
    return (
      <div className="mx-auto max-w-2xl">
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-16">
            <CheckCircle2 className="mb-4 h-16 w-16 text-green-500" />
            <h2 className="text-xl font-semibold text-foreground">认证已通过</h2>
            <p className="mt-2 text-sm text-muted-foreground">
              恭喜！您的企业认证已审核通过，可以开始使用平台服务。
            </p>
            {data?.kyc_completed_at && (
              <p className="mt-1 text-xs text-muted-foreground">
                通过时间：{data.kyc_completed_at}
              </p>
            )}
            <Button onClick={() => router.push("/dashboard")} className="mt-6" size="lg">
              进入概览
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (kycStatus === "AUTHING") {
    return (
      <div className="mx-auto max-w-2xl">
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-16">
            <div className="mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-accent">
              <Clock className="h-8 w-8 animate-pulse text-primary" />
            </div>
            <h2 className="text-xl font-semibold text-foreground">认证审核中</h2>
            <p className="mt-2 text-sm text-muted-foreground">
              您的认证资料已提交至审核机构，请耐心等待审核结果。
            </p>
            <p className="mt-1 text-xs text-muted-foreground">
              状态将每 15 秒自动更新
            </p>
            {data?.kyc_submitted_at && (
              <p className="mt-1 text-xs text-muted-foreground">
                提交时间：{data.kyc_submitted_at}
              </p>
            )}
            <Button variant="outline" onClick={() => router.push("/dashboard")} className="mt-6">
              返回概览
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (kycStatus === "AUTH_FAIL") {
    return (
      <div className="mx-auto max-w-2xl">
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-16">
            <XCircle className="mb-4 h-16 w-16 text-red-500" />
            <h2 className="text-xl font-semibold text-foreground">认证未通过</h2>
            <p className="mt-2 text-sm text-muted-foreground">
              很抱歉，您的认证未通过审核。
            </p>
            {data?.kyc_fail_reason && (
              <div className="mt-3 w-full max-w-sm rounded-md bg-red-50 p-3">
                <p className="text-sm font-medium text-red-700">失败原因：</p>
                <p className="mt-1 text-sm text-red-600">{data.kyc_fail_reason}</p>
              </div>
            )}
            <Button onClick={() => router.push("/kyc")} className="mt-6" size="lg">
              重新提交
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return null;
}
