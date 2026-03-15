"use client";

import { useBudgets } from "@/hooks/useBudgets";
import { getBudgetAlert, formatBudgetSeconds } from "@timetracking/core";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

export function BudgetStatusList() {
  const { data: statuses, isLoading } = useBudgets();

  if (isLoading)
    return <div className="text-sm text-gray-500">Loading budgets...</div>;
  if (!statuses?.length)
    return <div className="text-sm text-gray-500">No budget rules set.</div>;

  return (
    <div className="space-y-3">
      {statuses.map((status) => {
        const alert = getBudgetAlert(status);
        return (
          <Card key={status.rule.id}>
            <CardContent className="py-4">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-medium text-gray-900 capitalize">
                  {status.rule.period_type} budget
                </span>
                <Badge
                  variant={
                    alert === "exceeded"
                      ? "danger"
                      : alert === "warning"
                      ? "warning"
                      : "success"
                  }
                >
                  {Math.round(status.percent_used)}%
                </Badge>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className={`h-2 rounded-full transition-all ${
                    alert === "exceeded"
                      ? "bg-red-500"
                      : alert === "warning"
                      ? "bg-yellow-500"
                      : "bg-green-500"
                  }`}
                  style={{ width: `${Math.min(100, status.percent_used)}%` }}
                />
              </div>
              <div className="mt-1 flex justify-between text-xs text-gray-500">
                <span>{formatBudgetSeconds(status.used_sec)} used</span>
                <span>{formatBudgetSeconds(status.rule.budget_sec)} budget</span>
              </div>
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}
