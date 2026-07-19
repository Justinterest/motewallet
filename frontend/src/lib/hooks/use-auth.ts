import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { authApi } from "@/lib/api/auth";
import { useAuthStore } from "@/stores/auth-store";
import type {
  AuthChallenge,
  LoginRequest,
  RegisterRequest,
  User,
} from "@/types/auth";

function getRedirectPath(_user: User): string {
  return "/dashboard";
}

function completeAuth(
  challenge: AuthChallenge,
  setUser: (user: User) => void,
  queryClient: ReturnType<typeof useQueryClient>,
  router: ReturnType<typeof useRouter>
) {
  if (challenge.status !== "SUCCESS" || !challenge.merchant) {
    return challenge;
  }
  setUser(challenge.merchant);
  queryClient.setQueryData(["auth", "me"], challenge.merchant);
  router.push(getRedirectPath(challenge.merchant));
  return challenge;
}

export function useCurrentUser() {
  const { setUser } = useAuthStore();

  return useQuery({
    queryKey: ["auth", "me"],
    queryFn: async () => {
      const user = await authApi.getMe();
      setUser(user);
      return user;
    },
    retry: false,
    staleTime: 5 * 60 * 1000,
  });
}

export function useLogin() {
  return useMutation({
    mutationFn: (data: LoginRequest) => authApi.login(data),
  });
}

export function useVerify2FA() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { setUser } = useAuthStore();

  return useMutation({
    mutationFn: (data: { temp_token: string; code: string }) =>
      authApi.verify2FA(data),
    onSuccess: (challenge) => {
      completeAuth(challenge, setUser, queryClient, router);
    },
  });
}

export function useConfirm2FASetup() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { setUser } = useAuthStore();

  return useMutation({
    mutationFn: (data: { temp_token: string; code: string }) =>
      authApi.confirm2FASetup(data),
    onSuccess: (challenge) => {
      completeAuth(challenge, setUser, queryClient, router);
    },
  });
}

export function useSendVerificationCode() {
  return useMutation({
    mutationFn: (email: string) => authApi.sendVerificationCode({ email }),
  });
}

export function useRegister() {
  return useMutation({
    mutationFn: (data: RegisterRequest) => authApi.register(data),
  });
}

export function useTotpStatus() {
  return useQuery({
    queryKey: ["auth", "2fa", "status"],
    queryFn: () => authApi.getTotpStatus(),
  });
}

export function usePrepareTotpRebind() {
  return useMutation({
    mutationFn: (currentCode: string) => authApi.prepareTotpRebind(currentCode),
  });
}

export function useConfirmTotpRebind() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (code: string) => authApi.confirmTotpRebind(code),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["auth", "2fa", "status"] });
      queryClient.invalidateQueries({ queryKey: ["auth", "me"] });
    },
  });
}

export function useLogout() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { clearUser } = useAuthStore();

  return useMutation({
    mutationFn: () => authApi.logout(),
    onSuccess: () => {
      clearUser();
      queryClient.clear();
      router.push("/login");
    },
  });
}
