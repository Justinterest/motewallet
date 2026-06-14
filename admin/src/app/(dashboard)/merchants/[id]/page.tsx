"use client";

import { use, useEffect, useState } from "react";
import Link from "next/link";
import { ArrowLeft, Shield, ShieldOff, CheckCircle, XCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  useMerchant,
  useUpdateMerchantStatus,
  useUpdateMerchantFeeTemplate,
  useUpdateMerchantSupportedCurrencies,
  useApproveKyc,
  useRejectKyc,
} from "@/lib/hooks/use-merchants";
import { useFeeTemplates } from "@/lib/hooks/use-fee-templates";
import { formatAmount } from "@/lib/utils/format";
import { toast } from "@/hooks/use-toast";
import { Switch } from "@/components/ui/switch";

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

export default function MerchantDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id: idStr } = use(params);
  const id = parseInt(idStr, 10);
  const { data: merchant, isLoading } = useMerchant(id);
  const { data: templates } = useFeeTemplates();
  const updateStatusMutation = useUpdateMerchantStatus();
  const updateFeeTemplateMutation = useUpdateMerchantFeeTemplate();
  const updateSupportedCurrenciesMutation = useUpdateMerchantSupportedCurrencies();
  const approveKycMutation = useApproveKyc();
  const rejectKycMutation = useRejectKyc();

  const [freezeDialogOpen, setFreezeDialogOpen] = useState(false);
  const [kycRejectDialogOpen, setKycRejectDialogOpen] = useState(false);
  const [kycApproveDialogOpen, setKycApproveDialogOpen] = useState(false);
  const [rejectReason, setRejectReason] = useState("");
  const [selectedTemplateId, setSelectedTemplateId] = useState<string>("");
  const [selectedCrypto, setSelectedCrypto] = useState<string[]>([]);
  const [selectedFiat, setSelectedFiat] = useState<string[]>([]);

  useEffect(() => {
    if (!merchant) return;
    setSelectedCrypto(merchant.supported_crypto_currencies ?? []);
    setSelectedFiat(merchant.supported_fiat_currencies ?? []);
  }, [merchant]);

  const handleFreezeToggle = () => {
    const newStatus = merchant?.status === "FROZEN" ? "ACTIVE" : "FROZEN";
    updateStatusMutation.mutate(
      { id, status: newStatus },
      {
        onSuccess: () => {
          toast({ title: newStatus === "FROZEN" ? "已冻结" : "已解冻" });
          setFreezeDialogOpen(false);
        },
        onError: (error) => {
          toast({ title: "操作失败", description: error.message, variant: "destructive" });
        },
      }
    );
  };

  const handleAssignTemplate = () => {
    if (!selectedTemplateId) return;
    updateFeeTemplateMutation.mutate(
      { id, fee_template_id: parseInt(selectedTemplateId, 10) },
      {
        onSuccess: () => {
          toast({ title: "模板已分配" });
          setSelectedTemplateId("");
        },
        onError: (error) => {
          toast({ title: "分配失败", description: error.message, variant: "destructive" });
        },
      }
    );
  };

  const toggleCurrency = (
    currency: string,
    selected: string[],
    setSelected: (value: string[]) => void
  ) => {
    if (selected.includes(currency)) {
      setSelected(selected.filter((item) => item !== currency));
      return;
    }
    setSelected([...selected, currency]);
  };

  const handleSaveSupportedCurrencies = () => {
    if (selectedCrypto.length === 0 || selectedFiat.length === 0) {
      toast({
        title: "保存失败",
        description: "至少需保留一种数字币和一种法币",
        variant: "destructive",
      });
      return;
    }

    updateSupportedCurrenciesMutation.mutate(
      { id, crypto_currencies: selectedCrypto, fiat_currencies: selectedFiat },
      {
        onSuccess: () => {
          toast({ title: "币种配置已保存" });
        },
        onError: (error) => {
          toast({ title: "保存失败", description: error.message, variant: "destructive" });
        },
      }
    );
  };

  const handleApproveKyc = () => {
    approveKycMutation.mutate(id, {
      onSuccess: () => {
        toast({ title: "KYC 已通过" });
        setKycApproveDialogOpen(false);
      },
      onError: (error) => {
        toast({ title: "操作失败", description: error.message, variant: "destructive" });
      },
    });
  };

  const handleRejectKyc = () => {
    if (!rejectReason.trim()) return;
    rejectKycMutation.mutate(
      { id, reason: rejectReason },
      {
        onSuccess: () => {
          toast({ title: "KYC 已拒绝" });
          setKycRejectDialogOpen(false);
          setRejectReason("");
        },
        onError: (error) => {
          toast({ title: "操作失败", description: error.message, variant: "destructive" });
        },
      }
    );
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (!merchant) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/merchants">
              <ArrowLeft className="size-4" />
            </Link>
          </Button>
          <h1 className="text-2xl font-bold text-slate-900">商户不存在</h1>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" asChild>
          <Link href="/merchants">
            <ArrowLeft className="size-4" />
          </Link>
        </Button>
        <h1 className="text-2xl font-bold text-slate-900">商户详情</h1>
      </div>

      {/* Basic Info */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">基本信息</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 gap-4 md:grid-cols-3">
            <div>
              <p className="text-sm text-slate-500">邮箱</p>
              <p className="font-medium">{merchant.email}</p>
            </div>
            <div>
              <p className="text-sm text-slate-500">状态</p>
              <div className="mt-1">{getStatusBadge(merchant.status)}</div>
            </div>
            <div>
              <p className="text-sm text-slate-500">KYC 状态</p>
              <div className="mt-1">{getKycBadge(merchant.kyc_status)}</div>
            </div>
            <div>
              <p className="text-sm text-slate-500">手续费模板</p>
              <p className="font-medium">{merchant.fee_template_name || "未分配"}</p>
            </div>
            <div>
              <p className="text-sm text-slate-500">创建时间</p>
              <p className="font-medium">
                {new Date(merchant.created_at).toLocaleString("zh-CN")}
              </p>
            </div>
            {merchant.kyc_fail_reason && (
              <div className="col-span-full">
                <p className="text-sm text-slate-500">KYC 拒绝原因</p>
                <p className="font-medium text-red-600">{merchant.kyc_fail_reason}</p>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Wallets */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">钱包余额</CardTitle>
        </CardHeader>
        <CardContent>
          {merchant.wallets && merchant.wallets.length > 0 ? (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>账户类型</TableHead>
                  <TableHead>币种</TableHead>
                  <TableHead className="text-right">余额</TableHead>
                  <TableHead className="text-right">冻结金额</TableHead>
                  <TableHead className="text-right">可用余额</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {merchant.wallets.map((wallet, index) => (
                  <TableRow key={index}>
                    <TableCell className="font-medium">{wallet.account_type}</TableCell>
                    <TableCell>{wallet.currency}</TableCell>
                    <TableCell className="text-right">
                      {formatAmount(wallet.balance, wallet.currency)}
                    </TableCell>
                    <TableCell className="text-right">
                      {formatAmount(wallet.frozen_balance, wallet.currency)}
                    </TableCell>
                    <TableCell className="text-right">
                      {formatAmount(wallet.available_balance, wallet.currency)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          ) : (
            <p className="py-4 text-center text-sm text-slate-400">暂无钱包数据</p>
          )}
        </CardContent>
      </Card>

      {/* Supported Currencies */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">支持币种</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <div>
            <p className="mb-3 text-sm font-medium text-slate-700">数字货币</p>
            <div className="grid gap-3 sm:grid-cols-2 md:grid-cols-3">
              {(merchant.available_crypto_currencies ?? []).map((currency) => (
                <div key={currency} className="flex items-center justify-between rounded-lg border px-3 py-2">
                  <span className="text-sm font-medium">{currency}</span>
                  <Switch
                    checked={selectedCrypto.includes(currency)}
                    onCheckedChange={() => toggleCurrency(currency, selectedCrypto, setSelectedCrypto)}
                  />
                </div>
              ))}
            </div>
          </div>

          <Separator />

          <div>
            <p className="mb-3 text-sm font-medium text-slate-700">法币</p>
            <div className="grid gap-3 sm:grid-cols-2 md:grid-cols-3">
              {(merchant.available_fiat_currencies ?? []).map((currency) => (
                <div key={currency} className="flex items-center justify-between rounded-lg border px-3 py-2">
                  <span className="text-sm font-medium">{currency}</span>
                  <Switch
                    checked={selectedFiat.includes(currency)}
                    onCheckedChange={() => toggleCurrency(currency, selectedFiat, setSelectedFiat)}
                  />
                </div>
              ))}
            </div>
          </div>

          <div className="flex justify-end">
            <Button
              onClick={handleSaveSupportedCurrencies}
              disabled={updateSupportedCurrenciesMutation.isPending}
            >
              {updateSupportedCurrenciesMutation.isPending ? "保存中..." : "保存币种配置"}
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Timeline */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">时间线</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-sm text-slate-500">协议签署时间</span>
              <span className="text-sm font-medium">
                {merchant.agreement_signed_at
                  ? new Date(merchant.agreement_signed_at).toLocaleString("zh-CN")
                  : "未签署"}
              </span>
            </div>
            <Separator />
            <div className="flex items-center justify-between">
              <span className="text-sm text-slate-500">KYC 提交时间</span>
              <span className="text-sm font-medium">
                {merchant.kyc_submitted_at
                  ? new Date(merchant.kyc_submitted_at).toLocaleString("zh-CN")
                  : "未提交"}
              </span>
            </div>
            <Separator />
            <div className="flex items-center justify-between">
              <span className="text-sm text-slate-500">KYC 完成时间</span>
              <span className="text-sm font-medium">
                {merchant.kyc_completed_at
                  ? new Date(merchant.kyc_completed_at).toLocaleString("zh-CN")
                  : "未完成"}
              </span>
            </div>
            <Separator />
            <div className="flex items-center justify-between">
              <span className="text-sm text-slate-500">冻结时间</span>
              <span className="text-sm font-medium">
                {merchant.frozen_at
                  ? new Date(merchant.frozen_at).toLocaleString("zh-CN")
                  : "未冻结"}
              </span>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Actions */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">操作</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Freeze / Unfreeze */}
          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium">
                {merchant.status === "FROZEN" ? "解冻商户" : "冻结商户"}
              </p>
              <p className="text-sm text-slate-500">
                {merchant.status === "FROZEN"
                  ? "解冻后商户可正常使用系统"
                  : "冻结后商户将无法进行任何操作"}
              </p>
            </div>
            <Button
              variant={merchant.status === "FROZEN" ? "default" : "destructive"}
              onClick={() => setFreezeDialogOpen(true)}
            >
              {merchant.status === "FROZEN" ? (
                <>
                  <Shield className="size-4" />
                  解冻
                </>
              ) : (
                <>
                  <ShieldOff className="size-4" />
                  冻结
                </>
              )}
            </Button>
          </div>

          <Separator />

          {/* Assign fee template */}
          <div className="flex items-center justify-between gap-4">
            <div>
              <p className="font-medium">分配手续费模板</p>
              <p className="text-sm text-slate-500">
                当前模板: {merchant.fee_template_name || "未分配"}
              </p>
            </div>
            <div className="flex items-center gap-2">
              <Select value={selectedTemplateId} onValueChange={setSelectedTemplateId}>
                <SelectTrigger className="w-[200px]">
                  <SelectValue placeholder="选择模板" />
                </SelectTrigger>
                <SelectContent>
                  {templates?.map((t) => (
                    <SelectItem key={t.id} value={String(t.id)}>
                      {t.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Button
                variant="outline"
                disabled={!selectedTemplateId || updateFeeTemplateMutation.isPending}
                onClick={handleAssignTemplate}
              >
                {updateFeeTemplateMutation.isPending ? "分配中..." : "分配"}
              </Button>
            </div>
          </div>

          {/* KYC Actions */}
          {merchant.kyc_status === "PENDING" && (
            <>
              <Separator />
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">KYC 审核</p>
                  <p className="text-sm text-slate-500">商户已提交 KYC 材料，等待审核</p>
                </div>
                <div className="flex items-center gap-2">
                  <Button
                    variant="default"
                    className="bg-green-600 hover:bg-green-700"
                    onClick={() => setKycApproveDialogOpen(true)}
                  >
                    <CheckCircle className="size-4" />
                    通过
                  </Button>
                  <Button
                    variant="destructive"
                    onClick={() => setKycRejectDialogOpen(true)}
                  >
                    <XCircle className="size-4" />
                    拒绝
                  </Button>
                </div>
              </div>
            </>
          )}
        </CardContent>
      </Card>

      {/* Freeze/Unfreeze Dialog */}
      <AlertDialog open={freezeDialogOpen} onOpenChange={setFreezeDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              {merchant.status === "FROZEN" ? "确认解冻" : "确认冻结"}
            </AlertDialogTitle>
            <AlertDialogDescription>
              {merchant.status === "FROZEN"
                ? `确定要解冻商户 ${merchant.email} 吗？解冻后商户可正常使用系统。`
                : `确定要冻结商户 ${merchant.email} 吗？冻结后商户将无法进行任何操作。`}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleFreezeToggle}
              className={merchant.status === "FROZEN" ? "" : "bg-red-600 hover:bg-red-700"}
            >
              {updateStatusMutation.isPending
                ? "处理中..."
                : merchant.status === "FROZEN"
                  ? "确认解冻"
                  : "确认冻结"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* KYC Approve Dialog */}
      <AlertDialog open={kycApproveDialogOpen} onOpenChange={setKycApproveDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>确认通过 KYC</AlertDialogTitle>
            <AlertDialogDescription>
              确定要通过商户 {merchant.email} 的 KYC 审核吗？
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleApproveKyc}
              className="bg-green-600 hover:bg-green-700"
            >
              {approveKycMutation.isPending ? "处理中..." : "确认通过"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* KYC Reject Dialog */}
      <Dialog open={kycRejectDialogOpen} onOpenChange={setKycRejectDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>拒绝 KYC</DialogTitle>
            <DialogDescription>
              请填写拒绝原因，该原因将展示给商户。
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-2">
            <Label>拒绝原因</Label>
            <Textarea
              placeholder="请输入拒绝原因..."
              value={rejectReason}
              onChange={(e) => setRejectReason(e.target.value)}
            />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setKycRejectDialogOpen(false)}>
              取消
            </Button>
            <Button
              variant="destructive"
              disabled={!rejectReason.trim() || rejectKycMutation.isPending}
              onClick={handleRejectKyc}
            >
              {rejectKycMutation.isPending ? "处理中..." : "确认拒绝"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
