import { getKycFieldMeta, type KycFieldMetaKey } from "@/lib/kyc/field-meta";

export interface ParsedKycFieldError {
  path: string;
  message: string;
  label: string;
}

const MESSAGE_MAP: Record<string, string> = {
  "Country not supported": "暂不支持该国家/地区",
};

function translateMessage(message: string): string {
  if (MESSAGE_MAP[message]) return MESSAGE_MAP[message];
  if (message.startsWith("is not supported for ")) {
    return `不支持：${message.slice("is not supported for ".length)}`;
  }
  const requiredMatch = message.match(/^required property '([^']+)' not found$/);
  if (requiredMatch) return "此项为必填";
  return message;
}

function jsonPathToFormPath(path: string): string {
  return path.replace(/\[(\d+)\]/g, ".$1");
}

const PERSON_FIELD_KEYS = new Set([
  "gender",
  "idCard",
  "nameEN",
  "country",
  "nameCHS",
  "surname",
  "authType",
  "birthday",
  "surnameCHS",
  "certificateStart",
  "certificateEnd",
  "residenceAddress",
  "residenceCountry",
  "shareholdingRatio",
  "idHoldingPaths",
]);

function resolveFieldLabel(path: string): string {
  const parts = path.split(".");
  const field = parts[parts.length - 1];
  const section = parts[0];
  const indexPart = parts.length > 2 ? parts[1] : undefined;
  const index =
    indexPart !== undefined && /^\d+$/.test(indexPart)
      ? parseInt(indexPart, 10)
      : undefined;

  if (section === "shareholdersInfo" && PERSON_FIELD_KEYS.has(field)) {
    const meta = getKycFieldMeta(`person.${field}` as KycFieldMetaKey);
    return index !== undefined ? `股东 ${index + 1} · ${meta.label}` : meta.label;
  }
  if (section === "directorInfo" && PERSON_FIELD_KEYS.has(field)) {
    const meta = getKycFieldMeta(`person.${field}` as KycFieldMetaKey);
    return index !== undefined ? `董事 ${index + 1} · ${meta.label}` : meta.label;
  }
  if (section === "enterpriseInfo") {
    const meta = getKycFieldMeta(field as KycFieldMetaKey);
    return meta?.label ?? field;
  }
  return field;
}

export function parseKycFieldError(error: string): ParsedKycFieldError | null {
  if (!error.startsWith("$.")) return null;

  const rest = error.slice(2);

  const colonMatch = rest.match(/^([^:]+):\s*(.+)$/);
  if (colonMatch) {
    const objectPath = colonMatch[1];
    const msg = colonMatch[2];
    const propMatch = msg.match(/required property '([^']+)' not found/);
    if (propMatch) {
      const path = jsonPathToFormPath(`${objectPath}.${propMatch[1]}`);
      return {
        path,
        message: translateMessage(msg),
        label: resolveFieldLabel(path),
      };
    }
    const path = jsonPathToFormPath(objectPath);
    return {
      path,
      message: translateMessage(msg),
      label: resolveFieldLabel(path),
    };
  }

  const match = rest.match(/^([\w.[\]0-9]+)\s+(.+)$/);
  if (match) {
    const path = jsonPathToFormPath(match[1]);
    return {
      path,
      message: translateMessage(match[2]),
      label: resolveFieldLabel(path),
    };
  }

  return null;
}

export function parseKycFieldErrors(errors: string[]): ParsedKycFieldError[] {
  return errors
    .map(parseKycFieldError)
    .filter((item): item is ParsedKycFieldError => item !== null);
}
