"use client";

import { useRef, useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { ChevronDown, Menu, X } from "lucide-react";
import { cn } from "@/lib/utils";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { UserMenu } from "./user-menu";

const withdrawNavItems = [
  { label: "发起提现", href: "/withdraw" },
  { label: "数字货币地址", href: "/crypto-addresses" },
  { label: "银行账户", href: "/bank-accounts" },
] as const;

const navLinks = [
  { label: "概览", href: "/dashboard" },
  { label: "兑换", href: "/exchange" },
  { label: "交易记录", href: "/transactions" },
] as const;

const withdrawPaths = withdrawNavItems.map((item) => item.href);

function isWithdrawPath(pathname: string) {
  return withdrawPaths.some(
    (href) => pathname === href || pathname.startsWith(`${href}/`),
  );
}

function isLinkActive(pathname: string, href: string) {
  return pathname === href || (href !== "/dashboard" && pathname.startsWith(href));
}

export function TopNav() {
  const pathname = usePathname();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [withdrawOpen, setWithdrawOpen] = useState(false);
  const withdrawCloseTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const withdrawActive = isWithdrawPath(pathname);

  const openWithdrawMenu = () => {
    if (withdrawCloseTimer.current) {
      clearTimeout(withdrawCloseTimer.current);
      withdrawCloseTimer.current = null;
    }
    setWithdrawOpen(true);
  };

  const scheduleCloseWithdrawMenu = () => {
    withdrawCloseTimer.current = setTimeout(() => {
      setWithdrawOpen(false);
    }, 120);
  };

  return (
    <header className="fixed top-0 z-50 w-full border-b border-border bg-background/80 backdrop-blur-sm">
      <div className="mx-auto flex h-14 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
        <Link href="/dashboard" className="shrink-0">
          <span className="text-xl font-bold text-primary">Motewallet</span>
        </Link>

        <nav className="hidden items-center space-x-1 md:flex">
          <Link
            href="/dashboard"
            className={cn(
              "relative px-3 py-4 text-sm font-medium transition-colors",
              isLinkActive(pathname, "/dashboard")
                ? "text-primary"
                : "text-muted-foreground hover:text-foreground",
            )}
          >
            概览
            {isLinkActive(pathname, "/dashboard") && (
              <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary" />
            )}
          </Link>

          <DropdownMenu open={withdrawOpen} onOpenChange={setWithdrawOpen} modal={false}>
            <div onMouseEnter={openWithdrawMenu} onMouseLeave={scheduleCloseWithdrawMenu}>
              <DropdownMenuTrigger
                className={cn(
                  "relative inline-flex items-center gap-1 px-3 py-4 text-sm font-medium transition-colors outline-none",
                  withdrawActive
                    ? "text-primary"
                    : "text-muted-foreground hover:text-foreground",
                )}
              >
                提现
                <ChevronDown className="h-3.5 w-3.5 opacity-70" />
                {withdrawActive && (
                  <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary" />
                )}
              </DropdownMenuTrigger>
            </div>
            <DropdownMenuContent
              align="start"
              className="w-44"
              onMouseEnter={openWithdrawMenu}
              onMouseLeave={scheduleCloseWithdrawMenu}
            >
              {withdrawNavItems.map((item) => (
                <DropdownMenuItem key={item.href} asChild>
                  <Link
                    href={item.href}
                    className={cn(
                      isLinkActive(pathname, item.href) && "text-primary",
                    )}
                  >
                    {item.label}
                  </Link>
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>

          {navLinks.slice(1).map((item) => {
            const isActive = isLinkActive(pathname, item.href);
            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  "relative px-3 py-4 text-sm font-medium transition-colors",
                  isActive
                    ? "text-primary"
                    : "text-muted-foreground hover:text-foreground",
                )}
              >
                {item.label}
                {isActive && (
                  <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary" />
                )}
              </Link>
            );
          })}
        </nav>

        <div className="flex items-center gap-3">
          <UserMenu />
          <button
            className="p-1 text-muted-foreground hover:text-foreground md:hidden"
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
          >
            {mobileMenuOpen ? (
              <X className="h-5 w-5" />
            ) : (
              <Menu className="h-5 w-5" />
            )}
          </button>
        </div>
      </div>

      {mobileMenuOpen && (
        <nav className="border-t bg-background md:hidden">
          <div className="space-y-1 px-4 py-3">
            <Link
              href="/dashboard"
              onClick={() => setMobileMenuOpen(false)}
              className={cn(
                "block rounded-md px-3 py-2 text-sm font-medium",
                isLinkActive(pathname, "/dashboard")
                  ? "bg-accent text-primary"
                  : "text-muted-foreground hover:bg-muted hover:text-foreground",
              )}
            >
              概览
            </Link>

            <div className="pt-1">
              <p
                className={cn(
                  "px-3 py-2 text-xs font-semibold uppercase tracking-wide",
                  withdrawActive ? "text-primary" : "text-muted-foreground",
                )}
              >
                提现
              </p>
              <div className="space-y-1">
                {withdrawNavItems.map((item) => (
                  <Link
                    key={item.href}
                    href={item.href}
                    onClick={() => setMobileMenuOpen(false)}
                    className={cn(
                      "block rounded-md py-2 pl-6 pr-3 text-sm font-medium",
                      isLinkActive(pathname, item.href)
                        ? "bg-accent text-primary"
                        : "text-muted-foreground hover:bg-muted hover:text-foreground",
                    )}
                  >
                    {item.label}
                  </Link>
                ))}
              </div>
            </div>

            {navLinks.slice(1).map((item) => (
              <Link
                key={item.href}
                href={item.href}
                onClick={() => setMobileMenuOpen(false)}
                className={cn(
                  "block rounded-md px-3 py-2 text-sm font-medium",
                  isLinkActive(pathname, item.href)
                    ? "bg-accent text-primary"
                    : "text-muted-foreground hover:bg-muted hover:text-foreground",
                )}
              >
                {item.label}
              </Link>
            ))}
          </div>
        </nav>
      )}
    </header>
  );
}
