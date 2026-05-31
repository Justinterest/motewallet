"use client";

import type { KycFieldMeta } from "@/lib/kyc/field-meta";

interface FieldExampleLinkProps {
  href?: string;
  className?: string;
}

/** Shown inline next to a field label; opens the example asset in a new tab. */
export function FieldExampleLink({ href, className }: FieldExampleLinkProps) {
  if (!href) return null;

  return (
    <a
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      className={
        className ??
        "text-xs font-normal text-blue-700 hover:text-blue-800 hover:underline"
      }
    >
      查看示例
    </a>
  );
}

interface FieldHintProps {
  meta: Pick<KycFieldMeta, "description">;
  className?: string;
}

export function FieldHint({ meta, className }: FieldHintProps) {
  const { description } = meta;
  if (!description) return null;

  return (
    <p className={className ?? "text-xs leading-relaxed text-slate-500"}>
      {description}
    </p>
  );
}
