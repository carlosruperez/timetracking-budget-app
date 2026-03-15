package domain

import (
	"time"

	"github.com/google/uuid"
)

type TimerState string

const (
	TimerStateRunning TimerState = "running"
	TimerStatePaused  TimerState = "paused"
	TimerStateStopped TimerState = "stopped"
)

type TimerEntry struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	UserID      uuid.UUID  `db:"user_id" json:"user_id"`
	CategoryID  uuid.UUID  `db:"category_id" json:"category_id"`
	Description string     `db:"description" json:"description"`
	StartedAt   time.Time  `db:"started_at" json:"started_at"`
	EndedAt     *time.Time `db:"ended_at" json:"ended_at"`
	DurationSec int64      `db:"duration_sec" json:"duration_sec"`
	State       TimerState `db:"state" json:"state"`
	PausedAt    *time.Time `db:"paused_at" json:"paused_at"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// ElapsedSec computes current elapsed seconds including live time if running
func (t *TimerEntry) ElapsedSec() int64 {
	if t.State == TimerStateRunning {
		return t.DurationSec + int64(time.Since(t.StartedAt).Seconds())
	}
	return t.DurationSec
}
