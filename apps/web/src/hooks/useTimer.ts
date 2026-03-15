"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useState, useEffect } from "react";
import { api } from "@/lib/api";
import { computeElapsedSec } from "@timetracking/core";

export function useTimer() {
  const queryClient = useQueryClient();
  const [elapsed, setElapsed] = useState(0);

  const { data: activeTimer } = useQuery({
    queryKey: ["timer", "active"],
    queryFn: () => api.getActiveTimer(),
    refetchInterval: 30000,
  });

  useEffect(() => {
    if (!activeTimer) return;
    setElapsed(computeElapsedSec(activeTimer));
    if (activeTimer.state !== "running") return;

    const interval = setInterval(() => {
      setElapsed(computeElapsedSec(activeTimer));
    }, 1000);
    return () => clearInterval(interval);
  }, [activeTimer]);

  const invalidate = () => {
    queryClient.invalidateQueries({ queryKey: ["timer"] });
  };

  const start = useMutation({
    mutationFn: ({
      categoryId,
      description,
    }: {
      categoryId: string;
      description?: string;
    }) => api.startTimer(categoryId, description),
    onSuccess: invalidate,
  });

  const pause = useMutation({
    mutationFn: () => api.pauseTimer(),
    onSuccess: invalidate,
  });

  const resume = useMutation({
    mutationFn: () => api.resumeTimer(),
    onSuccess: invalidate,
  });

  const stop = useMutation({
    mutationFn: () => api.stopTimer(),
    onSuccess: invalidate,
  });

  return { activeTimer, elapsed, start, pause, resume, stop };
}
