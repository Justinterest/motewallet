"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Copy, KeyRound, Plus, ShieldOff } from "lucide-react";
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
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
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
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "@/hooks/use-toast";
import {
  useAdminEmployees,
  useCreateAdminEmployee,
  useResetAdmin2FA,
  useResetAdminPassword,
} from "@/lib/hooks/use-users";
import {
  createEmployeeSchema,
  type CreateEmployeeFormValues,
} from "@/lib/validations/auth";
import type { AdminEmployee } from "@/types/auth";

const ROLE_LABELS: Record<string, string> = {
  SUPER_ADMIN: "超级管理员",
  OPERATOR: "运营",
  FINANCE: "财务",
};

function formatDateTime(value?: string | null) {
  if (!value) return "-";
  return new Date(value).toLocaleString("zh-CN");
}

export default function EmployeesPage() {
  const { data: employees, isLoading, error } = useAdminEmployees();
  const createMutation = useCreateAdminEmployee();
  const resetPasswordMutation = useResetAdminPassword();
  const reset2FAMutation = useResetAdmin2FA();

  const [createOpen, setCreateOpen] = useState(false);
  const [createdPassword, setCreatedPassword] = useState<string | null>(null);
  const [resetPasswordTarget, setResetPasswordTarget] =
    useState<AdminEmployee | null>(null);
  const [resetPasswordResult, setResetPasswordResult] = useState<string | null>(
    null
  );
  const [reset2FATarget, setReset2FATarget] = useState<AdminEmployee | null>(
    null
  );

  const form = useForm<CreateEmployeeFormValues>({
    resolver: zodResolver(createEmployeeSchema),
    defaultValues: {
      username: "",
      email: "",
      role: "OPERATOR",
      password: "",
    },
  });

  function handleCreate(values: CreateEmployeeFormValues) {
    createMutation.mutate(
      {
        username: values.username,
        email: values.email,
        role: values.role,
        password: values.password || undefined,
      },
      {
        onSuccess: (result) => {
          setCreateOpen(false);
          form.reset();
          setCreatedPassword(result.initial_password || null);
          if (!result.initial_password) {
            toast({ title: "员工创建成功" });
          }
        },
        onError: (err) => {
          toast({
            title: "创建失败",
            description: err.message,
            variant: "destructive",
          });
        },
      }
    );
  }

  function handleResetPassword() {
    if (!resetPasswordTarget) return;
    resetPasswordMutation.mutate(resetPasswordTarget.id, {
      onSuccess: (result) => {
        setResetPasswordTarget(null);
        setResetPasswordResult(result.new_password);
      },
      onError: (err) => {
        toast({
          title: "重置密码失败",
          description: err.message,
          variant: "destructive",
        });
      },
    });
  }

  function handleReset2FA() {
    if (!reset2FATarget) return;
    reset2FAMutation.mutate(reset2FATarget.id, {
      onSuccess: () => {
        toast({
          title: "已重置两步验证",
          description: "该员工下次登录需重新绑定验证器",
        });
        setReset2FATarget(null);
      },
      onError: (err) => {
        toast({
          title: "重置两步验证失败",
          description: err.message,
          variant: "destructive",
        });
      },
    });
  }

  async function copyText(text: string) {
    try {
      await navigator.clipboard.writeText(text);
      toast({ title: "已复制到剪贴板" });
    } catch {
      toast({ title: "复制失败", variant: "destructive" });
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">员工管理</h1>
          <p className="mt-1 text-sm text-slate-500">
            添加员工、重置密码与两步验证（仅超级管理员）
          </p>
        </div>
        <Button onClick={() => setCreateOpen(true)}>
          <Plus className="size-4" />
          添加员工
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">员工列表</CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-3">
              {Array.from({ length: 4 }).map((_, i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          ) : error ? (
            <div className="py-12 text-center text-sm text-destructive">
              {error.message || "加载失败"}
            </div>
          ) : !employees || employees.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12 text-slate-400">
              暂无员工
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>用户名</TableHead>
                  <TableHead>邮箱</TableHead>
                  <TableHead>角色</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead>两步验证</TableHead>
                  <TableHead>最近登录</TableHead>
                  <TableHead className="text-right">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {employees.map((employee) => (
                  <TableRow key={employee.id}>
                    <TableCell className="font-medium">
                      {employee.username}
                    </TableCell>
                    <TableCell>{employee.email}</TableCell>
                    <TableCell>
                      {ROLE_LABELS[employee.role] || employee.role}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          employee.status === "ACTIVE" ? "default" : "secondary"
                        }
                      >
                        {employee.status === "ACTIVE" ? "正常" : employee.status}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {employee.totp_enabled ? (
                        <Badge variant="outline" className="text-green-700">
                          已开启
                        </Badge>
                      ) : (
                        <Badge variant="secondary">未开启</Badge>
                      )}
                    </TableCell>
                    <TableCell className="text-slate-500">
                      {formatDateTime(employee.last_login_at)}
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex justify-end gap-2">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => setResetPasswordTarget(employee)}
                        >
                          <KeyRound className="size-3.5" />
                          重置密码
                        </Button>
                        <Button
                          variant="outline"
                          size="sm"
                          disabled={!employee.totp_enabled}
                          onClick={() => setReset2FATarget(employee)}
                        >
                          <ShieldOff className="size-3.5" />
                          重置 2FA
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

      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>添加员工</DialogTitle>
            <DialogDescription>
              创建后员工首次登录必须修改初始密码并绑定两步验证。若不填写密码，将自动生成初始密码。
            </DialogDescription>
          </DialogHeader>
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(handleCreate)}
              className="space-y-4"
            >
              <FormField
                control={form.control}
                name="username"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>用户名</FormLabel>
                    <FormControl>
                      <Input placeholder="登录用户名" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>邮箱</FormLabel>
                    <FormControl>
                      <Input placeholder="name@example.com" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="role"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>角色</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="选择角色" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="OPERATOR">运营</SelectItem>
                        <SelectItem value="FINANCE">财务</SelectItem>
                        <SelectItem value="SUPER_ADMIN">超级管理员</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>初始密码（可选）</FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        placeholder="留空则自动生成"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <DialogFooter>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setCreateOpen(false)}
                >
                  取消
                </Button>
                <Button type="submit" disabled={createMutation.isPending}>
                  {createMutation.isPending ? "创建中..." : "创建"}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      <Dialog
        open={!!createdPassword}
        onOpenChange={(open) => {
          if (!open) setCreatedPassword(null);
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>员工创建成功</DialogTitle>
            <DialogDescription>
              请妥善保存初始密码并告知员工。员工首次登录必须立即修改密码并绑定两步验证，关闭后将无法再次查看。
            </DialogDescription>
          </DialogHeader>
          <div className="flex items-center gap-2 rounded-md border bg-slate-50 px-3 py-2">
            <code className="flex-1 break-all text-sm">{createdPassword}</code>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={() => createdPassword && copyText(createdPassword)}
            >
              <Copy className="size-4" />
            </Button>
          </div>
          <DialogFooter>
            <Button onClick={() => setCreatedPassword(null)}>我已保存</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <AlertDialog
        open={!!resetPasswordTarget}
        onOpenChange={(open) => {
          if (!open) setResetPasswordTarget(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>确认重置密码</AlertDialogTitle>
            <AlertDialogDescription>
              确定要重置员工 {resetPasswordTarget?.username}{" "}
              的密码吗？将生成新的随机密码。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleResetPassword}
              disabled={resetPasswordMutation.isPending}
            >
              {resetPasswordMutation.isPending ? "处理中..." : "确认重置"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <Dialog
        open={!!resetPasswordResult}
        onOpenChange={(open) => {
          if (!open) setResetPasswordResult(null);
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>密码已重置</DialogTitle>
            <DialogDescription>
              请将新密码告知员工，关闭后将无法再次查看。
            </DialogDescription>
          </DialogHeader>
          <div className="flex items-center gap-2 rounded-md border bg-slate-50 px-3 py-2">
            <code className="flex-1 break-all text-sm">
              {resetPasswordResult}
            </code>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={() =>
                resetPasswordResult && copyText(resetPasswordResult)
              }
            >
              <Copy className="size-4" />
            </Button>
          </div>
          <DialogFooter>
            <Button onClick={() => setResetPasswordResult(null)}>我已保存</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <AlertDialog
        open={!!reset2FATarget}
        onOpenChange={(open) => {
          if (!open) setReset2FATarget(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>确认重置两步验证</AlertDialogTitle>
            <AlertDialogDescription>
              确定要重置员工 {reset2FATarget?.username}{" "}
              的两步验证吗？重置后该员工下次登录必须重新绑定验证器。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleReset2FA}
              disabled={reset2FAMutation.isPending}
            >
              {reset2FAMutation.isPending ? "处理中..." : "确认重置"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
