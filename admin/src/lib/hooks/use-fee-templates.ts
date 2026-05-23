import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { feeTemplateApi } from "@/lib/api/fee-templates";
import type { CreateFeeTemplateRequest, UpdateFeeTemplateRequest } from "@/types/fee-template";

export function useFeeTemplates() {
  return useQuery({
    queryKey: ["fee-templates"],
    queryFn: () => feeTemplateApi.list(),
  });
}

export function useFeeTemplate(id: number) {
  return useQuery({
    queryKey: ["fee-templates", id],
    queryFn: () => feeTemplateApi.getById(id),
    enabled: !!id,
  });
}

export function useCreateFeeTemplate() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateFeeTemplateRequest) => feeTemplateApi.create(data),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["fee-templates"] }),
  });
}

export function useUpdateFeeTemplate() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateFeeTemplateRequest }) =>
      feeTemplateApi.update(id, data),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["fee-templates"] }),
  });
}

export function useDeleteFeeTemplate() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => feeTemplateApi.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["fee-templates"] }),
  });
}
