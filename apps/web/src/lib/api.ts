import { ApiClient } from "@timetracking/core";

let accessToken: string | null = null;

export function setAccessToken(token: string | null) {
  accessToken = token;
  if (typeof window !== "undefined") {
    if (token) {
      localStorage.setItem("access_token", token);
    } else {
      localStorage.removeItem("access_token");
    }
  }
}

export function getAccessToken(): string | null {
  if (accessToken) return accessToken;
  if (typeof window !== "undefined") {
    accessToken = localStorage.getItem("access_token");
  }
  return accessToken;
}

export const api = new ApiClient({
  baseUrl: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080",
  getToken: getAccessToken,
  onUnauthorized: () => {
    setAccessToken(null);
    if (typeof window !== "undefined") {
      localStorage.removeItem("refresh_token");
      window.location.href = "/login";
    }
  },
});
