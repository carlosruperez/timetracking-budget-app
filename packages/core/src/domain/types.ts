export type TimerState = "running" | "paused" | "stopped";
export type PeriodType = "daily" | "weekly" | "monthly";

export interface User {
  id: string;
  email: string;
  name: string;
  timezone: string;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: string;
  user_id: string;
  name: string;
  color: string;
  icon: string;
  created_at: string;
  updated_at: string;
}

export interface TimerEntry {
  id: string;
  user_id: string;
  category_id: string;
  description: string;
  started_at: string;
  ended_at: string | null;
  duration_sec: number;
  state: TimerState;
  paused_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface BudgetRule {
  id: string;
  user_id: string;
  category_id: string;
  period_type: PeriodType;
  budget_sec: number;
  active: boolean;
  created_at: string;
  updated_at: string;
}

export interface BudgetStatus {
  rule: BudgetRule;
  used_sec: number;
  remaining_sec: number;
  percent_used: number;
  period_start: string;
  period_end: string;
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
}

export interface PaginatedResponse<T> {
  entries: T[];
  total: number;
  page: number;
  limit: number;
}
