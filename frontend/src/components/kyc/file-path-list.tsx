"use client";

import { Plus, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

interface FilePathListProps {
  label: string;
  description?: string;
  paths: string[];
  onChange: (paths: string[]) => void;
  minItems?: number;
}

export function FilePathList({
  label,
  description,
  paths,
  onChange,
  minItems = 1,
}: FilePathListProps) {
  function updateAt(index: number, value: string) {
    const next = [...paths];
    next[index] = value;
    onChange(next);
  }

  function add() {
    onChange([...paths, ""]);
  }

  function remove(index: number) {
    if (paths.length <= minItems) return;
    onChange(paths.filter((_, i) => i !== index));
  }

  return (
    <div className="space-y-2">
      <Label>{label}</Label>
      {description && (
        <p className="text-xs text-slate-500">{description}</p>
      )}
      <div className="space-y-2">
        {paths.map((path, index) => (
          <div key={index} className="flex gap-2">
            <Input
              value={path}
              onChange={(e) => updateAt(index, e.target.value)}
              placeholder="鲲文件上传接口返回的 path"
              className="flex-1"
            />
            <Button
              type="button"
              variant="outline"
              size="icon"
              onClick={() => remove(index)}
              disabled={paths.length <= minItems}
              aria-label="删除"
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        ))}
      </div>
      <Button type="button" variant="outline" size="sm" onClick={add}>
        <Plus className="mr-1 h-4 w-4" />
        添加文件
      </Button>
    </div>
  );
}
