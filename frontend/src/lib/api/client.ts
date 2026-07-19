import axios from "axios";
import type { ApiResponse } from "@/types/api";
import { ApiError, extractFieldErrors } from "./api-error";

// 生产：留空走同域（nginx 把 /api/ 转到 backend）；本地开发默认直连 :8080
const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080",
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.response.use(
  (response) => {
    const data = response.data as ApiResponse<unknown>;

    if (data.code !== 0) {
      const fieldErrors = extractFieldErrors(data.data);
      return Promise.reject(
        new ApiError(data.message || "请求失败", {
          code: data.code,
          fieldErrors,
        })
      );
    }

    return data.data as never;
  },
  (error) => {
    if (error.response?.status === 401) {
      if (typeof window !== "undefined" && !window.location.pathname.startsWith("/login")) {
        window.location.href = "/login";
      }
    }

    const responseData = error.response?.data;
    const message =
      responseData?.message || error.message || "网络错误，请稍后重试";
    const fieldErrors = extractFieldErrors(responseData?.data);
    return Promise.reject(
      new ApiError(message, {
        code: responseData?.code,
        fieldErrors,
      })
    );
  }
);

export default apiClient;
