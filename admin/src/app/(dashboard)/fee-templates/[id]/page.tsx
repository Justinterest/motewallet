"use client";

import { use } from "react";
import { useRouter } from "next/navigation";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { FeeTemplateForm } from "@/components/fee-templates/fee-template-form";
import { useFeeTemplate, useUpdateFeeTemplate } from "@/lib/hooks/use-fee-templates";
import { toast } from "@/hooks/use-toast";
import type { UpdateFeeTemplateRequest } from "@/types/fee-template";

export default function FeeTemplateDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id: idStr } = use(params);
  const id = parseInt(idStr, 10);
  const router = useRouter();
  const { data: template, isLoading } = useFeeTemplate(id);
  const updateMutation = useUpdateFeeTemplate();

  const handleSubmit = (data: UpdateFeeTemplateRequest) => {
    updateMutation.mutate(
      { id, data },
      {
        onSuccess: () => {
          toast({ title: "保存成功" });
          router.push("/fee-templates");
        },
        onError: (error) => {
          toast({ title: "保存失败", description: error.message, variant: "destructive" });
        },
      }
    );
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (!template) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/fee-templates">
              <ArrowLeft className="size-4" />
            </Link>
          </Button>
          <h1 className="text-2xl font-bold text-slate-900">模板不存在</h1>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" asChild>
          <Link href="/fee-templates">
            <ArrowLeft className="size-4" />
          </Link>
        </Button>
        <h1 className="text-2xl font-bold text-slate-900">编辑手续费模板</h1>
      </div>

      <FeeTemplateForm
        defaultValues={{
          name: template.name,
          description: template.description || "",
          is_default: template.is_default,
        }}
        defaultExchangeItems={template.exchange_items}
        defaultCryptoItems={template.crypto_withdrawal_items}
        defaultFiatItems={template.fiat_withdrawal_items}
        defaultExchangeFeeDeductionMethod={template.exchange_fee_deduction_method}
        defaultCryptoFeeDeductionMethod={template.crypto_withdrawal_fee_deduction_method}
        defaultFiatFeeDeductionMethod={template.fiat_withdrawal_fee_deduction_method}
        onSubmit={handleSubmit}
        isPending={updateMutation.isPending}
        submitLabel="保存修改"
      />
    </div>
  );
}
