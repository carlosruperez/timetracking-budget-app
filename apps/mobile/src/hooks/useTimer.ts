import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useState, useEffect, useRef } from "react";
import { AppState } from "react-native";
import { api } from "@/lib/api";
import { computeElapsedSec } from "@timetracking/core";

export function useTimer() {
  const queryClient = useQueryClient();
  const [elapsed, setElapsed] = useState(0);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const { data: activeTimer, refetch } = useQuery({
    queryKey: ["timer", "active"],
    queryFn: () => api.getActiveTimer(),
    refetchInterval: 30000,
  });

  // Recompute elapsed when app comes to foreground
  useEffect(() => {
    const sub = AppState.addEventListener("change", (state) => {
      if (state === "active") {
        refetch();
      }
    });
    return () => sub.remove();
  }, [refetch]);

  useEffect(() => {
    if (intervalRef.current) clearInterval(intervalRef.current);
    if (!activeTimer) {
      setElapsed(0);
      return;
    }
    setElapsed(computeElapsedSec(activeTimer));
    if (activeTimer.state !== "running") return;

    intervalRef.current = setInterval(() => {
      setElapsed(computeElapsedSec(activeTimer));
    }, 1000);

    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
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
