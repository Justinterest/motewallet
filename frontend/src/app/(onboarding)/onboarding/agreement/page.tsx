"use client";

import { useRouter } from "next/navigation";
import { Loader2, FileText, ScrollText } from "lucide-react";
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
  const { data, isLoading, isError } = useAgreements();
  const signMutation = useSignAgreements();

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

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-48 w-full rounded-lg" />
        <Skeleton className="h-48 w-full rounded-lg" />
        <Skeleton className="h-12 w-full rounded-lg" />
      </div>
    );
  }

  if (isError) {
    return (
      <Card>
        <CardContent className="flex flex-col items-center justify-center py-12">
          <p className="text-sm text-red-600">加载协议失败，请刷新页面重试。</p>
        </CardContent>
      </Card>
    );
  }

  const agreements = data?.agreements || [];

  return (
    <div className="space-y-4">
      <div className="text-center">
        <h2 className="text-lg font-semibold text-slate-900">服务协议</h2>
        <p className="mt-1 text-sm text-slate-500">
          请仔细阅读以下协议内容，确认后签署
        </p>
      </div>

      {agreements.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <ScrollText className="mb-3 h-10 w-10 text-slate-300" />
            <p className="text-sm text-slate-400">暂无需要签署的协议</p>
          </CardContent>
        </Card>
      ) : (
        agreements.map((agreement) => (
          <Card key={agreement.id}>
            <CardHeader className="pb-3">
              <CardTitle className="flex items-center gap-2 text-base">
                <FileText className="h-4 w-4 text-blue-700" />
                {agreement.title}
                {agreement.required && (
                  <span className="text-xs font-normal text-red-500">
                    *必读
                  </span>
                )}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="max-h-60 overflow-y-auto rounded-md bg-slate-50 p-4 text-sm leading-relaxed text-slate-700 whitespace-pre-wrap">
                {agreement.content}
              </div>
            </CardContent>
          </Card>
        ))
      )}

      <Button
        onClick={handleSign}
        disabled={signMutation.isPending || agreements.length === 0}
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
