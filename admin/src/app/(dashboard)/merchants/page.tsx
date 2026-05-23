"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
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
import { useMerchants } from "@/lib/hooks/use-merchants";

const STATUS_OPTIONS = [
  { value: "ALL", label: "全部状态" },
  { value: "PENDING_AGREEMENT", label: "待签协议" },
  { value: "PENDING_KYC", label: "待KYC" },
  { value: "ACTIVE", label: "活跃" },
  { value: "FROZEN", label: "已冻结" },
];

const KYC_STATUS_OPTIONS = [
  { value: "ALL", label: "全部KYC状态" },
  { value: "NONE", label: "未提交" },
  { value: "PENDING", label: "审核中" },
  { value: "AUTH_SUC", label: "已通过" },
  { value: "AUTH_FAIL", label: "已拒绝" },
];

function getStatusBadge(status: string) {
  switch (status) {
    case "ACTIVE":
      return <Badge className="bg-green-100 text-green-700 border-green-200">活跃</Badge>;
    case "FROZEN":
      return <Badge className="bg-red-100 text-red-700 border-red-200">已冻结</Badge>;
    case "PENDING_AGREEMENT":
      return <Badge className="bg-yellow-100 text-yellow-700 border-yellow-200">待签协议</Badge>;
    case "PENDING_KYC":
      return <Badge className="bg-blue-100 text-blue-700 border-blue-200">待KYC</Badge>;
    default:
      return <Badge variant="outline">{status}</Badge>;
  }
}

function getKycBadge(kycStatus: string) {
  switch (kycStatus) {
    case "AUTH_SUC":
      return <Badge className="bg-green-100 text-green-700 border-green-200">已通过</Badge>;
    case "AUTH_FAIL":
      return <Badge className="bg-red-100 text-red-700 border-red-200">已拒绝</Badge>;
    case "PENDING":
      return <Badge className="bg-yellow-100 text-yellow-700 border-yellow-200">审核中</Badge>;
    case "NONE":
      return <Badge className="bg-slate-100 text-slate-500 border-slate-200">未提交</Badge>;
    default:
      return <Badge variant="outline">{kycStatus}</Badge>;
  }
}

export default function MerchantsPage() {
  const router = useRouter();
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [status, setStatus] = useState<string>("ALL");
  const [kycStatus, setKycStatus] = useState<string>("ALL");
  const [search, setSearch] = useState("");
  const [searchInput, setSearchInput] = useState("");

  const { data, isLoading } = useMerchants({
    page,
    page_size: pageSize,
    status: status === "ALL" ? undefined : status,
    kyc_status: kycStatus === "ALL" ? undefined : kycStatus,
    search: search || undefined,
  });

  const totalPages = data ? Math.ceil(data.total / pageSize) : 0;

  const handleSearch = () => {
    setSearch(searchInput);
    setPage(1);
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-slate-900">商户管理</h1>

      {/* Filters */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-wrap items-center gap-4">
            <div className="flex flex-1 items-center gap-2">
              <div className="relative flex-1 max-w-sm">
                <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-slate-400" />
                <Input
                  placeholder="搜索邮箱..."
                  value={searchInput}
                  onChange={(e) => setSearchInput(e.target.value)}
                  onKeyDown={(e) => e.key === "Enter" && handleSearch()}
                  className="pl-9"
                />
              </div>
              <Button variant="outline" onClick={handleSearch}>
                搜索
              </Button>
            </div>
            <Select
              value={status}
              onValueChange={(v) => { setStatus(v); setPage(1); }}
            >
              <SelectTrigger className="w-[160px]">
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
            <Select
              value={kycStatus}
              onValueChange={(v) => { setKycStatus(v); setPage(1); }}
            >
              <SelectTrigger className="w-[160px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {KYC_STATUS_OPTIONS.map((opt) => (
                  <SelectItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Table */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">
            商户列表
            {data && (
              <span className="ml-2 text-sm font-normal text-slate-500">
                共 {data.total} 条
              </span>
            )}
          </CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-3">
              {Array.from({ length: 5 }).map((_, i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          ) : !data || data.list.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12 text-slate-400">
              <p className="text-sm">暂无商户数据</p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>邮箱</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead>KYC 状态</TableHead>
                  <TableHead>手续费模板</TableHead>
                  <TableHead>创建时间</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {data.list.map((merchant) => (
                  <TableRow
                    key={merchant.id}
                    className="cursor-pointer"
                    onClick={() => router.push(`/merchants/${merchant.id}`)}
                  >
                    <TableCell className="font-medium">{merchant.id}</TableCell>
                    <TableCell>{merchant.email}</TableCell>
                    <TableCell>{getStatusBadge(merchant.status)}</TableCell>
                    <TableCell>{getKycBadge(merchant.kyc_status)}</TableCell>
                    <TableCell className="text-slate-500">
                      {merchant.fee_template_name || "-"}
                    </TableCell>
                    <TableCell className="text-slate-500">
                      {new Date(merchant.created_at).toLocaleDateString("zh-CN")}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="mt-4 flex items-center justify-between">
              <p className="text-sm text-slate-500">
                第 {page} / {totalPages} 页
              </p>
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  disabled={page <= 1}
                  onClick={() => setPage((p) => p - 1)}
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
        </CardContent>
      </Card>
    </div>
  );
}
