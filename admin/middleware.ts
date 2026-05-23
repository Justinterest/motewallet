import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const adminToken = request.cookies.get("admin_token");

  // 已登录用户访问登录页面，重定向到仪表盘
  if (pathname.startsWith("/login") && adminToken) {
    return NextResponse.redirect(new URL("/dashboard", request.url));
  }

  // 未登录用户访问仪表盘，重定向到登录页面
  if (pathname.startsWith("/dashboard") && !adminToken) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/dashboard/:path*", "/login"],
};
