import { create } from "zustand";
import AsyncStorage from "@react-native-async-storage/async-storage";
import type { User } from "@timetracking/core";

interface AuthState {
  user: User | null;
  refreshToken: string | null;
  isLoading: boolean;
  setAuth: (user: User, accessToken: string, refreshToken: string) => void;
  clearAuth: () => void;
  setLoading: (v: boolean) => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  refreshToken: null,
  isLoading: true,
  setAuth: (user, accessToken, refreshToken) => {
    AsyncStorage.setItem("refresh_token", refreshToken);
    set({ user, refreshToken });
  },
  clearAuth: () => {
    AsyncStorage.removeItem("refresh_token");
    AsyncStorage.removeItem("access_token");
    set({ user: null, refreshToken: null });
  },
  setLoading: (v) => set({ isLoading: v }),
}));
