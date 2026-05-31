import apiClient from "./client";

export interface PresignKycFileRequest {
  filename: string;
  content_type: string;
}

export interface PresignKycFileResponse {
  upload_url: string;
  object_key: string;
  expires_in: number;
}

export interface PresignKycFileAccessRequest {
  object_key: string;
}

export interface PresignKycFileAccessResponse {
  access_url: string;
  expires_in: number;
}

export const kycFilesApi = {
  presign: (data: PresignKycFileRequest) =>
    apiClient.post<never, PresignKycFileResponse>("/api/v1/onboarding/files/presign", data),
  access: (data: PresignKycFileAccessRequest) =>
    apiClient.post<never, PresignKycFileAccessResponse>(
      "/api/v1/onboarding/files/access",
      data
    ),
};
