import { useEffect } from "react";
import { Stack } from "expo-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StatusBar } from "expo-status-bar";
import { loadAccessToken } from "@/lib/api";
import { useAuthStore } from "@/store/auth";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { api, setAccessToken } from "@/lib/api";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { staleTime: 30_000, retry: 1 },
  },
});

export default function RootLayout() {
  const { setAuth, setLoading } = useAuthStore();

  useEffect(() => {
    async function bootstrap() {
      await loadAccessToken();
      // Try to restore session from stored refresh token
      const refreshToken = await AsyncStorage.getItem("refresh_token");
      if (refreshToken) {
        try {
          const tokens = await api.refresh(refreshToken);
          setAccessToken(tokens.access_token);
          const user = await api.getMe();
          setAuth(user, tokens.access_token, tokens.refresh_token);
        } catch {
          // Refresh failed — stay logged out
        }
      }
      setLoading(false);
    }
    bootstrap();
  }, []);

  return (
    <QueryClientProvider client={queryClient}>
      <StatusBar style="auto" />
      <Stack screenOptions={{ headerShown: false }} />
    </QueryClientProvider>
  );
}
