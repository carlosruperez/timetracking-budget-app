import { BudgetStatus } from "./types";

export function getBudgetAlert(
  status: BudgetStatus
): "ok" | "warning" | "exceeded" {
  if (status.percent_used >= 100) return "exceeded";
  if (status.percent_used >= 80) return "warning";
  return "ok";
}

export function formatBudgetSeconds(sec: number): string {
  const h = Math.floor(sec / 3600);
  const m = Math.floor((sec % 3600) / 60);
  if (h === 0) return `${m}m`;
  if (m === 0) return `${h}h`;
  return `${h}h ${m}m`;
}
