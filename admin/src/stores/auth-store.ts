import { create } from "zustand";
import type { AdminUser } from "@/types/auth";

interface AuthState {
  admin: AdminUser | null;
  isAuthenticated: boolean;
  setAdmin: (admin: AdminUser) => void;
  clearAdmin: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  admin: null,
  isAuthenticated: false,
  setAdmin: (admin) => set({ admin, isAuthenticated: true }),
  clearAdmin: () => set({ admin: null, isAuthenticated: false }),
}));
