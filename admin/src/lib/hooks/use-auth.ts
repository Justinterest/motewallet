"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { adminAuthApi } from "@/lib/api/auth";
import { useAuthStore } from "@/stores/auth-store";
import type { AdminLoginRequest } from "@/types/auth";

export function useCurrentAdmin() {
  const { setAdmin } = useAuthStore();

  return useQuery({
    queryKey: ["admin", "me"],
    queryFn: async () => {
      const data = await adminAuthApi.getMe();
      setAdmin(data);
      return data;
    },
    retry: false,
    staleTime: 5 * 60 * 1000,
  });
}

export function useAdminLogin() {
  const { setAdmin } = useAuthStore();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: AdminLoginRequest) => adminAuthApi.login(data),
    onSuccess: (data) => {
      setAdmin(data);
      queryClient.setQueryData(["admin", "me"], data);
    },
  });
}

export function useAdminLogout() {
  const { clearAdmin } = useAuthStore();
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: () => adminAuthApi.logout(),
    onSuccess: () => {
      clearAdmin();
      queryClient.clear();
      router.push("/login");
    },
    onError: () => {
      // 即使退出失败也清除本地状态
      clearAdmin();
      queryClient.clear();
      router.push("/login");
    },
  });
}
