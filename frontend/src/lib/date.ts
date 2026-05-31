import { format, isValid, parse } from "date-fns"

const STORAGE_FORMAT = "yyyy-MM-dd"
const DISPLAY_FORMAT = "yyyy年M月d日"

export function parseDateString(value?: string): Date | undefined {
  if (!value) return undefined

  const parsed = parse(value, STORAGE_FORMAT, new Date())
  return isValid(parsed) ? parsed : undefined
}

export function formatDateString(date: Date): string {
  return format(date, STORAGE_FORMAT)
}

export function formatDateDisplay(value?: string): string {
  const date = parseDateString(value)
  if (!date) return ""
  return format(date, DISPLAY_FORMAT)
}
