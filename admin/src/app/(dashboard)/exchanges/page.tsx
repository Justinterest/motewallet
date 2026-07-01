"use client";

import { useState } from "react";
import Link from "next/link";
import { Loader2, RefreshCw, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { useAdminExchanges, useSyncExchangeStatus } from "@/lib/hooks/use-exchanges";
import { formatAmount } from "@/lib/utils/format";
import { useToast } from "@/hooks/use-toast";

const CURRENCY_OPTIONS = [
  { value: "ALL", label: "全部币种" },
  { value: "USDT", label: "USDT" },
  { value: "USDC", label: "USDC" },
  { value: "USD", label: "USD" },
];

const STATUS_OPTIONS = [
  { value: "ALL", label: "全部状态" },
  { value: "COMPLETED", label: "已完成" },
  { value: "PROCESSING", label: "处理中" },
  { value: "FAILED", label: "失败" },
];

function getStatusBadge(status: string) {
  switch (status) {
    case "COMPLETED":
      return <Badge className="border-green-200 bg-green-100 text-green-700">已完成</Badge>;
    case "PROCESSING":
    case "PENDING":
      return <Badge className="border-amber-200 bg-amber-100 text-amber-700">处理中</Badge>;
    case "FAILED":
      return <Badge className="border-red-200 bg-red-100 text-red-700">失败</Badge>;
    default:
      return <Badge variant="outline">{status}</Badge>;
  }
}

function canSync(status: string) {
  return status === "PROCESSING" || status === "PENDING" || status === "FAILED";
}

export default function ExchangesPage() {
  const [page, setPage] = useState(1);
  const [merchantEmail, setMerchantEmail] = useState("");
  const [searchEmail, setSearchEmail] = useState("");
  const [currency, setCurrency] = useState("ALL");
  const [status, setStatus] = useState("ALL");
  const [syncingId, setSyncingId] = useState<number | null>(null);

  const { data, isLoading } = useAdminExchanges({
    page,
    merchantEmail: searchEmail || undefined,
    currency,
    status,
  });
  const syncMutation = useSyncExchangeStatus();
  const { toast } = useToast();

  const exchanges = data?.exchanges || [];
  const total = data?.total ?? 0;
  const pageSize = 20;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  function handleSearch() {
    setSearchEmail(merchantEmail.trim());
    setPage(1);
  }

  function handleSync(id: number) {
    setSyncingId(id);
    syncMutation.mutate(id, {
      onSuccess: (result) => {
        toast({
          title: result.updated ? "同步成功" : "状态未变化",
          description: result.updated
            ? result.status === "FAILED" && result.fail_reason
              ? `订单 #${id} 失败原因已更新：${result.fail_reason}`
              : `订单 #${id} 已更新为 ${result.status}${result.status === "FAILED" && result.fail_reason ? `：${result.fail_reason}` : ""}`
            : `KUN 状态：${result.kun_status || "—"}，当前：${result.status}`,
        });
      },
      onError: (err) => {
        toast({ variant: "destructive", title: "同步失败", description: err.message });
      },
      onSettled: () => setSyncingId(null),
    });
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-slate-900">兑换记录</h1>
        <p className="mt-1 text-sm text-slate-500">查看商户兑换订单，支持手动同步 KUN 状态</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">筛选</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap items-end gap-3">
            <div className="min-w-[220px] flex-1">
              <label className="mb-1.5 block text-xs font-medium text-slate-600">商户邮箱</label>
              <Input
                placeholder="搜索商户邮箱"
                value={merchantEmail}
                onChange={(e) => setMerchantEmail(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && handleSearch()}
              />
            </div>
            <div className="w-[120px]">
              <label className="mb-1.5 block text-xs font-medium text-slate-600">币种</label>
              <Select value={currency} onValueChange={(v) => { setCurrency(v); setPage(1); }}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {CURRENCY_OPTIONS.map((opt) => (
                    <SelectItem key={opt.value} value={opt.value}>
                      {opt.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="w-[120px]">
              <label className="mb-1.5 block text-xs font-medium text-slate-600">状态</label>
              <Select value={status} onValueChange={(v) => { setStatus(v); setPage(1); }}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {STATUS_OPTIONS.map((opt) => (
                    <SelectItem key={opt.value} value={opt.value}>
                      {opt.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <Button onClick={handleSearch}>
              <Search className="mr-1 h-4 w-4" />
              查询
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">
            兑换列表
            {total > 0 && (
              <span className="ml-2 text-sm font-normal text-slate-500">共 {total} 条</span>
            )}
          </CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-3">
              <Skeleton className="h-12 w-full" />
              <Skeleton className="h-12 w-full" />
              <Skeleton className="h-12 w-full" />
            </div>
          ) : exchanges.length === 0 ? (
            <p className="py-12 text-center text-sm text-slate-400">暂无兑换记录</p>
          ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>订单号</TableHead>
                    <TableHead>商户</TableHead>
                    <TableHead>兑换</TableHead>
                    <TableHead>卖出</TableHead>
                    <TableHead>买入</TableHead>
                    <TableHead>手续费</TableHead>
                    <TableHead>状态</TableHead>
                    <TableHead>失败原因</TableHead>
                    <TableHead>时间</TableHead>
                    <TableHead className="text-right">操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {exchanges.map((item) => (
                    <TableRow key={item.id}>
                      <TableCell className="font-mono text-xs">#{item.id}</TableCell>
                      <TableCell>
                        <Link
                          href={`/merchants/${item.merchant_id}`}
                          className="text-sm text-[#1E40AF] hover:underline"
                        >
                          {item.merchant_email}
                        </Link>
                      </TableCell>
                      <TableCell>
                        <Badge variant="secondary">{item.exchange_type}</Badge>
                      </TableCell>
                      <TableCell className="font-medium">
                        {formatAmount(item.from_amount, item.from_currency)} {item.from_currency}
                      </TableCell>
                      <TableCell>
                        {item.to_amount
                          ? `${formatAmount(item.to_amount, item.to_currency)} ${item.to_currency}`
                          : "—"}
                      </TableCell>
                      <TableCell className="text-slate-500">
                        {formatAmount(item.platform_fee, item.from_currency)}
                      </TableCell>
                      <TableCell>{getStatusBadge(item.status)}</TableCell>
                      <TableCell className="max-w-[200px] text-xs text-red-600">
                        {item.status === "FAILED" && item.fail_reason ? item.fail_reason : "—"}
                      </TableCell>
                      <TableCell className="text-xs text-slate-500">
                        {new Date(item.created_at).toLocaleString("zh-CN")}
                      </TableCell>
                      <TableCell className="text-right">
                        {canSync(item.status) ? (
                          <Button
                            variant="outline"
                            size="sm"
                            disabled={syncingId === item.id}
                            onClick={() => handleSync(item.id)}
                          >
                            {syncingId === item.id ? (
                              <Loader2 className="mr-1 h-3.5 w-3.5 animate-spin" />
                            ) : (
                              <RefreshCw className="mr-1 h-3.5 w-3.5" />
                            )}
                            同步状态
                          </Button>
                        ) : (
                          <span className="text-xs text-slate-400">—</span>
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>

              {totalPages > 1 && (
                <div className="mt-4 flex items-center justify-between">
                  <p className="text-sm text-slate-500">
                    第 {page} / {totalPages} 页
                  </p>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      disabled={page <= 1}
                      onClick={() => setPage((p) => Math.max(1, p - 1))}
                    >
                      上一页
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      disabled={page >= totalPages}
                      onClick={() => setPage((p) => p + 1)}
                    >
                      下一页
                    </Button>
                  </div>
                </div>
              )}
            </>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
