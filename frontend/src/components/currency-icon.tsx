import { clsx } from "clsx";

type CurrencyIconProps = {
  currency: string;
  className?: string;
};

type SvgProps = {
  className?: string;
};

function IconBase({
  className,
  children,
  bg,
}: SvgProps & { children: React.ReactNode; bg: string }) {
  return (
    <svg
      viewBox="0 0 32 32"
      aria-hidden="true"
      className={clsx("h-8 w-8", className)}
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <circle cx="16" cy="16" r="16" fill={bg} />
      {children}
    </svg>
  );
}

function GenericCoinIcon({ className }: SvgProps) {
  return (
    <IconBase className={className} bg="#CBD5E1">
      <circle cx="16" cy="16" r="10.5" fill="#94A3B8" />
      <path
        d="M11.5 16H20.5M16 11.5V20.5"
        stroke="white"
        strokeWidth="2"
        strokeLinecap="round"
      />
    </IconBase>
  );
}

function UsdtIcon({ className }: SvgProps) {
  return (
    <IconBase className={className} bg="#26A17B">
      <path
        d="M9 10.3H23V12.7H17.6V13.95C21.3 14.2 24 14.95 24 15.85C24 16.75 21.3 17.5 17.6 17.75V21.7H14.4V17.75C10.7 17.5 8 16.75 8 15.85C8 14.95 10.7 14.2 14.4 13.95V12.7H9V10.3Z"
        fill="white"
      />
      <ellipse cx="16" cy="15.85" rx="7.8" ry="1.45" fill="#DDF7EE" />
    </IconBase>
  );
}

function UsdcIcon({ className }: SvgProps) {
  return (
    <IconBase className={className} bg="#2775CA">
      <circle cx="16" cy="16" r="9.25" stroke="white" strokeWidth="2" />
      <path
        d="M18.8 12.7C18.15 12 17.2 11.65 16.15 11.7C14.65 11.8 13.45 12.7 13.45 13.95C13.45 15.2 14.5 15.75 16.2 16.1C17.5 16.35 18.05 16.75 18.05 17.55C18.05 18.35 17.25 19.05 16.05 19.05C14.95 19.05 14.05 18.65 13.4 18"
        stroke="white"
        strokeWidth="1.8"
        strokeLinecap="round"
      />
      <path
        d="M16 10.6V12M16 20V21.4M9.9 12.6C9.15 13.55 8.7 14.72 8.7 16C8.7 17.28 9.15 18.45 9.9 19.4M22.1 12.6C22.85 13.55 23.3 14.72 23.3 16C23.3 17.28 22.85 18.45 22.1 19.4"
        stroke="white"
        strokeWidth="1.8"
        strokeLinecap="round"
      />
    </IconBase>
  );
}

function BtcIcon({ className }: SvgProps) {
  return (
    <IconBase className={className} bg="#F7931A">
      <path
        d="M13.15 9.4H16.1C18.85 9.4 20.7 10.75 20.7 12.95C20.7 14.15 20.05 15.1 18.95 15.55C20.4 15.95 21.25 17 21.25 18.45C21.25 21.05 19.15 22.6 15.9 22.6H13.15V9.4Z"
        fill="white"
      />
      <path
        d="M15.4 11.7V14.7H16.1C17.55 14.7 18.45 14.2 18.45 13.2C18.45 12.2 17.75 11.7 16.35 11.7H15.4ZM15.4 17V20.2H16.45C18.15 20.2 19.05 19.65 19.05 18.5C19.05 17.35 18.05 17 16.45 17H15.4Z"
        fill="#F7931A"
      />
      <path
        d="M14.4 8.5V23.5M17.2 8.5V23.5"
        stroke="white"
        strokeWidth="1.35"
        strokeLinecap="round"
      />
    </IconBase>
  );
}

function UsdIcon({ className }: SvgProps) {
  return (
    <IconBase className={className} bg="#16A34A">
      <circle cx="16" cy="16" r="9.5" fill="#15803D" />
      <path
        d="M18.25 12.65C17.65 12.05 16.85 11.75 15.95 11.8C14.65 11.85 13.65 12.6 13.65 13.75C13.65 14.95 14.6 15.45 16.1 15.75C17.3 16 17.85 16.35 17.85 17.1C17.85 17.85 17.05 18.45 15.95 18.45C14.95 18.45 14.15 18.1 13.55 17.45"
        stroke="white"
        strokeWidth="1.8"
        strokeLinecap="round"
      />
      <path
        d="M16 10.9V12.05M16 18.4V19.55"
        stroke="white"
        strokeWidth="1.8"
        strokeLinecap="round"
      />
    </IconBase>
  );
}

function HkdIcon({ className }: SvgProps) {
  return (
    <IconBase className={className} bg="#E11D48">
      <path d="M16 10.25C17.3 11.15 18.1 12.4 18.35 14C16.8 13.95 15.55 13.4 14.6 12.35C14.8 11.5 15.25 10.8 16 10.25Z" fill="white" />
      <path d="M20.95 13.35C21 14.85 20.55 16.15 19.6 17.2C18.75 15.95 18.45 14.65 18.7 13.35C19.45 12.95 20.2 12.95 20.95 13.35Z" fill="white" />
      <path d="M19.1 18.65C18.2 19.85 17 20.6 15.5 20.9C15.55 19.4 16.05 18.2 17 17.25C17.8 17.35 18.5 17.8 19.1 18.65Z" fill="white" />
      <path d="M12.9 20.45C11.55 19.95 10.55 19.1 9.9 17.9C11.3 17.45 12.55 17.55 13.75 18.2C13.9 19 13.55 19.75 12.9 20.45Z" fill="white" />
      <path d="M10.95 14.15C11.6 12.8 12.65 11.85 14.05 11.3C14.35 12.7 14.1 13.95 13.35 15.05C12.55 15.15 11.75 14.85 10.95 14.15Z" fill="white" />
      <circle cx="16" cy="16" r="1.3" fill="#EF4444" />
    </IconBase>
  );
}

function EurIcon({ className }: SvgProps) {
  return (
    <IconBase className={className} bg="#1D4ED8">
      <path
        d="M20.8 11.75C19.85 10.75 18.45 10.15 16.9 10.15C14.35 10.15 12.15 11.65 11.15 13.95H9.7V15.4H10.65C10.6 15.65 10.6 15.9 10.6 16.15C10.6 16.45 10.6 16.7 10.65 16.95H9.7V18.4H11.15C12.15 20.7 14.35 22.2 16.9 22.2C18.45 22.2 19.85 21.6 20.8 20.6"
        stroke="#FCD34D"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M10 15.4H16.8M10 18.4H16.8"
        stroke="#FCD34D"
        strokeWidth="1.9"
        strokeLinecap="round"
      />
      <circle cx="16" cy="7.9" r="0.75" fill="#FCD34D" />
      <circle cx="12.35" cy="8.85" r="0.75" fill="#FCD34D" />
      <circle cx="19.65" cy="8.85" r="0.75" fill="#FCD34D" />
      <circle cx="22.15" cy="11.35" r="0.75" fill="#FCD34D" />
      <circle cx="23.1" cy="15" r="0.75" fill="#FCD34D" />
      <circle cx="22.15" cy="18.65" r="0.75" fill="#FCD34D" />
      <circle cx="19.65" cy="21.15" r="0.75" fill="#FCD34D" />
      <circle cx="16" cy="22.1" r="0.75" fill="#FCD34D" />
      <circle cx="12.35" cy="21.15" r="0.75" fill="#FCD34D" />
    </IconBase>
  );
}

const ICON_MAP: Record<string, (props: SvgProps) => JSX.Element> = {
  USD: UsdIcon,
  USDT: UsdtIcon,
  USDC: UsdcIcon,
  BTC: BtcIcon,
  HKD: HkdIcon,
  EUR: EurIcon,
};

export function CurrencyIcon({ currency, className }: CurrencyIconProps) {
  const Icon = ICON_MAP[currency] ?? GenericCoinIcon;

  return <Icon className={className} />;
}
