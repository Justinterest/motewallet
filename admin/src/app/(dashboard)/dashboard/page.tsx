"use client";

import {
  Building2,
  ArrowLeftRight,
  Clock,
  DollarSign,
} from "lucide-react";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

const statsCards = [
  {
    title: "商户总数",
    icon: Building2,
  },
  {
    title: "今日交易",
    icon: ArrowLeftRight,
  },
  {
    title: "待审核提现",
    icon: Clock,
  },
  {
    title: "今日手续费收入",
    icon: DollarSign,
  },
];

export default function DashboardPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-slate-900">仪表盘</h1>

      {/* 统计卡片 */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {statsCards.map((card) => (
          <Card key={card.title}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-slate-500">
                {card.title}
              </CardTitle>
              <card.icon className="size-4 text-slate-400" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-8 w-24" />
              <Skeleton className="mt-2 h-3 w-32" />
            </CardContent>
          </Card>
        ))}
      </div>

      {/* 最近交易 */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">最近交易</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            <div className="flex items-center gap-4">
              <Skeleton className="h-4 w-full" />
            </div>
            <div className="flex items-center gap-4">
              <Skeleton className="h-4 w-full" />
            </div>
            <div className="flex items-center gap-4">
              <Skeleton className="h-4 w-full" />
            </div>
            <div className="flex items-center gap-4">
              <Skeleton className="h-4 w-full" />
            </div>
            <div className="flex items-center gap-4">
              <Skeleton className="h-4 w-full" />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* 待处理事项 */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">待处理事项</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col items-center justify-center py-8 text-slate-400">
            <Clock className="mb-2 size-8" />
            <p className="text-sm">暂无待处理事项</p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
