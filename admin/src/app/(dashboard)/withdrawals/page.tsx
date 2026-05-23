"use client";

import { useState } from "react";
import { CheckCircle, XCircle, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
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
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import {
  usePendingWithdrawals,
  useApproveWithdrawal,
  useRejectWithdrawal,
} from "@/lib/hooks/use-withdrawals";
import { useToast } from "@/hooks/use-toast";

export default function WithdrawalsPage() {
  const [page] = useState(1);
  const [rejectDialog, setRejectDialog] = useState<{ open: boolean; id: number | null }>({ open: false, id: null });
  const [rejectReason, setRejectReason] = useState("");

  const { data, isLoading } = usePendingWithdrawals(page);
  const approveMutation = useApproveWithdrawal();
  const rejectMutation = useRejectWithdrawal();
  const { toast } = useToast();

  const orders = data?.orders || [];

  function handleApprove(id: number) {
    approveMutation.mutate(id, {
      onSuccess: () => toast({ title: "已批准", description: `提现订单 #${id} 已批准并提交至支付通道。` }),
      onError: (err) => toast({ variant: "destructive", title: "批准失败", description: err.message }),
    });
  }

  function handleRejectConfirm() {
    if (!rejectDialog.id) return;
    rejectMutation.mutate(
      { id: rejectDialog.id, reason: rejectReason },
      {
        onSuccess: () => {
          toast({ title: "已拒绝", description: `提现订单 #${rejectDialog.id} 已拒绝，冻结金额已解冻。` });
          setRejectDialog({ open: false, id: null });
          setRejectReason("");
        },
        onError: (err) => toast({ variant: "destructive", title: "拒绝失败", description: err.message }),
      }
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-slate-900">提现审核</h1>
        <p className="mt-1 text-sm text-slate-500">审核待处理的提现申请</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">
            待审核列表
            {data?.total ? <span className="ml-2 text-sm font-normal text-slate-500">共 {data.total} 条</span> : null}
          </CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-3">
              <Skeleton className="h-12 w-full" />
              <Skeleton className="h-12 w-full" />
              <Skeleton className="h-12 w-full" />
            </div>
          ) : orders.length === 0 ? (
            <p className="py-12 text-center text-sm text-slate-400">暂无待审核提现申请</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>类型</TableHead>
                  <TableHead>币种</TableHead>
                  <TableHead>金额</TableHead>
                  <TableHead>手续费</TableHead>
                  <TableHead>收款信息</TableHead>
                  <TableHead>时间</TableHead>
                  <TableHead className="text-right">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {orders.map((order) => (
                  <TableRow key={order.id}>
                    <TableCell className="font-mono text-sm">#{order.id}</TableCell>
                    <TableCell>
                      <Badge variant="secondary">
                        {order.type === "CRYPTO" ? "加密" : "法币"}
                      </Badge>
                    </TableCell>
                    <TableCell>{order.currency}</TableCell>
                    <TableCell className="font-medium">{order.amount}</TableCell>
                    <TableCell className="text-slate-500">{order.platform_fee}</TableCell>
                    <TableCell className="max-w-[200px] truncate text-xs font-mono">
                      {order.to_address || "—"}
                    </TableCell>
                    <TableCell className="text-xs text-slate-500">
                      {new Date(order.created_at).toLocaleString("zh-CN")}
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex justify-end gap-2">
                        <Button
                          size="sm"
                          variant="outline"
                          className="text-green-700 hover:bg-green-50"
                          onClick={() => handleApprove(order.id)}
                          disabled={approveMutation.isPending}
                        >
                          {approveMutation.isPending ? (
                            <Loader2 className="mr-1 h-3 w-3 animate-spin" />
                          ) : (
                            <CheckCircle className="mr-1 h-3 w-3" />
                          )}
                          批准
                        </Button>
                        <Button
                          size="sm"
                          variant="outline"
                          className="text-red-700 hover:bg-red-50"
                          onClick={() => setRejectDialog({ open: true, id: order.id })}
                        >
                          <XCircle className="mr-1 h-3 w-3" />
                          拒绝
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      <Dialog open={rejectDialog.open} onOpenChange={(open) => { if (!open) { setRejectDialog({ open: false, id: null }); setRejectReason(""); } }}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>拒绝提现申请</DialogTitle>
          </DialogHeader>
          <div className="py-4">
            <label className="mb-1.5 block text-sm font-medium text-slate-700">拒绝原因</label>
            <Input
              placeholder="请输入拒绝原因（可选）"
              value={rejectReason}
              onChange={(e) => setRejectReason(e.target.value)}
            />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setRejectDialog({ open: false, id: null })}>
              取消
            </Button>
            <Button
              variant="destructive"
              onClick={handleRejectConfirm}
              disabled={rejectMutation.isPending}
            >
              {rejectMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              确认拒绝
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
