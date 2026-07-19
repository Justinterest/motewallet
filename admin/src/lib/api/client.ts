import axios from "axios";
import type { ApiResponse } from "@/types/api";

// 生产：留空走同域（nginx 把 /api/ 转到 backend）；本地开发默认直连 :8080
const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080",
  withCredentials: true,
  timeout: 15000,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.response.use(
  (response) => {
    const res = response.data as ApiResponse<unknown>;

    if (res.code !== 0) {
      // 业务错误
      return Promise.reject(new Error(res.message || "请求失败"));
    }

    // 解包，返回 data 字段
    return res.data as ReturnType<typeof response.data>;
  },
  (error) => {
    if (error.response?.status === 401) {
      // 未授权，跳转登录页
      if (typeof window !== "undefined") {
        window.location.href = "/login";
      }
      return Promise.reject(new Error("未授权，请重新登录"));
    }

    const message =
      error.response?.data?.message || error.message || "网络请求失败";
    return Promise.reject(new Error(message));
  }
);

export default apiClient;
