"use client";

import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/auth";
import { api, setAccessToken } from "@/lib/api";

export function useAuth() {
  const router = useRouter();
  const { user, setAuth, clearAuth } = useAuthStore();

  const login = async (email: string, password: string) => {
    const { user, tokens } = await api.login(email, password);
    setAuth(user, tokens.access_token, tokens.refresh_token);
    setAccessToken(tokens.access_token);
    localStorage.setItem("refresh_token", tokens.refresh_token);
    router.push("/");
  };

  const register = async (
    email: string,
    password: string,
    name: string
  ) => {
    const { user, tokens } = await api.register(email, password, name);
    setAuth(user, tokens.access_token, tokens.refresh_token);
    setAccessToken(tokens.access_token);
    localStorage.setItem("refresh_token", tokens.refresh_token);
    router.push("/");
  };

  const logout = async () => {
    const refreshToken = localStorage.getItem("refresh_token");
    if (refreshToken) {
      await api.logout(refreshToken).catch(() => {});
    }
    clearAuth();
    setAccessToken(null);
    router.push("/login");
  };

  return { user, login, register, logout, isAuthenticated: !!user };
}
