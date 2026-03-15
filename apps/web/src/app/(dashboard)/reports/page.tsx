"use client";

import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { formatDuration } from "@timetracking/core";
import { Card, CardHeader, CardContent } from "@/components/ui/card";

function getLastWeekRange() {
  const to = new Date();
  const from = new Date();
  from.setDate(from.getDate() - 7);
  return { from: from.toISOString(), to: to.toISOString() };
}

export default function ReportsPage() {
  const { from, to } = getLastWeekRange();

  const { data: summary = [], isLoading } = useQuery({
    queryKey: ["reports", "summary", from, to],
    queryFn: () => api.getSummary(from, to),
  });

  const totalSec = summary.reduce((acc, s) => acc + s.total_sec, 0);

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">Reports</h1>
      <p className="text-sm text-gray-500">Last 7 days</p>

      <Card>
        <CardHeader>
          <h2 className="font-semibold text-gray-900">Total time tracked</h2>
          <div className="text-3xl font-mono font-bold text-indigo-600 mt-1">
            {formatDuration(totalSec)}
          </div>
        </CardHeader>
      </Card>

      <Card>
        <CardHeader>
          <h2 className="font-semibold text-gray-900">By category</h2>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="text-gray-500 text-sm">Loading...</div>
          ) : summary.length === 0 ? (
            <div className="text-gray-500 text-sm">
              No data for this period.
            </div>
          ) : (
            <div className="space-y-3">
              {summary.map((item) => (
                <div
                  key={item.category_id}
                  className="flex items-center justify-between"
                >
                  <div>
                    <div className="font-medium text-gray-900">
                      {item.category_name}
                    </div>
                    <div className="text-xs text-gray-500">
                      {item.entry_count} sessions
                    </div>
                  </div>
                  <div className="font-mono text-sm text-gray-700">
                    {formatDuration(item.total_sec)}
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
