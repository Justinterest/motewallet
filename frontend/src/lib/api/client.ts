import axios from "axios";
import type { ApiResponse } from "@/types/api";

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080",
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.response.use(
  (response) => {
    const data = response.data as ApiResponse<unknown>;

    if (data.code !== 0) {
      return Promise.reject(new Error(data.message || "请求失败"));
    }

    return data.data as never;
  },
  (error) => {
    if (error.response?.status === 401) {
      if (typeof window !== "undefined" && !window.location.pathname.startsWith("/login")) {
        window.location.href = "/login";
      }
    }

    const message =
      error.response?.data?.message || error.message || "网络错误，请稍后重试";
    return Promise.reject(new Error(message));
  }
);

export default apiClient;
