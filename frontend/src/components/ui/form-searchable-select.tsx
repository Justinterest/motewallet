"use client";

import * as React from "react";
import { CheckIcon, ChevronDownIcon } from "lucide-react";

import { cn } from "@/lib/utils";
import { Input } from "@/components/ui/input";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import type { SelectOption } from "@/components/ui/select";

interface FormSearchableSelectProps {
  value?: string;
  onChange: (value: string) => void;
  onBlur?: () => void;
  options: SelectOption[] | readonly SelectOption[];
  placeholder?: string;
  searchPlaceholder?: string;
  emptyText?: string;
  disabled?: boolean;
  className?: string;
}

function filterOptions(
  options: readonly SelectOption[],
  query: string
): SelectOption[] {
  const q = query.trim().toLowerCase();
  if (!q) return [...options];
  return options.filter(
    (opt) =>
      opt.label.toLowerCase().includes(q) ||
      opt.value.toLowerCase().includes(q)
  );
}

function FormSearchableSelect({
  value,
  onChange,
  onBlur,
  options,
  placeholder = "请选择",
  searchPlaceholder = "搜索…",
  emptyText = "无匹配结果",
  disabled,
  className,
}: FormSearchableSelectProps) {
  const [open, setOpen] = React.useState(false);
  const [search, setSearch] = React.useState("");
  const listId = React.useId();
  const inputRef = React.useRef<HTMLInputElement>(null);

  const selectedLabel = options.find((o) => o.value === value)?.label;
  const filtered = React.useMemo(
    () => filterOptions(options, search),
    [options, search]
  );

  return (
    <Popover
      open={open}
      onOpenChange={(nextOpen) => {
        setOpen(nextOpen);
        if (!nextOpen) {
          setSearch("");
          onBlur?.();
        } else {
          requestAnimationFrame(() => inputRef.current?.focus());
        }
      }}
    >
      <PopoverTrigger asChild>
        <button
          type="button"
          role="combobox"
          aria-expanded={open}
          aria-controls={listId}
          disabled={disabled}
          className={cn(
            "border-input focus-visible:border-ring focus-visible:ring-ring/50 flex h-9 w-full items-center justify-between gap-2 rounded-md border bg-card px-3 py-2 text-sm shadow-xs transition-[color,box-shadow] outline-none focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50",
            !selectedLabel && "text-muted-foreground",
            className
          )}
        >
          <span className="truncate">{selectedLabel ?? placeholder}</span>
          <ChevronDownIcon className="size-4 shrink-0 opacity-50" />
        </button>
      </PopoverTrigger>
      <PopoverContent
        className="w-[var(--radix-popover-trigger-width)] p-0"
        align="start"
      >
        <div className="border-b p-2">
          <Input
            ref={inputRef}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder={searchPlaceholder}
            className="h-8"
          />
        </div>
        <div id={listId} role="listbox" className="max-h-60 overflow-y-auto p-1">
          {filtered.length === 0 ? (
            <p className="text-muted-foreground py-6 text-center text-sm">
              {emptyText}
            </p>
          ) : (
            filtered.map((option) => (
              <button
                key={option.value}
                type="button"
                className={cn(
                  "hover:bg-accent hover:text-accent-foreground relative flex w-full cursor-default items-center rounded-sm py-1.5 pr-8 pl-2 text-sm outline-hidden select-none",
                  value === option.value && "bg-accent text-accent-foreground"
                )}
                onClick={() => {
                  onChange(option.value);
                  setOpen(false);
                }}
              >
                {option.label}
                {value === option.value && (
                  <CheckIcon className="absolute right-2 size-4" />
                )}
              </button>
            ))
          )}
        </div>
      </PopoverContent>
    </Popover>
  );
}

export { FormSearchableSelect };
