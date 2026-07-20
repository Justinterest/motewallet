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

/** Official Tether token mark (cryptocurrency-icons / Tether brand green). */
function UsdtIcon({ className }: SvgProps) {
  return (
    <IconBase className={className} bg="#26A17B">
      <path
        fill="#FFF"
        d="M17.922 17.383v-.002c-.11.008-.677.042-1.942.042-1.01 0-1.721-.03-1.971-.042v.003c-3.888-.171-6.79-.848-6.79-1.658 0-.809 2.902-1.486 6.79-1.66v2.644c.254.018.982.061 1.988.061 1.207 0 1.812-.05 1.925-.06v-2.643c3.88.173 6.775.85 6.775 1.658 0 .81-2.895 1.485-6.775 1.657m0-3.59v-2.366h5.414V7.819H8.595v3.608h5.414v2.365c-4.4.202-7.709 1.074-7.709 2.118 0 1.044 3.309 1.915 7.709 2.118v7.582h3.913v-7.584c4.393-.202 7.694-1.073 7.694-2.116 0-1.043-3.301-1.914-7.694-2.117"
      />
    </IconBase>
  );
}

/** Official USDC token mark (Circle brand blue). */
function UsdcIcon({ className }: SvgProps) {
  return (
    <IconBase className={className} bg="#2775CA">
      <path
        fill="#FFF"
        d="M20.022 18.124c0-2.124-1.28-2.852-3.84-3.156-1.828-.243-2.193-.728-2.193-1.578 0-.85.61-1.396 1.828-1.396 1.097 0 1.707.364 2.011 1.275a.458.458 0 00.427.303h.975a.416.416 0 00.427-.425v-.06a3.04 3.04 0 00-2.743-2.489V9.142c0-.243-.183-.425-.487-.486h-.915c-.243 0-.426.182-.487.486v1.396c-1.829.242-2.986 1.456-2.986 2.974 0 2.002 1.218 2.791 3.778 3.095 1.707.303 2.255.668 2.255 1.639 0 .97-.853 1.638-2.011 1.638-1.585 0-2.133-.667-2.316-1.578-.06-.242-.244-.364-.427-.364h-1.036a.416.416 0 00-.426.425v.06c.243 1.518 1.219 2.61 3.23 2.914v1.457c0 .242.183.425.487.485h.915c.243 0 .426-.182.487-.485V21.34c1.829-.303 3.047-1.578 3.047-3.217z"
      />
      <path
        fill="#FFF"
        d="M12.892 24.497c-4.754-1.7-7.192-6.98-5.424-11.653.914-2.55 2.925-4.491 5.424-5.402.244-.121.365-.303.365-.607v-.85c0-.242-.121-.424-.365-.485-.061 0-.183 0-.244.06a10.895 10.895 0 00-7.13 13.717c1.096 3.4 3.717 6.01 7.13 7.102.244.121.488 0 .548-.243.061-.06.061-.122.061-.243v-.85c0-.182-.182-.424-.365-.546zm6.46-18.936c-.244-.122-.488 0-.548.242-.061.061-.061.122-.061.243v.85c0 .243.182.485.365.607 4.754 1.7 7.192 6.98 5.424 11.653-.914 2.55-2.925 4.491-5.424 5.402-.244.121-.365.303-.365.607v.85c0 .242.121.424.365.485.061 0 .183 0 .244-.06a10.895 10.895 0 007.13-13.717c-1.096-3.46-3.778-6.07-7.13-7.162z"
      />
    </IconBase>
  );
}

/** Official Bitcoin logo (bitcoin.org orange). */
function BtcIcon({ className }: SvgProps) {
  return (
    <IconBase className={className} bg="#F7931A">
      <path
        fill="#FFF"
        fillRule="nonzero"
        d="M23.189 14.02c.314-2.096-1.283-3.223-3.465-3.975l.708-2.84-1.728-.43-.69 2.765c-.454-.114-.92-.22-1.385-.326l.695-2.783L15.596 6l-.708 2.839c-.376-.086-.746-.17-1.104-.26l.002-.009-2.384-.595-.46 1.846s1.283.294 1.256.312c.7.175.826.638.805 1.006l-.806 3.235c.048.012.11.03.18.057l-.183-.045-1.13 4.532c-.086.212-.303.531-.793.41.018.025-1.256-.313-1.256-.313l-.858 1.978 2.25.561c.418.105.828.215 1.231.318l-.715 2.872 1.727.43.708-2.84c.472.127.93.245 1.378.357l-.706 2.828 1.728.43.715-2.866c2.948.558 5.164.333 6.097-2.333.752-2.146-.037-3.385-1.588-4.192 1.13-.26 1.98-1.003 2.207-2.538zm-3.95 5.538c-.533 2.147-4.148.986-5.32.695l.95-3.805c1.172.293 4.929.872 4.37 3.11zm.535-5.569c-.487 1.953-3.495.96-4.47.717l.86-3.45c.975.243 4.118.696 3.61 2.733z"
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

const ICON_MAP: Record<string, (props: SvgProps) => React.ReactElement> = {
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
