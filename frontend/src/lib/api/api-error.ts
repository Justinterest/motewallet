export interface KycFieldErrorData {
  errors: string[];
}

export class ApiError extends Error {
  code?: number;
  fieldErrors?: string[];

  constructor(message: string, options?: { code?: number; fieldErrors?: string[] }) {
    super(message);
    this.name = "ApiError";
    this.code = options?.code;
    this.fieldErrors = options?.fieldErrors;
  }
}

export function extractFieldErrors(data: unknown): string[] | undefined {
  if (!data || typeof data !== "object") return undefined;
  const errors = (data as KycFieldErrorData).errors;
  return Array.isArray(errors) && errors.length > 0 ? errors : undefined;
}
