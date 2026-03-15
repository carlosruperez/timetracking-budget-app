"use client";

import { useTimer } from "@/hooks/useTimer";
import { Button } from "@/components/ui/button";
import type { TimerEntry } from "@timetracking/core";

interface Props {
  timer: TimerEntry;
}

export function TimerControls({ timer }: Props) {
  const { pause, resume, stop } = useTimer();

  return (
    <div className="flex gap-2">
      {timer.state === "running" && (
        <Button
          variant="secondary"
          onClick={() => pause.mutate()}
          disabled={pause.isPending}
        >
          Pause
        </Button>
      )}
      {timer.state === "paused" && (
        <Button
          variant="primary"
          onClick={() => resume.mutate()}
          disabled={resume.isPending}
        >
          Resume
        </Button>
      )}
      <Button
        variant="danger"
        onClick={() => stop.mutate()}
        disabled={stop.isPending}
      >
        Stop
      </Button>
    </div>
  );
}
