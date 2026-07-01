"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  LayoutDashboard,
  Building2,
  Receipt,
  Activity,
  ClipboardCheck,
  ArrowDownToLine,
  ArrowUpFromLine,
  ArrowRightLeft,
  Settings,
  Users,
  FileText,
  Webhook,
} from "lucide-react";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { useAuthStore } from "@/stores/auth-store";

const mainNavItems = [
  {
    title: "仪表盘",
    href: "/dashboard",
    icon: LayoutDashboard,
  },
  {
    title: "商户管理",
    href: "/merchants",
    icon: Building2,
  },
  {
    title: "手续费模板",
    href: "/fee-templates",
    icon: Receipt,
  },
  {
    title: "交易监控",
    href: "#",
    icon: Activity,
  },
];

const operationNavItems = [
  {
    title: "数币充值记录",
    href: "/deposits",
    icon: ArrowDownToLine,
  },
  {
    title: "提现记录",
    href: "/withdrawal-records",
    icon: ArrowUpFromLine,
  },
  {
    title: "兑换记录",
    href: "/exchanges",
    icon: ArrowRightLeft,
  },
  {
    title: "提现审核",
    href: "/withdrawals",
    icon: ClipboardCheck,
  },
  {
    title: "系统设置",
    href: "#",
    icon: Settings,
  },
];

const systemNavItems = [
  {
    title: "管理员",
    href: "#",
    icon: Users,
  },
  {
    title: "审计日志",
    href: "#",
    icon: FileText,
  },
  {
    title: "Webhook 日志",
    href: "#",
    icon: Webhook,
  },
];

export function AppSidebar() {
  const pathname = usePathname();
  const { admin } = useAuthStore();

  return (
    <Sidebar className="border-r border-slate-200 bg-white">
      <SidebarHeader className="border-b border-slate-200 px-6 py-4">
        <Link href="/dashboard" className="flex items-center gap-2">
          <span className="text-lg font-bold text-[#1E40AF]">Motewallet</span>
        </Link>
      </SidebarHeader>

      <SidebarContent className="px-2 py-2">
        <SidebarGroup>
          <SidebarGroupLabel>主菜单</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {mainNavItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton
                    asChild
                    isActive={pathname === item.href || (item.href !== "#" && item.href !== "/dashboard" && pathname.startsWith(item.href))}
                    tooltip={item.title}
                  >
                    <Link href={item.href}>
                      <item.icon className="size-4" />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarGroup>
          <SidebarGroupLabel>运营</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {operationNavItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton
                    asChild
                    isActive={pathname === item.href || (item.href !== "#" && pathname.startsWith(item.href))}
                    tooltip={item.title}
                  >
                    <Link href={item.href}>
                      <item.icon className="size-4" />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarGroup>
          <SidebarGroupLabel>系统</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {systemNavItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton
                    asChild
                    isActive={pathname === item.href}
                    tooltip={item.title}
                  >
                    <Link href={item.href}>
                      <item.icon className="size-4" />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter className="border-t border-slate-200 px-4 py-3">
        <div className="flex items-center gap-3">
          <div className="flex size-8 items-center justify-center rounded-full bg-[#1E40AF] text-xs font-medium text-white">
            {admin?.username?.charAt(0)?.toUpperCase() || "A"}
          </div>
          <div className="flex flex-col">
            <span className="text-sm font-medium text-slate-900">
              {admin?.username || "管理员"}
            </span>
            <span className="text-xs text-slate-500">
              {admin?.role || "admin"}
            </span>
          </div>
        </div>
      </SidebarFooter>
    </Sidebar>
  );
}
