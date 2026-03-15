import { ApiClient } from "@timetracking/core";
import AsyncStorage from "@react-native-async-storage/async-storage";

let accessToken: string | null = null;

export async function loadAccessToken() {
  accessToken = await AsyncStorage.getItem("access_token");
}

export function setAccessToken(token: string | null) {
  accessToken = token;
  if (token) {
    AsyncStorage.setItem("access_token", token);
  } else {
    AsyncStorage.removeItem("access_token");
  }
}

export function getAccessToken(): string | null {
  return accessToken;
}

const BASE_URL =
  process.env.EXPO_PUBLIC_API_URL ?? "http://localhost:8080";

export const api = new ApiClient({
  baseUrl: BASE_URL,
  getToken: getAccessToken,
  onUnauthorized: () => {
    setAccessToken(null);
    AsyncStorage.removeItem("refresh_token");
  },
});
