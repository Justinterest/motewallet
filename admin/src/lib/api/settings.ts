import apiClient from "./client";
import type { SystemCurrencyConfig, UpdateSystemCurrencyConfigPayload } from "@/types/currency-config";

export const settingsApi = {
  getCurrencyConfig: () =>
    apiClient.get<never, SystemCurrencyConfig>("/api/v1/admin/settings/currencies"),
  updateCurrencyConfig: (data: UpdateSystemCurrencyConfigPayload) =>
    apiClient.put<never, void>("/api/v1/admin/settings/currencies", data),
};
