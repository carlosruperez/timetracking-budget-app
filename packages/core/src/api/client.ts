import type {
  User,
  Category,
  TimerEntry,
  BudgetRule,
  BudgetStatus,
  TokenPair,
  PaginatedResponse,
} from "../domain/types";

export interface ApiClientConfig {
  baseUrl: string;
  getToken: () => string | null;
  onUnauthorized?: () => void;
}

export class ApiClient {
  private baseUrl: string;
  private getToken: () => string | null;
  private onUnauthorized?: () => void;

  constructor(config: ApiClientConfig) {
    this.baseUrl = config.baseUrl.replace(/\/$/, "");
    this.getToken = config.getToken;
    this.onUnauthorized = config.onUnauthorized;
  }

  private async request<T>(path: string, init?: RequestInit): Promise<T> {
    const token = this.getToken();
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      ...(init?.headers as Record<string, string>),
    };
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const res = await fetch(`${this.baseUrl}${path}`, { ...init, headers });

    if (res.status === 401) {
      this.onUnauthorized?.();
      throw new Error("Unauthorized");
    }
    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: "Request failed" }));
      throw new Error((err as { error?: string }).error ?? "Request failed");
    }
    if (res.status === 204) return undefined as unknown as T;
    return res.json() as Promise<T>;
  }

  // Auth
  async register(
    email: string,
    password: string,
    name: string,
    timezone?: string
  ) {
    return this.request<{ user: User; tokens: TokenPair }>(
      "/api/v1/auth/register",
      {
        method: "POST",
        body: JSON.stringify({
          email,
          password,
          name,
          timezone: timezone ?? "UTC",
        }),
      }
    );
  }

  async login(email: string, password: string) {
    return this.request<{ user: User; tokens: TokenPair }>(
      "/api/v1/auth/login",
      {
        method: "POST",
        body: JSON.stringify({ email, password }),
      }
    );
  }

  async refresh(refreshToken: string) {
    return this.request<TokenPair>("/api/v1/auth/refresh", {
      method: "POST",
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
  }

  async logout(refreshToken: string) {
    return this.request<void>("/api/v1/auth/logout", {
      method: "POST",
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
  }

  async getMe() {
    return this.request<User>("/api/v1/auth/me");
  }

  // Categories
  async listCategories() {
    return this.request<Category[]>("/api/v1/categories");
  }

  async createCategory(name: string, color?: string, icon?: string) {
    return this.request<Category>("/api/v1/categories", {
      method: "POST",
      body: JSON.stringify({ name, color, icon }),
    });
  }

  async updateCategory(
    id: string,
    data: Partial<Pick<Category, "name" | "color" | "icon">>
  ) {
    return this.request<Category>(`/api/v1/categories/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async deleteCategory(id: string) {
    return this.request<void>(`/api/v1/categories/${id}`, {
      method: "DELETE",
    });
  }

  // Timer
  async getActiveTimer() {
    return this.request<TimerEntry | null>("/api/v1/timer/active");
  }

  async startTimer(categoryId: string, description?: string) {
    return this.request<TimerEntry>("/api/v1/timer/start", {
      method: "POST",
      body: JSON.stringify({
        category_id: categoryId,
        description: description ?? "",
      }),
    });
  }

  async pauseTimer() {
    return this.request<TimerEntry>("/api/v1/timer/pause", { method: "POST" });
  }

  async resumeTimer() {
    return this.request<TimerEntry>("/api/v1/timer/resume", {
      method: "POST",
    });
  }

  async stopTimer() {
    return this.request<TimerEntry>("/api/v1/timer/stop", { method: "POST" });
  }

  async listTimerEntries(params?: {
    from?: string;
    to?: string;
    category_id?: string;
    page?: number;
    limit?: number;
  }) {
    const q = new URLSearchParams();
    if (params?.from) q.set("from", params.from);
    if (params?.to) q.set("to", params.to);
    if (params?.category_id) q.set("category_id", params.category_id);
    if (params?.page) q.set("page", String(params.page));
    if (params?.limit) q.set("limit", String(params.limit));
    const qs = q.toString() ? `?${q}` : "";
    return this.request<PaginatedResponse<TimerEntry>>(
      `/api/v1/timer/entries${qs}`
    );
  }

  // Budgets
  async listBudgets() {
    return this.request<BudgetRule[]>("/api/v1/budgets");
  }

  async createBudget(data: {
    category_id: string;
    period_type: string;
    budget_sec: number;
  }) {
    return this.request<BudgetRule>("/api/v1/budgets", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async getBudgetStatus(timezone?: string) {
    const tz =
      timezone ?? Intl.DateTimeFormat().resolvedOptions().timeZone;
    return this.request<BudgetStatus[]>(
      `/api/v1/budgets/status?timezone=${encodeURIComponent(tz)}`
    );
  }

  // Reports
  async getSummary(from: string, to: string) {
    return this.request<
      Array<{
        category_id: string;
        category_name: string;
        total_sec: number;
        entry_count: number;
      }>
    >(`/api/v1/reports/summary?from=${from}&to=${to}`);
  }

  async getDaily(from: string, to: string, timezone?: string) {
    const tz =
      timezone ?? Intl.DateTimeFormat().resolvedOptions().timeZone;
    return this.request<Array<{ date: string; total_sec: number }>>(
      `/api/v1/reports/daily?from=${from}&to=${to}&timezone=${encodeURIComponent(tz)}`
    );
  }

  // SSE stream
  createTimerStream(
    onTick: (data: {
      entry_id: string;
      state: string;
      elapsed_sec: number;
      category_id: string;
    }) => void
  ): EventSource {
    const url = `${this.baseUrl}/api/v1/timer/stream`;
    const es = new EventSource(url);
    es.addEventListener("tick", ((e: MessageEvent) => {
      try {
        onTick(JSON.parse(e.data));
      } catch {
        // ignore parse errors
      }
    }) as EventListener);
    return es;
  }
}
