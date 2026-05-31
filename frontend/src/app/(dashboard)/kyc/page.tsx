"use client";

import { KycWizard } from "@/components/kyc/kyc-wizard";

export default function KycPage() {
  return (
    <div className="mx-auto w-full max-w-5xl space-y-4">
      <div>
        <h1 className="text-[32px] font-bold tracking-tight text-foreground">
          实名认证
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          分步填写企业、管理人、股东及董事信息
        </p>
      </div>
      <KycWizard />
    </div>
  );
}
