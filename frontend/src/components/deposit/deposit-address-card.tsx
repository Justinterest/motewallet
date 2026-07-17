"use client";

import { useState } from "react";
import { QRCodeSVG } from "qrcode.react";
import { Check, Copy } from "lucide-react";
import { Button } from "@/components/ui/button";

interface DepositAddressCardProps {
  address: string;
  currency: string;
  network?: string;
}

export function DepositAddressCard({
  address,
  currency,
  network,
}: DepositAddressCardProps) {
  const [copied, setCopied] = useState(false);

  async function handleCopy() {
    try {
      await navigator.clipboard.writeText(address);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // Fallback for older browsers / insecure contexts
      const textarea = document.createElement("textarea");
      textarea.value = address;
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand("copy");
      document.body.removeChild(textarea);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  }

  return (
    <div className="space-y-4 rounded-lg border bg-slate-50 p-4">
      <div className="flex justify-center">
        <div className="rounded-xl bg-white p-3 shadow-sm ring-1 ring-slate-200">
          <QRCodeSVG value={address} size={168} level="M" includeMargin={false} />
        </div>
      </div>
      <p className="text-center text-xs text-slate-500">扫码向该地址转入 {currency}</p>

      <div>
        <p className="mb-1 text-xs font-medium text-slate-500">充值地址</p>
        <p className="break-all font-mono text-sm leading-relaxed text-slate-900">
          {address}
        </p>
        {network && (
          <p className="mt-1.5 text-xs text-slate-500">网络：{network}</p>
        )}
      </div>

      <Button
        type="button"
        variant="outline"
        className="w-full"
        onClick={handleCopy}
      >
        {copied ? <Check className="mr-1.5 h-4 w-4" /> : <Copy className="mr-1.5 h-4 w-4" />}
        {copied ? "已复制" : "复制地址"}
      </Button>
    </div>
  );
}
