"use client";

import { useState } from "react";
import Link from "next/link";
import { Search } from "lucide-react";
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
import { useAdminDeposits } from "@/lib/hooks/use-deposits";
import { formatAmount } from "@/lib/utils/format";

const CURRENCY_OPTIONS = [
  { value: "ALL", label: "全部币种" },
  { value: "USDT", label: "USDT" },
  { value: "USDC", label: "USDC" },
  { value: "BTC", label: "BTC" },
];

const STATUS_OPTIONS = [
  { value: "ALL", label: "全部状态" },
  { value: "COMPLETED", label: "已到账" },
  { value: "PROCESSING", label: "处理中" },
  { value: "PENDING", label: "待处理" },
  { value: "FAILED", label: "失败" },
];

function getStatusBadge(status: string) {
  switch (status) {
    case "COMPLETED":
      return <Badge className="bg-green-100 text-green-700 border-green-200">已到账</Badge>;
    case "PROCESSING":
    case "PENDING":
      return <Badge className="bg-amber-100 text-amber-700 border-amber-200">处理中</Badge>;
    case "FAILED":
      return <Badge className="bg-red-100 text-red-700 border-red-200">失败</Badge>;
    default:
      return <Badge variant="outline">{status}</Badge>;
  }
}

export default function DepositsPage() {
  const [page, setPage] = useState(1);
  const [merchantEmail, setMerchantEmail] = useState("");
  const [searchEmail, setSearchEmail] = useState("");
  const [currency, setCurrency] = useState("ALL");
  const [status, setStatus] = useState("ALL");

  const { data, isLoading } = useAdminDeposits({
    page,
    merchantEmail: searchEmail || undefined,
    currency,
    status,
  });

  const deposits = data?.deposits || [];
  const total = data?.total ?? 0;
  const pageSize = 20;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  function handleSearch() {
    setSearchEmail(merchantEmail.trim());
    setPage(1);
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-slate-900">数币充值记录</h1>
        <p className="mt-1 text-sm text-slate-500">查看商户数字货币充值到账记录</p>
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
            <div className="w-[140px]">
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
            <div className="w-[140px]">
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
            充值列表
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
          ) : deposits.length === 0 ? (
            <p className="py-12 text-center text-sm text-slate-400">暂无充值记录</p>
          ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>平台订单号</TableHead>
                    <TableHead>商户</TableHead>
                    <TableHead>金额</TableHead>
                    <TableHead>网络</TableHead>
                    <TableHead>交易哈希</TableHead>
                    <TableHead>状态</TableHead>
                    <TableHead>充值时间</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {deposits.map((item) => (
                    <TableRow key={item.id}>
                      <TableCell className="font-mono text-xs">{item.platform_order_id}</TableCell>
                      <TableCell>
                        <Link
                          href={`/merchants/${item.merchant_id}`}
                          className="text-sm text-[#1E40AF] hover:underline"
                        >
                          {item.merchant_email}
                        </Link>
                      </TableCell>
                      <TableCell className="font-medium">
                        {formatAmount(item.amount, item.currency)} {item.currency}
                      </TableCell>
                      <TableCell className="text-sm">{item.network}</TableCell>
                      <TableCell className="max-w-[160px] truncate font-mono text-xs" title={item.tx_hash}>
                        {item.tx_hash || "—"}
                      </TableCell>
                      <TableCell>{getStatusBadge(item.status)}</TableCell>
                      <TableCell className="text-xs text-slate-500">
                        {new Date(item.created_at).toLocaleString("zh-CN")}
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
