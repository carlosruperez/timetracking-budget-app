package domain

import (
	"time"

	"github.com/google/uuid"
)

type PeriodType string

const (
	PeriodDaily   PeriodType = "daily"
	PeriodWeekly  PeriodType = "weekly"
	PeriodMonthly PeriodType = "monthly"
)

type BudgetRule struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	UserID     uuid.UUID  `db:"user_id" json:"user_id"`
	CategoryID uuid.UUID  `db:"category_id" json:"category_id"`
	PeriodType PeriodType `db:"period_type" json:"period_type"`
	BudgetSec  int64      `db:"budget_sec" json:"budget_sec"`
	Active     bool       `db:"active" json:"active"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}

type BudgetStatus struct {
	Rule         BudgetRule `json:"rule"`
	UsedSec      int64      `json:"used_sec"`
	RemainingSec int64      `json:"remaining_sec"`
	PercentUsed  float64    `json:"percent_used"`
	PeriodStart  time.Time  `json:"period_start"`
	PeriodEnd    time.Time  `json:"period_end"`
}
