"use client";

import { useRouter } from "next/navigation";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { FeeTemplateForm } from "@/components/fee-templates/fee-template-form";
import { useCreateFeeTemplate } from "@/lib/hooks/use-fee-templates";
import { toast } from "@/hooks/use-toast";
import type { CreateFeeTemplateRequest } from "@/types/fee-template";

export default function CreateFeeTemplatePage() {
  const router = useRouter();
  const createMutation = useCreateFeeTemplate();

  const handleSubmit = (data: CreateFeeTemplateRequest) => {
    createMutation.mutate(data, {
      onSuccess: () => {
        toast({ title: "创建成功" });
        router.push("/fee-templates");
      },
      onError: (error) => {
        toast({ title: "创建失败", description: error.message, variant: "destructive" });
      },
    });
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" asChild>
          <Link href="/fee-templates">
            <ArrowLeft className="size-4" />
          </Link>
        </Button>
        <h1 className="text-2xl font-bold text-slate-900">创建手续费模板</h1>
      </div>

      <FeeTemplateForm
        onSubmit={handleSubmit}
        isPending={createMutation.isPending}
        submitLabel="创建模板"
      />
    </div>
  );
}
