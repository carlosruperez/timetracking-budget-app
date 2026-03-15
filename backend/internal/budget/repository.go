package budget

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rupi/timetracking/internal/domain"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, rule *domain.BudgetRule) error {
	return r.db.QueryRowxContext(ctx,
		`INSERT INTO budget_rules (user_id, category_id, period_type, budget_sec, active)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING *`,
		rule.UserID, rule.CategoryID, rule.PeriodType, rule.BudgetSec, rule.Active,
	).StructScan(rule)
}

func (r *Repository) List(ctx context.Context, userID uuid.UUID) ([]domain.BudgetRule, error) {
	rules := make([]domain.BudgetRule, 0)
	err := r.db.SelectContext(ctx, &rules,
		`SELECT * FROM budget_rules WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	return rules, err
}

func (r *Repository) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.BudgetRule, error) {
	rule := &domain.BudgetRule{}
	err := r.db.QueryRowxContext(ctx,
		`SELECT * FROM budget_rules WHERE id = $1 AND user_id = $2`, id, userID,
	).StructScan(rule)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return rule, nil
}

func (r *Repository) Update(ctx context.Context, rule *domain.BudgetRule) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE budget_rules SET period_type=$1, budget_sec=$2, active=$3, updated_at=NOW()
         WHERE id=$4 AND user_id=$5`,
		rule.PeriodType, rule.BudgetSec, rule.Active, rule.ID, rule.UserID,
	)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM budget_rules WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// GetUsedSecondsForPeriod returns total duration_sec for a category in the given time window.
// For running timers, it adds the live elapsed time.
func (r *Repository) GetUsedSecondsForPeriod(ctx context.Context, userID, categoryID uuid.UUID, from, to time.Time) (int64, error) {
	var used int64
	err := r.db.QueryRowContext(ctx, fmt.Sprintf(`
        SELECT COALESCE(SUM(
            CASE
                WHEN state = '%s' THEN duration_sec + GREATEST(0, EXTRACT(EPOCH FROM (NOW() - started_at))::BIGINT)
                ELSE duration_sec
            END
        ), 0)
        FROM timer_entries
        WHERE user_id = $1
          AND category_id = $2
          AND started_at >= $3
          AND started_at < $4
    `, domain.TimerStateRunning), userID, categoryID, from, to).Scan(&used)
	return used, err
}
