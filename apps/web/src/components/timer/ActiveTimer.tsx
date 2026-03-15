"use client";

import { useTimer } from "@/hooks/useTimer";
import { formatDuration } from "@timetracking/core";
import { TimerControls } from "./TimerControls";
import { Card, CardContent } from "@/components/ui/card";

export function ActiveTimer() {
  const { activeTimer, elapsed } = useTimer();

  if (!activeTimer) {
    return (
      <Card>
        <CardContent className="py-8 text-center text-gray-500">
          No active timer. Start one below.
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardContent className="flex items-center justify-between py-6">
        <div>
          <div className="text-4xl font-mono font-bold text-gray-900">
            {formatDuration(elapsed)}
          </div>
          <div className="mt-1 text-sm text-gray-500">
            {activeTimer.description || "No description"}
          </div>
          <div className="mt-1 flex items-center gap-2">
            <span
              className={`inline-flex h-2 w-2 rounded-full ${
                activeTimer.state === "running"
                  ? "bg-green-500 animate-pulse"
                  : "bg-yellow-500"
              }`}
            />
            <span className="text-xs text-gray-500 capitalize">
              {activeTimer.state}
            </span>
          </div>
        </div>
        <TimerControls timer={activeTimer} />
      </CardContent>
    </Card>
  );
}
