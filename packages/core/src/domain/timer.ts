import { TimerEntry, TimerState } from "./types";

export function computeElapsedSec(entry: TimerEntry): number {
  if (entry.state === "running") {
    const startedAt = new Date(entry.started_at).getTime();
    const now = Date.now();
    return entry.duration_sec + Math.floor((now - startedAt) / 1000);
  }
  return entry.duration_sec;
}

export function canTransition(
  from: TimerState,
  action: "pause" | "resume" | "stop"
): boolean {
  const transitions: Record<TimerState, string[]> = {
    running: ["pause", "stop"],
    paused: ["resume", "stop"],
    stopped: [],
  };
  return transitions[from]?.includes(action) ?? false;
}

export function formatDuration(seconds: number): string {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = seconds % 60;
  return [h, m, s].map((v) => String(v).padStart(2, "0")).join(":");
}
