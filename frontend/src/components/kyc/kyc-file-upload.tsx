"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { FileText, ImageIcon, Loader2, Trash2, Upload, X, ZoomIn } from "lucide-react";
import { FieldExampleLink, FieldHint } from "@/components/kyc/field-hint";
import { Label } from "@/components/ui/label";
import {
  getKycFileAccessUrl,
  isStagedObjectKey,
  uploadKycFile,
  validateKycFile,
} from "@/lib/kyc/upload-file";
import { cn } from "@/lib/utils";
import { toast } from "@/hooks/use-toast";

interface FilePreviewMeta {
  name: string;
  /** Presigned S3 GET URL for preview and display */
  accessUrl: string;
  previewUrl: string | null;
  isImage: boolean;
  isPdf: boolean;
}

interface PendingUpload {
  id: string;
  name: string;
  previewUrl: string | null;
  isImage: boolean;
  isPdf: boolean;
}

interface KycFileUploadProps {
  label: string;
  /** 显示在 label 下方的补充说明 */
  description?: string;
  /** @deprecated 使用 description */
  hint?: string;
  exampleImage?: string;
  paths: string[];
  onChange: (paths: string[]) => void;
  minItems?: number;
  accept?: string;
  disabled?: boolean;
}

function isImageFile(file: File): boolean {
  return file.type.startsWith("image/");
}

function isPdfFile(file: File): boolean {
  return file.type === "application/pdf";
}

function buildPreviewFromFile(file: File): Pick<FilePreviewMeta, "previewUrl" | "isImage" | "isPdf"> {
  const isImage = isImageFile(file);
  const isPdf = isPdfFile(file);
  return {
    isImage,
    isPdf,
    previewUrl: isImage ? URL.createObjectURL(file) : null,
  };
}

/** S3 object keys use UUID filenames; hide those from the user. */
function displayFileName(name: string): string {
  if (/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\./i.test(name)) {
    return "已上传文件";
  }
  return name;
}

function inferFileKind(filenameOrKey: string): Pick<FilePreviewMeta, "isImage" | "isPdf"> {
  const probe = filenameOrKey.split("/").pop() ?? filenameOrKey;
  const isPdf = /\.pdf$/i.test(probe);
  const isImage = !isPdf && /\.(jpe?g|png|gif|bmp|webp)$/i.test(probe);
  return { isImage, isPdf };
}

function buildMetaFromAccess(
  name: string,
  accessUrl: string,
  objectKey?: string
): FilePreviewMeta {
  const { isImage, isPdf } = inferFileKind(objectKey ?? name);
  return {
    name: displayFileName(name),
    accessUrl,
    previewUrl: isImage ? accessUrl : null,
    isImage,
    isPdf,
  };
}

function revokeBlob(url: string | null) {
  if (url?.startsWith("blob:")) {
    URL.revokeObjectURL(url);
  }
}

export function KycFileUpload({
  label,
  description,
  hint,
  exampleImage,
  paths,
  onChange,
  minItems = 1,
  accept = ".pdf,.jpg,.jpeg,.png,.gif,.bmp,.webp,application/pdf,image/*",
  disabled = false,
}: KycFileUploadProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const dragCounterRef = useRef(0);
  const fileMetaRef = useRef<Record<string, FilePreviewMeta>>({});
  const pathsRef = useRef(paths);

  const [isDragging, setIsDragging] = useState(false);
  const [uploadingCount, setUploadingCount] = useState(0);
  const [fileMeta, setFileMeta] = useState<Record<string, FilePreviewMeta>>({});
  const [pending, setPending] = useState<PendingUpload[]>([]);
  const [lightboxUrl, setLightboxUrl] = useState<string | null>(null);

  fileMetaRef.current = fileMeta;
  pathsRef.current = paths;

  const filledPaths = paths.filter(Boolean);
  const isUploading = uploadingCount > 0;
  const remaining = Math.max(0, minItems - filledPaths.length - pending.length);

  const commitPaths = useCallback(
    (nextFilled: string[]) => {
      if (nextFilled.length === 0) {
        onChange(minItems > 0 ? Array(minItems).fill("") : []);
        return;
      }
      onChange(nextFilled);
    },
    [minItems, onChange]
  );

  const uploadOne = useCallback(
    async (file: File) => {
      const validationError = validateKycFile(file);
      if (validationError) {
        toast({
          variant: "destructive",
          title: "无法上传",
          description: validationError,
        });
        return;
      }

      const pendingId = crypto.randomUUID();
      const localPreview = buildPreviewFromFile(file);

      setPending((prev) => [
        ...prev,
        {
          id: pendingId,
          name: file.name,
          previewUrl: localPreview.previewUrl,
          isImage: localPreview.isImage,
          isPdf: localPreview.isPdf,
        },
      ]);
      setUploadingCount((c) => c + 1);

      try {
        const staged = await uploadKycFile(file);
        revokeBlob(localPreview.previewUrl);
        const meta = buildMetaFromAccess(file.name, staged.accessUrl, staged.objectKey);

        setFileMeta((prev) => ({ ...prev, [staged.objectKey]: meta }));
        const nextFilled = [...pathsRef.current.filter(Boolean), staged.objectKey];
        pathsRef.current = nextFilled;
        commitPaths(nextFilled);
        toast({ title: "上传成功", description: file.name });
      } catch (error) {
        toast({
          variant: "destructive",
          title: "上传失败",
          description: error instanceof Error ? error.message : "请重试",
        });
      } finally {
        setPending((prev) => {
          const item = prev.find((p) => p.id === pendingId);
          revokeBlob(item?.previewUrl ?? null);
          return prev.filter((p) => p.id !== pendingId);
        });
        setUploadingCount((c) => c - 1);
        if (inputRef.current) inputRef.current.value = "";
      }
    },
    [commitPaths]
  );

  const processFileList = useCallback(
    async (fileList: FileList | null) => {
      if (!fileList?.length || disabled) return;
      const files = Array.from(fileList);
      for (const file of files) {
        await uploadOne(file);
      }
    },
    [disabled, uploadOne]
  );

  function removeAt(index: number) {
    const path = filledPaths[index];
    if (!path) return;

    const meta = fileMeta[path];
    revokeBlob(meta?.previewUrl ?? null);
    setFileMeta((prev) => {
      const next = { ...prev };
      delete next[path];
      return next;
    });
    const nextFilled = filledPaths.filter((_, i) => i !== index);
    pathsRef.current = nextFilled;
    commitPaths(nextFilled);
  }

  useEffect(() => {
    let cancelled = false;

    for (const path of filledPaths) {
      if (!isStagedObjectKey(path)) continue;

      void getKycFileAccessUrl(path)
        .then((accessUrl) => {
          if (cancelled) return;
          const rawName = path.split("/").pop() ?? "";
          setFileMeta((prev) => {
            if (prev[path]?.accessUrl) return prev;
            return {
              ...prev,
              [path]: buildMetaFromAccess(rawName || "已上传文件", accessUrl, path),
            };
          });
        })
        .catch(() => {
          /* ignore; user can re-upload */
        });
    }

    setFileMeta((prev) => {
      const next = { ...prev };
      let changed = false;
      for (const key of Object.keys(next)) {
        if (!filledPaths.includes(key)) {
          revokeBlob(next[key].previewUrl);
          delete next[key];
          changed = true;
        }
      }
      return changed ? next : prev;
    });

    return () => {
      cancelled = true;
    };
  }, [filledPaths.join("|")]);

  useEffect(() => {
    return () => {
      for (const meta of Object.values(fileMetaRef.current)) {
        revokeBlob(meta.previewUrl);
      }
    };
  }, []);

  const handleDragEnter = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (disabled || isUploading) return;
    dragCounterRef.current += 1;
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    dragCounterRef.current -= 1;
    if (dragCounterRef.current <= 0) {
      dragCounterRef.current = 0;
      setIsDragging(false);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    dragCounterRef.current = 0;
    setIsDragging(false);
    void processFileList(e.dataTransfer.files);
  };

  const hasPreviews = filledPaths.length > 0 || pending.length > 0;
  const inlineLayout = minItems > 1 || hasPreviews;
  const uploadDescription = description ?? hint;

  const dropZone = (
    <div
      role="button"
      tabIndex={disabled || isUploading ? -1 : 0}
      aria-label="上传文件，支持拖放"
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          inputRef.current?.click();
        }
      }}
      onClick={() => {
        if (!disabled && !isUploading) inputRef.current?.click();
      }}
      onDragEnter={handleDragEnter}
      onDragLeave={handleDragLeave}
      onDragOver={handleDragOver}
      onDrop={handleDrop}
      className={cn(
        "relative flex shrink-0 cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed text-center transition-colors",
        inlineLayout
          ? "h-[7.75rem] w-[6.75rem] px-2 py-3 sm:w-28"
          : "min-h-[140px] w-full px-4 py-8",
        isDragging && !disabled && !isUploading
          ? "border-blue-500 bg-blue-50/80"
          : "border-slate-200 bg-slate-50/50 hover:border-slate-300 hover:bg-slate-50",
        (disabled || isUploading) && "pointer-events-none cursor-not-allowed opacity-60"
      )}
    >
      {isUploading ? (
        <Loader2
          className={cn(
            "animate-spin text-blue-600",
            inlineLayout ? "h-5 w-5" : "h-8 w-8"
          )}
        />
      ) : (
        <Upload
          className={cn("text-slate-400", inlineLayout ? "h-5 w-5" : "h-8 w-8")}
        />
      )}
      <p
        className={cn(
          "font-medium text-slate-700",
          inlineLayout ? "mt-1.5 text-[10px] leading-tight" : "mt-3 text-sm"
        )}
      >
        {isUploading
          ? "上传中…"
          : isDragging
            ? "松开上传"
            : inlineLayout
              ? "添加文件"
              : "拖放文件到此处，或点击选择"}
      </p>
      {!inlineLayout && (
        <p className="mt-1 text-xs text-slate-500">
          支持 PDF、JPG、PNG 等，单文件不超过 10 MB
        </p>
      )}
      {minItems > 1 && remaining > 0 && !isUploading && (
        <p
          className={cn(
            "font-medium text-amber-700",
            inlineLayout ? "mt-1 text-[10px]" : "mt-2 text-xs"
          )}
        >
          还需 {remaining} 个
        </p>
      )}
    </div>
  );

  return (
    <div className="space-y-3">
      <div className="flex flex-wrap items-center gap-x-2 gap-y-0.5">
        <Label className="mb-0">{label}</Label>
        <FieldExampleLink href={exampleImage} />
      </div>
      <FieldHint meta={{ description: uploadDescription }} />
      <p className="text-xs text-slate-500">
        支持 PDF、JPG、PNG 等，单文件不超过 10 MB
      </p>

      <div className={cn(inlineLayout && "flex flex-wrap items-start gap-3")}>
        {hasPreviews && (
          <ul className="flex flex-wrap gap-3">
            {pending.map((item) => (
              <li key={item.id} className="w-[6.75rem] shrink-0 sm:w-28">
                <PreviewCard
                  compact
                  name={item.name}
                  previewUrl={item.previewUrl}
                  isImage={item.isImage}
                  isPdf={item.isPdf}
                  onPreview={setLightboxUrl}
                  uploading
                />
              </li>
            ))}
            {filledPaths.map((path, index) => {
              const meta = fileMeta[path];
              const kind = meta ?? inferFileKind(path);
              return (
                <li key={path} className="w-[6.75rem] shrink-0 sm:w-28">
                  <PreviewCard
                    compact
                    name={meta?.name ?? displayFileName(path.split("/").pop() ?? "已上传文件")}
                    accessUrl={meta?.accessUrl}
                    previewUrl={meta?.previewUrl ?? null}
                    isImage={kind.isImage}
                    isPdf={kind.isPdf}
                    onPreview={setLightboxUrl}
                    onRemove={() => removeAt(index)}
                    removeDisabled={disabled || isUploading}
                  />
                </li>
              );
            })}
          </ul>
        )}
        {dropZone}
      </div>

      <input
        ref={inputRef}
        type="file"
        className="hidden"
        accept={accept}
        multiple
        disabled={disabled || isUploading}
        onChange={(e) => void processFileList(e.target.files)}
      />

      {lightboxUrl && (
        <ImageLightbox url={lightboxUrl} onClose={() => setLightboxUrl(null)} />
      )}
    </div>
  );
}

function ImageLightbox({ url, onClose }: { url: string; onClose: () => void }) {
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === "Escape") onClose();
    }
    document.addEventListener("keydown", onKey);
    const prev = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => {
      document.removeEventListener("keydown", onKey);
      document.body.style.overflow = prev;
    };
  }, [onClose]);

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/75 p-4"
      role="dialog"
      aria-modal="true"
      aria-label="图片预览"
      onClick={onClose}
    >
      <button
        type="button"
        onClick={onClose}
        className="absolute right-4 top-4 rounded-full bg-black/50 p-2 text-white hover:bg-black/70"
        aria-label="关闭预览"
      >
        <X className="h-5 w-5" />
      </button>
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src={url}
        alt="预览"
        className="max-h-[90vh] max-w-full rounded-md object-contain shadow-2xl"
        onClick={(e) => e.stopPropagation()}
      />
    </div>
  );
}

interface PreviewCardProps {
  name: string;
  accessUrl?: string;
  previewUrl: string | null;
  isImage: boolean;
  isPdf: boolean;
  compact?: boolean;
  uploading?: boolean;
  onPreview?: (url: string) => void;
  onRemove?: () => void;
  removeDisabled?: boolean;
}

function PreviewCard({
  name,
  accessUrl,
  previewUrl,
  isImage,
  isPdf,
  compact,
  uploading,
  onPreview,
  onRemove,
  removeDisabled,
}: PreviewCardProps) {
  const displayUrl = previewUrl ?? accessUrl ?? null;
  const imageLoading = isImage && !displayUrl && !uploading;

  return (
    <div className="group relative overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm">
      <div
        className={cn(
          "relative bg-slate-100",
          compact ? "aspect-square" : "aspect-[4/3]"
        )}
      >
        {isImage && displayUrl ? (
          <button
            type="button"
            className="relative block h-full w-full cursor-zoom-in"
            onClick={() => onPreview?.(displayUrl)}
            aria-label={`预览 ${name}`}
          >
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={displayUrl}
              alt={name}
              className="h-full w-full object-cover"
            />
            <span className="absolute inset-0 flex items-center justify-center bg-black/0 opacity-0 transition-opacity group-hover:bg-black/25 group-hover:opacity-100">
              <ZoomIn className="h-6 w-6 text-white drop-shadow" />
            </span>
          </button>
        ) : isImage && imageLoading ? (
          <div className="flex h-full flex-col items-center justify-center gap-1.5 p-2 text-slate-400">
            <Loader2 className="h-5 w-5 animate-spin text-blue-600" />
            <span className={cn("text-slate-500", compact ? "text-[10px]" : "text-xs")}>
              加载预览…
            </span>
          </div>
        ) : (
          <div className="flex h-full flex-col items-center justify-center gap-2 p-3 text-slate-500">
            {isPdf ? (
              <FileText
                className={cn(
                  "text-red-500/80",
                  compact ? "h-7 w-7" : "h-10 w-10"
                )}
              />
            ) : (
              <ImageIcon className={compact ? "h-7 w-7" : "h-10 w-10"} />
            )}
            {!compact && (
              <>
                <span className="text-xs font-medium uppercase tracking-wide text-slate-400">
                  {isPdf ? "PDF" : "文件"}
                </span>
                {isPdf && accessUrl && (
                  <a
                    href={accessUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    onClick={(e) => e.stopPropagation()}
                    className="text-xs text-blue-700 hover:underline"
                  >
                    在新标签页打开
                  </a>
                )}
              </>
            )}
          </div>
        )}

        {uploading && (
          <div className="absolute inset-0 flex items-center justify-center bg-white/70">
            <Loader2 className="h-6 w-6 animate-spin text-blue-600" />
          </div>
        )}

        {onRemove && !uploading && (
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              onRemove();
            }}
            disabled={removeDisabled}
            aria-label={`删除 ${name}`}
            className="absolute right-1.5 top-1.5 rounded-md bg-black/55 p-1.5 text-white opacity-0 transition-opacity hover:bg-black/70 group-hover:opacity-100 focus:opacity-100 disabled:opacity-40"
          >
            <Trash2 className="h-3.5 w-3.5" />
          </button>
        )}
      </div>
      <p
        className={cn(
          "truncate px-1.5 text-slate-600",
          compact ? "py-1 text-[10px]" : "py-1.5 text-xs"
        )}
      >
        {name}
      </p>
    </div>
  );
}
