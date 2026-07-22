"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { adminUsersApi } from "@/lib/api/users";
import type { CreateAdminEmployeeRequest } from "@/types/auth";

export function useAdminEmployees() {
  return useQuery({
    queryKey: ["admin", "users"],
    queryFn: () => adminUsersApi.list(),
  });
}

export function useCreateAdminEmployee() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateAdminEmployeeRequest) => adminUsersApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    },
  });
}

export function useResetAdminPassword() {
  return useMutation({
    mutationFn: (id: number) => adminUsersApi.resetPassword(id),
  });
}

export function useResetAdmin2FA() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => adminUsersApi.reset2FA(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
    },
  });
}
