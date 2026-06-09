import { KYC_DRAFT_KEY } from "@/lib/kyc/constants";
import { createKycFormDefaults } from "@/lib/kyc/form-defaults";
import type { KycFormValues } from "@/lib/validations/onboarding";

export interface KycDraft {
  step: number;
  values: KycFormValues;
}

export function mergeKycFormValues(
  defaults: KycFormValues,
  draft: Partial<KycFormValues>
): KycFormValues {
  const shareholderDefaults = defaults.shareholdersInfo[0];
  const directorDefaults = defaults.directorInfo[0];

  return {
    enterpriseInfo: {
      ...defaults.enterpriseInfo,
      ...(draft.enterpriseInfo ?? {}),
    },
    shareholdersInfo:
      draft.shareholdersInfo && draft.shareholdersInfo.length > 0
        ? draft.shareholdersInfo.map((person) => ({
            ...shareholderDefaults,
            ...person,
          }))
        : defaults.shareholdersInfo,
    directorInfo:
      draft.directorInfo && draft.directorInfo.length > 0
        ? draft.directorInfo.map((person) => ({
            ...directorDefaults,
            ...person,
          }))
        : defaults.directorInfo,
  };
}

export function loadKycDraft(): KycDraft | null {
  try {
    const raw = localStorage.getItem(KYC_DRAFT_KEY);
    if (!raw) return null;

    const parsed = JSON.parse(raw) as
      | KycDraft
      | (Partial<KycFormValues> & { step?: number });

    const defaults = createKycFormDefaults();

    if (parsed && typeof parsed === "object" && "values" in parsed && parsed.values) {
      return {
        step:
          typeof parsed.step === "number"
            ? Math.max(0, Math.min(parsed.step, 3))
            : 0,
        values: mergeKycFormValues(defaults, parsed.values),
      };
    }

    if (parsed && typeof parsed === "object" && "enterpriseInfo" in parsed) {
      return {
        step:
          typeof parsed.step === "number"
            ? Math.max(0, Math.min(parsed.step, 3))
            : 0,
        values: mergeKycFormValues(defaults, parsed),
      };
    }

    return null;
  } catch {
    return null;
  }
}

export function saveKycDraft(step: number, values: KycFormValues): void {
  try {
    const draft: KycDraft = { step, values };
    localStorage.setItem(KYC_DRAFT_KEY, JSON.stringify(draft));
  } catch {
    /* quota exceeded */
  }
}

export function clearKycDraft(): void {
  localStorage.removeItem(KYC_DRAFT_KEY);
}
