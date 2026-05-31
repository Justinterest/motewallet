"use client"

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  type SelectOption,
} from "@/components/ui/select"

interface FormSelectProps {
  value?: string
  onChange: (value: string) => void
  onBlur?: () => void
  options: SelectOption[] | readonly SelectOption[]
  placeholder?: string
  disabled?: boolean
  className?: string
}

function FormSelect({
  value,
  onChange,
  onBlur,
  options,
  placeholder = "请选择",
  disabled,
  className,
}: FormSelectProps) {
  return (
    <Select
      value={value || undefined}
      onValueChange={onChange}
      disabled={disabled}
    >
      <SelectTrigger className={className} onBlur={onBlur}>
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        {options.map((option) => (
          <SelectItem key={option.value} value={option.value}>
            {option.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}

export { FormSelect }
