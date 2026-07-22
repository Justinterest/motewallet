"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { adminAuthApi } from "@/lib/api/auth";
import { useAuthStore } from "@/stores/auth-store";
import type {
  AdminAuthChallenge,
  AdminChangePasswordRequest,
  AdminLoginRequest,
  AdminUser,
} from "@/types/auth";

export function completeAdminAuth(
  challenge: AdminAuthChallenge,
  setAdmin: (admin: AdminUser) => void,
  queryClient: ReturnType<typeof useQueryClient>,
  router: ReturnType<typeof useRouter>
) {
  if (challenge.status !== "SUCCESS" || !challenge.admin) {
    return challenge;
  }
  setAdmin(challenge.admin);
  queryClient.setQueryData(["admin", "me"], challenge.admin);
  router.push("/dashboard");
  return challenge;
}

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
  return useMutation({
    mutationFn: (data: AdminLoginRequest) => adminAuthApi.login(data),
  });
}

export function useVerifyAdmin2FA() {
  return useMutation({
    mutationFn: (data: { temp_token: string; code: string }) =>
      adminAuthApi.verify2FA(data),
  });
}

export function useConfirmAdmin2FASetup() {
  return useMutation({
    mutationFn: (data: { temp_token: string; code: string }) =>
      adminAuthApi.confirm2FASetup(data),
  });
}

export function useChangeAdminPassword() {
  return useMutation({
    mutationFn: (data: AdminChangePasswordRequest) =>
      adminAuthApi.changePassword(data),
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
      clearAdmin();
      queryClient.clear();
      router.push("/login");
    },
  });
}
