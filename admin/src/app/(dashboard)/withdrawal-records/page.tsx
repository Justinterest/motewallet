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
import { useAdminWithdrawals } from "@/lib/hooks/use-withdrawals";
import { formatAmount } from "@/lib/utils/format";

const CURRENCY_OPTIONS = [
  { value: "ALL", label: "全部币种" },
  { value: "USDT", label: "USDT" },
  { value: "USDC", label: "USDC" },
  { value: "BTC", label: "BTC" },
  { value: "USD", label: "USD" },
];

const STATUS_OPTIONS = [
  { value: "ALL", label: "全部状态" },
  { value: "COMPLETED", label: "已完成" },
  { value: "PROCESSING", label: "处理中" },
  { value: "PENDING", label: "待处理" },
  { value: "FAILED", label: "失败" },
];

const REVIEW_STATUS_OPTIONS = [
  { value: "ALL", label: "全部审核" },
  { value: "PENDING_REVIEW", label: "待审核" },
  { value: "APPROVED", label: "已通过" },
  { value: "REJECTED", label: "已拒绝" },
];

const TYPE_OPTIONS = [
  { value: "ALL", label: "全部类型" },
  { value: "CRYPTO", label: "数币" },
  { value: "FIAT", label: "法币" },
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

function getReviewStatusBadge(status: string) {
  switch (status) {
    case "PENDING_REVIEW":
      return <Badge className="border-amber-200 bg-amber-100 text-amber-700">待审核</Badge>;
    case "APPROVED":
      return <Badge className="border-green-200 bg-green-100 text-green-700">已通过</Badge>;
    case "REJECTED":
      return <Badge className="border-red-200 bg-red-100 text-red-700">已拒绝</Badge>;
    default:
      return <Badge variant="outline">{status}</Badge>;
  }
}

export default function WithdrawalRecordsPage() {
  const [page, setPage] = useState(1);
  const [merchantEmail, setMerchantEmail] = useState("");
  const [searchEmail, setSearchEmail] = useState("");
  const [currency, setCurrency] = useState("ALL");
  const [status, setStatus] = useState("ALL");
  const [reviewStatus, setReviewStatus] = useState("ALL");
  const [type, setType] = useState("ALL");

  const { data, isLoading } = useAdminWithdrawals({
    page,
    merchantEmail: searchEmail || undefined,
    currency,
    status,
    reviewStatus,
    type,
  });

  const withdrawals = data?.withdrawals || [];
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
        <h1 className="text-2xl font-bold text-slate-900">提现记录</h1>
        <p className="mt-1 text-sm text-slate-500">查看商户提现申请及处理记录</p>
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
              <label className="mb-1.5 block text-xs font-medium text-slate-600">类型</label>
              <Select value={type} onValueChange={(v) => { setType(v); setPage(1); }}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {TYPE_OPTIONS.map((opt) => (
                    <SelectItem key={opt.value} value={opt.value}>
                      {opt.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
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
            <div className="w-[120px]">
              <label className="mb-1.5 block text-xs font-medium text-slate-600">审核</label>
              <Select value={reviewStatus} onValueChange={(v) => { setReviewStatus(v); setPage(1); }}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {REVIEW_STATUS_OPTIONS.map((opt) => (
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
            提现列表
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
          ) : withdrawals.length === 0 ? (
            <p className="py-12 text-center text-sm text-slate-400">暂无提现记录</p>
          ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>订单号</TableHead>
                    <TableHead>商户</TableHead>
                    <TableHead>类型</TableHead>
                    <TableHead>金额</TableHead>
                    <TableHead>手续费</TableHead>
                    <TableHead>收款信息</TableHead>
                    <TableHead>审核</TableHead>
                    <TableHead>状态</TableHead>
                    <TableHead>时间</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {withdrawals.map((item) => (
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
                        <Badge variant="secondary">
                          {item.type === "CRYPTO" ? "数币" : "法币"}
                        </Badge>
                      </TableCell>
                      <TableCell className="font-medium">
                        {formatAmount(item.amount, item.currency)} {item.currency}
                      </TableCell>
                      <TableCell className="text-slate-500">
                        {formatAmount(item.platform_fee, item.currency)}
                      </TableCell>
                      <TableCell className="max-w-[160px] truncate font-mono text-xs" title={item.to_address}>
                        {item.to_address || "—"}
                      </TableCell>
                      <TableCell>{getReviewStatusBadge(item.review_status)}</TableCell>
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
