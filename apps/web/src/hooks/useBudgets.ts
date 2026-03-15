"use client";

import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";

export function useBudgets() {
  return useQuery({
    queryKey: ["budgets", "status"],
    queryFn: () => api.getBudgetStatus(),
    refetchInterval: 60000,
  });
}
