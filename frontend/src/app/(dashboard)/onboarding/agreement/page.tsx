"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";
import { ExternalLink, Loader2, FileText } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useAgreements, useSignAgreements } from "@/lib/hooks/use-onboarding";
import { toast } from "@/hooks/use-toast";

export default function AgreementPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { data, isLoading, isError } = useAgreements();
  const signMutation = useSignAgreements();

  useEffect(() => {
    if (!data?.signed) {
      return;
    }
    void queryClient.invalidateQueries({ queryKey: ["auth", "me"] });
    router.replace("/kyc");
  }, [data?.signed, queryClient, router]);

  function handleSign() {
    signMutation.mutate(undefined, {
      onSuccess: () => {
        toast({
          title: "签署成功",
          description: "协议已签署，请继续完成实名认证。",
        });
        router.push("/kyc");
      },
      onError: (error) => {
        toast({
          variant: "destructive",
          title: "签署失败",
          description: error.message || "请稍后重试",
        });
      },
    });
  }

  if (isLoading || data?.signed) {
    return (
      <div className="mx-auto w-full max-w-2xl space-y-4">
        <Skeleton className="h-10 w-48" />
        <Skeleton className="h-48 w-full rounded-lg" />
        <Skeleton className="h-12 w-full rounded-lg" />
      </div>
    );
  }

  if (isError) {
    return (
      <Card className="mx-auto w-full max-w-2xl">
        <CardContent className="flex flex-col items-center justify-center py-12">
          <p className="text-sm text-red-600">加载协议失败，请刷新页面重试。</p>
        </CardContent>
      </Card>
    );
  }

  const agreements = data?.agreements ?? [];

  if (agreements.length === 0) {
    return null;
  }

  return (
    <div className="mx-auto w-full max-w-2xl space-y-4">
      <div>
        <h1 className="text-[32px] font-bold tracking-tight text-foreground">
          服务协议
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          请仔细阅读以下协议内容，确认后签署
        </p>
      </div>

      {agreements.map((agreement) => (
          <Card key={agreement.id}>
            <CardHeader className="pb-3">
              <CardTitle className="flex flex-wrap items-center gap-2 text-base">
                <FileText className="h-4 w-4 text-blue-700" />
                {agreement.title}
                {agreement.version && (
                  <span className="text-xs font-normal text-muted-foreground">
                    v{agreement.version}
                  </span>
                )}
                {agreement.required && (
                  <span className="text-xs font-normal text-red-500">
                    *必读
                  </span>
                )}
              </CardTitle>
            </CardHeader>
            <CardContent>
              {agreement.url ? (
                <a
                  href={agreement.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-1.5 text-sm font-medium text-blue-700 hover:text-blue-800 hover:underline"
                >
                  查看协议全文
                  <ExternalLink className="h-3.5 w-3.5" />
                </a>
              ) : (
                <p className="text-sm text-muted-foreground">暂无协议链接</p>
              )}
            </CardContent>
          </Card>
        ))}

      <Button
        onClick={handleSign}
        disabled={signMutation.isPending}
        className="w-full bg-blue-700 hover:bg-blue-800 text-white"
        size="lg"
      >
        {signMutation.isPending && (
          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
        )}
        全部同意并签署
      </Button>
    </div>
  );
}
