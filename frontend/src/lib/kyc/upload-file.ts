import { kycFilesApi } from "@/lib/api/kyc-files";

const MAX_FILE_BYTES = 10 * 1024 * 1024;

const ALLOWED_TYPES = new Set([
  "application/pdf",
  "image/jpeg",
  "image/jpg",
  "image/png",
  "image/gif",
  "image/bmp",
  "image/webp",
]);

export interface StagedKycFile {
  objectKey: string;
  accessUrl: string;
}

export function validateKycFile(file: File): string | null {
  if (!ALLOWED_TYPES.has(file.type)) {
    return "仅支持 PDF、JPEG、PNG、GIF、BMP、WEBP";
  }
  if (file.size > MAX_FILE_BYTES) {
    return "文件不能超过 10 MB";
  }
  return null;
}

export function isStagedObjectKey(path: string): boolean {
  return Boolean(path) && !path.startsWith("http://") && !path.startsWith("https://");
}

/** Fetch a presigned GET URL for a staged S3 object key. */
export async function getKycFileAccessUrl(objectKey: string): Promise<string> {
  const access = await kycFilesApi.access({ object_key: objectKey });
  return access.access_url;
}

/**
 * Stage a KYC file in S3 via presigned upload, then return object key and presigned access URL.
 * KUN upload happens when the full KYC form is submitted.
 */
export async function uploadKycFile(file: File): Promise<StagedKycFile> {
  const validationError = validateKycFile(file);
  if (validationError) {
    throw new Error(validationError);
  }

  const presign = await kycFilesApi.presign({
    filename: file.name,
    content_type: file.type,
  });

  const putRes = await fetch(presign.upload_url, {
    method: "PUT",
    headers: { "Content-Type": file.type },
    body: file,
  });
  if (!putRes.ok) {
    throw new Error("上传到存储失败，请重试");
  }

  const accessUrl = await getKycFileAccessUrl(presign.object_key);

  return {
    objectKey: presign.object_key,
    accessUrl,
  };
}
