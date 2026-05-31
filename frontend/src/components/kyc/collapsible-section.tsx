"use client";

import { useState, type ReactNode } from "react";
import { ChevronDown } from "lucide-react";
import { cn } from "@/lib/utils";

interface CollapsibleSectionProps {
  title: string;
  children: ReactNode;
  /** Header actions (e.g. delete), clicks do not toggle expand */
  actions?: ReactNode;
  defaultOpen?: boolean;
}

export function CollapsibleSection({
  title,
  children,
  actions,
  defaultOpen = true,
}: CollapsibleSectionProps) {
  const [open, setOpen] = useState(defaultOpen);

  return (
    <div className="rounded-lg border border-slate-200 bg-white">
      <div className="flex items-center justify-between gap-2 px-4 py-3">
        <button
          type="button"
          onClick={() => setOpen((v) => !v)}
          className="flex min-w-0 flex-1 items-center gap-2 rounded-md text-left hover:bg-slate-50"
          aria-expanded={open}
        >
          <ChevronDown
            className={cn(
              "h-4 w-4 shrink-0 text-slate-500 transition-transform",
              open && "rotate-180"
            )}
          />
          <span className="truncate font-medium text-slate-800">{title}</span>
        </button>
        {actions ? (
          <div className="flex shrink-0 items-center" onClick={(e) => e.stopPropagation()}>
            {actions}
          </div>
        ) : null}
      </div>
      {open ? (
        <div className="border-t border-slate-200 px-4 pb-4 pt-2">{children}</div>
      ) : null}
    </div>
  );
}
