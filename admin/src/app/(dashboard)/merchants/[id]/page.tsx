"use client";

import { use, useEffect, useState } from "react";
import Link from "next/link";
import { ArrowLeft, Shield, ShieldOff, CheckCircle, XCircle, RefreshCw } from "lucide-react";
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
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  useMerchant,
  useUpdateMerchantStatus,
  useUpdateMerchantFeeTemplate,
  useUpdateMerchantSupportedCurrencies,
  useSyncKUNBalances,
  useSyncDeposits,
  useMerchantLedger,
  useApproveKyc,
  useRejectKyc,
} from "@/lib/hooks/use-merchants";
import { useFeeTemplates } from "@/lib/hooks/use-fee-templates";
import { useMerchantDeposits } from "@/lib/hooks/use-deposits";
import { useMerchantWithdrawals } from "@/lib/hooks/use-withdrawals";
import { useAdminExchanges, useSyncExchangeStatus } from "@/lib/hooks/use-exchanges";
import { formatAmount } from "@/lib/utils/format";
import { toast } from "@/hooks/use-toast";
import type { KUNWalletBalance, MerchantWallet } from "@/types/merchant";
import { Switch } from "@/components/ui/switch";
import { formatChainLabel } from "@/lib/utils/chain";

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

function walletRowKey(accountType: string, currency: string) {
  return `${accountType}:${currency}`;
}

function buildMergedWalletRows(
  platform: MerchantWallet[] | undefined,
  kun: KUNWalletBalance[]
) {
  type Row = {
    account_type: string;
    currency: string;
    balance?: string;
    frozen_balance?: string;
    available_balance?: string;
    kun_balance?: string;
  };
  const rows = new Map<string, Row>();
  for (const wallet of platform ?? []) {
    const key = walletRowKey(wallet.account_type, wallet.currency);
    rows.set(key, {
      account_type: wallet.account_type,
      currency: wallet.currency,
      balance: wallet.balance,
      frozen_balance: wallet.frozen_balance,
      available_balance: wallet.available_balance,
    });
  }
  for (const wallet of kun) {
    const key = walletRowKey(wallet.account_type, wallet.currency);
    const existing = rows.get(key);
    if (existing) {
      existing.kun_balance = wallet.balance;
    } else {
      rows.set(key, {
        account_type: wallet.account_type,
        currency: wallet.currency,
        kun_balance: wallet.balance,
      });
    }
  }
  return Array.from(rows.values()).sort((a, b) => {
    if (a.account_type !== b.account_type) {
      return a.account_type.localeCompare(b.account_type);
    }
    return a.currency.localeCompare(b.currency);
  });
}

function getDepositStatusBadge(status: string) {
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

function getWithdrawalStatusBadge(status: string) {
  switch (status) {
    case "COMPLETED":
      return <Badge className="bg-green-100 text-green-700 border-green-200">已完成</Badge>;
    case "PROCESSING":
    case "PENDING":
      return <Badge className="bg-amber-100 text-amber-700 border-amber-200">处理中</Badge>;
    case "FAILED":
      return <Badge className="bg-red-100 text-red-700 border-red-200">失败</Badge>;
    default:
      return <Badge variant="outline">{status}</Badge>;
  }
}

function getExchangeStatusBadge(status: string) {
  return getWithdrawalStatusBadge(status);
}

function canSyncExchange(status: string) {
  return status === "PROCESSING" || status === "PENDING" || status === "FAILED";
}

function getReviewStatusBadge(status: string) {
  switch (status) {
    case "PENDING_REVIEW":
      return <Badge className="bg-amber-100 text-amber-700 border-amber-200">待审核</Badge>;
    case "APPROVED":
      return <Badge className="bg-green-100 text-green-700 border-green-200">已通过</Badge>;
    case "REJECTED":
      return <Badge className="bg-red-100 text-red-700 border-red-200">已拒绝</Badge>;
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
  const syncKUNBalancesMutation = useSyncKUNBalances();
  const syncDepositsMutation = useSyncDeposits();
  const { data: ledgerData, isLoading: ledgerLoading } = useMerchantLedger(id, {
    page: 1,
    page_size: 50,
  });
  const { data: depositsData, isLoading: depositsLoading } = useMerchantDeposits(id);
  const { data: withdrawalsData, isLoading: withdrawalsLoading } = useMerchantWithdrawals(id);
  const { data: exchangesData, isLoading: exchangesLoading } = useAdminExchanges({ merchantId: id });
  const syncExchangeMutation = useSyncExchangeStatus();
  const approveKycMutation = useApproveKyc();
  const rejectKycMutation = useRejectKyc();

  const [freezeDialogOpen, setFreezeDialogOpen] = useState(false);
  const [kycRejectDialogOpen, setKycRejectDialogOpen] = useState(false);
  const [kycApproveDialogOpen, setKycApproveDialogOpen] = useState(false);
  const [rejectReason, setRejectReason] = useState("");
  const [selectedTemplateId, setSelectedTemplateId] = useState<string>("");
  const [selectedCrypto, setSelectedCrypto] = useState<string[]>([]);
  const [selectedFiat, setSelectedFiat] = useState<string[]>([]);
  const [selectedChains, setSelectedChains] = useState<Record<string, string[]>>({});
  const [defaultChains, setDefaultChains] = useState<Record<string, string>>({});
  const [kunBalances, setKunBalances] = useState<KUNWalletBalance[]>([]);
  const [kunSyncedAt, setKunSyncedAt] = useState<string | null>(null);
  const [depositsSyncedAt, setDepositsSyncedAt] = useState<string | null>(null);
  const [syncingExchangeId, setSyncingExchangeId] = useState<number | null>(null);

  useEffect(() => {
    if (!merchant) return;
    setSelectedCrypto(merchant.supported_crypto_currencies ?? []);
    setSelectedFiat(merchant.supported_fiat_currencies ?? []);
    setSelectedChains(merchant.supported_crypto_chains ?? {});
    setDefaultChains(merchant.default_crypto_chains ?? {});
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
    if (!selectedChains[currency]?.length) {
      const available = merchant?.available_crypto_chains?.[currency] ?? [];
      if (available.length > 0) {
        setSelectedChains({ ...selectedChains, [currency]: available });
        setDefaultChains({
          ...defaultChains,
          [currency]: merchant?.available_default_chains?.[currency] ?? available[0],
        });
      }
    }
  };

  const toggleChain = (currency: string, chain: string) => {
    const current = selectedChains[currency] ?? [];
    const next = current.includes(chain)
      ? current.filter((item) => item !== chain)
      : [...current, chain];
    setSelectedChains({ ...selectedChains, [currency]: next });
    if (!next.includes(defaultChains[currency] ?? "") && next.length > 0) {
      setDefaultChains({ ...defaultChains, [currency]: next[0] });
    }
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
    for (const currency of selectedCrypto) {
      if ((selectedChains[currency] ?? []).length === 0) {
        toast({
          title: "保存失败",
          description: `${currency} 至少需要选择一条支持链`,
          variant: "destructive",
        });
        return;
      }
    }

    updateSupportedCurrenciesMutation.mutate(
      {
        id,
        crypto_currencies: selectedCrypto,
        fiat_currencies: selectedFiat,
        crypto_chains: selectedChains,
        default_chains: defaultChains,
      },
      {
        onSuccess: () => {
          toast({ title: "币种与链配置已保存" });
        },
        onError: (error) => {
          toast({ title: "保存失败", description: error.message, variant: "destructive" });
        },
      }
    );
  };

  const handleSyncKUNBalances = () => {
    syncKUNBalancesMutation.mutate(id, {
      onSuccess: (data) => {
        setKunBalances(data.kun_balances);
        setKunSyncedAt(data.synced_at);
        toast({ title: "KUN 余额已同步" });
      },
      onError: (error) => {
        toast({ title: "同步失败", description: error.message, variant: "destructive" });
      },
    });
  };

  const handleSyncDeposits = () => {
    syncDepositsMutation.mutate(id, {
      onSuccess: (data) => {
        setDepositsSyncedAt(data.synced_at);
        toast({
          title: "充值记录已同步",
          description: `新增 ${data.synced_count} 条，更新 ${data.updated_count} 条，共拉取 ${data.total_fetched} 条`,
        });
      },
      onError: (error) => {
        toast({ title: "同步失败", description: error.message, variant: "destructive" });
      },
    });
  };

  const handleSyncExchange = (orderId: number) => {
    setSyncingExchangeId(orderId);
    syncExchangeMutation.mutate(orderId, {
      onSuccess: (result) => {
        toast({
          title: result.updated ? "兑换状态已同步" : "状态未变化",
          description: result.updated
            ? result.status === "FAILED" && result.fail_reason
              ? `订单 #${orderId} 失败原因已更新：${result.fail_reason}`
              : `订单 #${orderId} 已更新为 ${result.status}${result.status === "FAILED" && result.fail_reason ? `：${result.fail_reason}` : ""}`
            : `KUN 状态：${result.kun_status || "—"}`,
        });
      },
      onError: (error) => {
        toast({ title: "同步失败", description: error.message, variant: "destructive" });
      },
      onSettled: () => setSyncingExchangeId(null),
    });
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

  const mergedWalletRows = buildMergedWalletRows(merchant.wallets, kunBalances);

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center gap-4">
        <Button variant="ghost" size="icon" asChild>
          <Link href="/merchants">
            <ArrowLeft className="size-4" />
          </Link>
        </Button>
        <div className="min-w-0 flex-1">
          <h1 className="truncate text-2xl font-bold text-slate-900">{merchant.email}</h1>
          <div className="mt-2 flex flex-wrap items-center gap-2">
            {getStatusBadge(merchant.status)}
            {getKycBadge(merchant.kyc_status)}
            {merchant.fee_template_name && (
              <Badge variant="outline">{merchant.fee_template_name}</Badge>
            )}
          </div>
        </div>
      </div>

      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">概览</TabsTrigger>
          <TabsTrigger value="funds">资金</TabsTrigger>
          <TabsTrigger value="transactions">交易记录</TabsTrigger>
          <TabsTrigger value="settings">配置</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
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
        </TabsContent>

        <TabsContent value="funds" className="space-y-6">
      {/* Platform + KUN Wallets */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0">
          <div>
            <CardTitle className="text-base">平台余额</CardTitle>
            {kunSyncedAt && (
              <p className="mt-1 text-xs text-slate-500">
                KUN 最近同步：{new Date(kunSyncedAt).toLocaleString("zh-CN")}
              </p>
            )}
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={handleSyncKUNBalances}
            disabled={syncKUNBalancesMutation.isPending || !merchant.kun_sub_customer_no}
          >
            <RefreshCw className={`mr-2 size-4 ${syncKUNBalancesMutation.isPending ? "animate-spin" : ""}`} />
            {syncKUNBalancesMutation.isPending ? "同步中..." : "同步 KUN 余额"}
          </Button>
        </CardHeader>
        <CardContent>
          {mergedWalletRows.length > 0 ? (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>账户类型</TableHead>
                  <TableHead>币种</TableHead>
                  <TableHead className="text-right">平台余额</TableHead>
                  <TableHead className="text-right">冻结金额</TableHead>
                  <TableHead className="text-right">可用余额</TableHead>
                  <TableHead className="text-right">KUN 余额</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {mergedWalletRows.map((row) => (
                  <TableRow key={walletRowKey(row.account_type, row.currency)}>
                    <TableCell className="font-medium">{row.account_type}</TableCell>
                    <TableCell>{row.currency}</TableCell>
                    <TableCell className="text-right">
                      {row.balance != null
                        ? formatAmount(row.balance, row.currency)
                        : "—"}
                    </TableCell>
                    <TableCell className="text-right">
                      {row.frozen_balance != null
                        ? formatAmount(row.frozen_balance, row.currency)
                        : "—"}
                    </TableCell>
                    <TableCell className="text-right">
                      {row.available_balance != null
                        ? formatAmount(row.available_balance, row.currency)
                        : "—"}
                    </TableCell>
                    <TableCell className="text-right">
                      {row.kun_balance != null
                        ? formatAmount(row.kun_balance, row.currency)
                        : "—"}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          ) : (
            <p className="py-4 text-center text-sm text-slate-400">暂无余额数据</p>
          )}
          {!merchant.kun_sub_customer_no && (
            <p className="mt-3 text-center text-xs text-slate-400">
              商户尚未完成 KUN 入网，无法同步 KUN 余额
            </p>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">资金变化</CardTitle>
        </CardHeader>
        <CardContent>
          {ledgerLoading ? (
            <Skeleton className="h-24 w-full" />
          ) : (ledgerData?.entries?.length ?? 0) > 0 ? (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>时间</TableHead>
                  <TableHead>类型</TableHead>
                  <TableHead>账户</TableHead>
                  <TableHead>业务</TableHead>
                  <TableHead>币种</TableHead>
                  <TableHead className="text-right">变动金额</TableHead>
                  <TableHead className="text-right">余额后</TableHead>
                  <TableHead className="text-right">冻结后</TableHead>
                  <TableHead>订单号</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {ledgerData?.entries.map((entry) => {
                  const entryLabel =
                    entry.entry_type === "CREDIT"
                      ? "入账"
                      : entry.entry_type === "FREEZE"
                        ? "冻结"
                        : entry.entry_type === "UNFREEZE"
                          ? "解冻"
                          : entry.entry_type === "DEDUCT_FROZEN"
                            ? "扣款"
                            : entry.entry_type;
                  const bizLabel =
                    entry.biz_type === "DEPOSIT"
                      ? "充值"
                      : entry.biz_type === "WITHDRAWAL"
                        ? "提现"
                        : entry.biz_type === "EXCHANGE"
                          ? "兑换"
                          : entry.biz_type === "TRANSFER"
                            ? "划转"
                            : entry.biz_type || "—";
                  const signedPrefix =
                    entry.entry_type === "CREDIT" || entry.entry_type === "UNFREEZE"
                      ? "+"
                      : entry.entry_type === "DEDUCT_FROZEN" || entry.entry_type === "FREEZE"
                        ? "-"
                        : "";
                  return (
                    <TableRow key={entry.id}>
                      <TableCell className="whitespace-nowrap text-sm">
                        {new Date(entry.created_at).toLocaleString("zh-CN")}
                      </TableCell>
                      <TableCell>{entryLabel}</TableCell>
                      <TableCell>
                        {entry.account_type === "FUNDING" ? "资金" : entry.account_type === "TRADING" ? "交易" : entry.account_type}
                      </TableCell>
                      <TableCell>{bizLabel}</TableCell>
                      <TableCell>{entry.currency}</TableCell>
                      <TableCell className="text-right font-medium">
                        {signedPrefix}
                        {formatAmount(entry.amount, entry.currency)}
                      </TableCell>
                      <TableCell className="text-right">
                        {formatAmount(entry.balance_after, entry.currency)}
                      </TableCell>
                      <TableCell className="text-right">
                        {formatAmount(entry.frozen_after, entry.currency)}
                      </TableCell>
                      <TableCell className="max-w-[160px] truncate font-mono text-xs text-slate-500">
                        {entry.platform_order_id || "—"}
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          ) : (
            <p className="py-4 text-center text-sm text-slate-400">暂无资金变化记录</p>
          )}
        </CardContent>
      </Card>
        </TabsContent>

        <TabsContent value="transactions" className="space-y-4">
          <Tabs defaultValue="deposit">
            <TabsList>
              <TabsTrigger value="deposit">充值</TabsTrigger>
              <TabsTrigger value="withdrawal">提现</TabsTrigger>
              <TabsTrigger value="exchange">兑换</TabsTrigger>
              <TabsTrigger value="transfer">划转</TabsTrigger>
            </TabsList>

            <TabsContent value="deposit" className="mt-4">
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0">
                  <div>
                    <CardTitle className="text-base">充值记录</CardTitle>
                    {depositsSyncedAt && (
                      <p className="mt-1 text-xs text-slate-500">
                        最近同步：{new Date(depositsSyncedAt).toLocaleString("zh-CN")}
                      </p>
                    )}
                  </div>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleSyncDeposits}
                    disabled={syncDepositsMutation.isPending || !merchant.kun_sub_customer_no}
                  >
                    <RefreshCw className={`mr-2 size-4 ${syncDepositsMutation.isPending ? "animate-spin" : ""}`} />
                    {syncDepositsMutation.isPending ? "同步中..." : "同步充值记录"}
                  </Button>
                </CardHeader>
                <CardContent>
                  {depositsLoading ? (
                    <Skeleton className="h-24 w-full" />
                  ) : (depositsData?.deposits?.length ?? 0) > 0 ? (
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead>时间</TableHead>
                          <TableHead>币种</TableHead>
                          <TableHead>网络</TableHead>
                          <TableHead className="text-right">金额</TableHead>
                          <TableHead>状态</TableHead>
                          <TableHead>交易哈希</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {depositsData?.deposits.map((deposit) => (
                          <TableRow key={deposit.id}>
                            <TableCell className="whitespace-nowrap text-sm">
                              {new Date(deposit.created_at).toLocaleString("zh-CN")}
                            </TableCell>
                            <TableCell>{deposit.currency}</TableCell>
                            <TableCell>{deposit.network}</TableCell>
                            <TableCell className="text-right">
                              {formatAmount(deposit.amount, deposit.currency)}
                            </TableCell>
                            <TableCell>{getDepositStatusBadge(deposit.status)}</TableCell>
                            <TableCell className="max-w-[180px] truncate font-mono text-xs text-slate-500">
                              {deposit.tx_hash || "—"}
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  ) : (
                    <p className="py-4 text-center text-sm text-slate-400">
                      暂无充值记录，点击「同步充值记录」从 KUN 拉取
                    </p>
                  )}
                  {!merchant.kun_sub_customer_no && (
                    <p className="mt-3 text-center text-xs text-slate-400">
                      商户尚未完成 KUN 入网，无法同步充值记录
                    </p>
                  )}
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="withdrawal" className="mt-4">
              <Card>
                <CardHeader>
                  <CardTitle className="text-base">提现记录</CardTitle>
                </CardHeader>
                <CardContent>
                  {withdrawalsLoading ? (
                    <Skeleton className="h-24 w-full" />
                  ) : (withdrawalsData?.withdrawals?.length ?? 0) > 0 ? (
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead>时间</TableHead>
                          <TableHead>类型</TableHead>
                          <TableHead>币种</TableHead>
                          <TableHead className="text-right">金额</TableHead>
                          <TableHead>审核</TableHead>
                          <TableHead>状态</TableHead>
                          <TableHead>收款信息</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {withdrawalsData?.withdrawals.map((withdrawal) => (
                          <TableRow key={withdrawal.id}>
                            <TableCell className="whitespace-nowrap text-sm">
                              {new Date(withdrawal.created_at).toLocaleString("zh-CN")}
                            </TableCell>
                            <TableCell>
                              <Badge variant="secondary">
                                {withdrawal.type === "CRYPTO" ? "数币" : "法币"}
                              </Badge>
                            </TableCell>
                            <TableCell>{withdrawal.currency}</TableCell>
                            <TableCell className="text-right">
                              {formatAmount(withdrawal.amount, withdrawal.currency)}
                            </TableCell>
                            <TableCell>{getReviewStatusBadge(withdrawal.review_status)}</TableCell>
                            <TableCell>{getWithdrawalStatusBadge(withdrawal.status)}</TableCell>
                            <TableCell className="max-w-[180px] truncate font-mono text-xs text-slate-500">
                              {withdrawal.to_address || "—"}
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  ) : (
                    <p className="py-4 text-center text-sm text-slate-400">暂无提现记录</p>
                  )}
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="exchange" className="mt-4">
              <Card>
                <CardHeader>
                  <CardTitle className="text-base">兑换记录</CardTitle>
                </CardHeader>
                <CardContent>
                  {exchangesLoading ? (
                    <Skeleton className="h-24 w-full" />
                  ) : (exchangesData?.exchanges?.length ?? 0) > 0 ? (
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead>时间</TableHead>
                          <TableHead>卖出</TableHead>
                          <TableHead>买入</TableHead>
                          <TableHead className="text-right">手续费</TableHead>
                          <TableHead>状态</TableHead>
                          <TableHead>失败原因</TableHead>
                          <TableHead className="text-right">操作</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {exchangesData?.exchanges.map((exchange) => (
                          <TableRow key={exchange.id}>
                            <TableCell className="whitespace-nowrap text-sm">
                              {new Date(exchange.created_at).toLocaleString("zh-CN")}
                            </TableCell>
                            <TableCell>
                              {formatAmount(exchange.from_amount, exchange.from_currency)} {exchange.from_currency}
                            </TableCell>
                            <TableCell>
                              {exchange.to_amount
                                ? `${formatAmount(exchange.to_amount, exchange.to_currency)} ${exchange.to_currency}`
                                : "—"}
                            </TableCell>
                            <TableCell className="text-right text-slate-500">
                              {formatAmount(exchange.platform_fee, exchange.from_currency)}
                            </TableCell>
                            <TableCell>{getExchangeStatusBadge(exchange.status)}</TableCell>
                            <TableCell className="max-w-[180px] text-xs text-red-600">
                              {exchange.status === "FAILED" && exchange.fail_reason ? exchange.fail_reason : "—"}
                            </TableCell>
                            <TableCell className="text-right">
                              {canSyncExchange(exchange.status) ? (
                                <Button
                                  variant="outline"
                                  size="sm"
                                  disabled={syncingExchangeId === exchange.id}
                                  onClick={() => handleSyncExchange(exchange.id)}
                                >
                                  <RefreshCw className={`mr-1 h-3.5 w-3.5 ${syncingExchangeId === exchange.id ? "animate-spin" : ""}`} />
                                  同步
                                </Button>
                              ) : (
                                <span className="text-xs text-slate-400">—</span>
                              )}
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  ) : (
                    <p className="py-4 text-center text-sm text-slate-400">暂无兑换记录</p>
                  )}
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="transfer" className="mt-4">
              <Card>
                <CardContent className="py-10 text-center text-sm text-slate-400">
                  划转记录即将上线
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        </TabsContent>

        <TabsContent value="settings" className="space-y-6">
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

          <Separator />

          <div className="space-y-4">
            <p className="text-sm font-medium text-slate-700">数字货币支持链与默认链</p>
            {selectedCrypto.map((currency) => {
              const available = merchant.available_crypto_chains?.[currency] ?? [];
              const enabled = selectedChains[currency] ?? [];
              return (
                <div key={currency} className="space-y-3 rounded-lg border p-4">
                  <div className="flex flex-wrap items-center justify-between gap-3">
                    <p className="text-sm font-semibold text-slate-800">{currency}</p>
                    <div className="flex items-center gap-2">
                      <span className="text-xs text-slate-500">默认链</span>
                      <Select
                        value={defaultChains[currency] || enabled[0] || ""}
                        onValueChange={(value) =>
                          setDefaultChains({ ...defaultChains, [currency]: value })
                        }
                      >
                        <SelectTrigger className="w-[200px]">
                          <SelectValue placeholder="选择默认链" />
                        </SelectTrigger>
                        <SelectContent>
                          {enabled.map((chain) => (
                            <SelectItem key={chain} value={chain}>
                              {formatChainLabel(chain)}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                  </div>
                  <div className="grid gap-2 sm:grid-cols-2 md:grid-cols-3">
                    {available.map((chain) => (
                      <div key={chain} className="flex items-center justify-between rounded-md border px-3 py-2">
                        <span className="text-sm">{formatChainLabel(chain)}</span>
                        <Switch
                          checked={enabled.includes(chain)}
                          onCheckedChange={() => toggleChain(currency, chain)}
                        />
                      </div>
                    ))}
                  </div>
                </div>
              );
            })}
          </div>

          <div className="flex justify-end">
            <Button
              onClick={handleSaveSupportedCurrencies}
              disabled={updateSupportedCurrenciesMutation.isPending}
            >
              {updateSupportedCurrenciesMutation.isPending ? "保存中..." : "保存币种与链配置"}
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Fee template */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">手续费模板</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div>
              <p className="text-sm text-slate-500">当前模板</p>
              <p className="font-medium">{merchant.fee_template_name || "未分配"}</p>
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
        </CardContent>
      </Card>
        </TabsContent>
      </Tabs>

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
