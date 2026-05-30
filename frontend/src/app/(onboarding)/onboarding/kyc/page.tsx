"use client";

import { KycWizard } from "@/components/kyc/kyc-wizard";

export default function KycPage() {
  return (
    <div className="space-y-4">
      <div className="text-center">
        <h2 className="text-lg font-semibold text-slate-900">实名认证</h2>
        <p className="mt-1 text-sm text-slate-500">
          分步填写企业、管理人、股东与董事信息，与鲲入网认证接口字段一致
        </p>
      </div>
      <KycWizard />
    </div>
  );
}
