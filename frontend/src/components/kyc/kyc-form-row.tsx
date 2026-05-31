import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

/** Spans full width of parent 2-col grid; children form a 2-column row on sm+. */
export function KycFormRow({
  children,
  className,
}: {
  children: ReactNode;
  className?: string;
}) {
  return (
    <div className={cn("grid gap-4 sm:col-span-2 sm:grid-cols-2", className)}>
      {children}
    </div>
  );
}
