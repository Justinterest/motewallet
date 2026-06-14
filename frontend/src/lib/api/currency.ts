import apiClient from "./client";
import type { SupportedCurrencies } from "@/types/currency";

export const currencyApi = {
  getSupported: () =>
    apiClient.get<never, SupportedCurrencies>("/api/v1/account/supported-currencies"),
};
