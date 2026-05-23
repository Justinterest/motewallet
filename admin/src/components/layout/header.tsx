"use client";

import { LogOut } from "lucide-react";
import { SidebarTrigger } from "@/components/ui/sidebar";
import { Button } from "@/components/ui/button";
import { useCurrentAdmin, useAdminLogout } from "@/lib/hooks/use-auth";

export function Header() {
  const { data: admin } = useCurrentAdmin();
  const logoutMutation = useAdminLogout();

  const handleLogout = () => {
    logoutMutation.mutate();
  };

  return (
    <header className="flex h-14 items-center justify-between border-b bg-white px-6">
      <div className="flex items-center gap-2">
        <SidebarTrigger />
      </div>
      <div className="flex items-center gap-4">
        {admin && (
          <>
            <span className="text-sm text-slate-700">{admin.username}</span>
            <span className="rounded-full bg-[#1E40AF]/10 px-2 py-0.5 text-xs font-medium text-[#1E40AF]">
              {admin.role}
            </span>
          </>
        )}
        <Button
          variant="ghost"
          size="sm"
          onClick={handleLogout}
          disabled={logoutMutation.isPending}
          className="text-slate-500 hover:text-slate-700"
        >
          <LogOut className="mr-1 size-4" />
          退出登录
        </Button>
      </div>
    </header>
  );
}
