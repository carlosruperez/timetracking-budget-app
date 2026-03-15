package timer

import (
	"context"
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

func (r *Repository) Create(ctx context.Context, entry *domain.TimerEntry) error {
	return r.db.QueryRowxContext(ctx,
		`INSERT INTO timer_entries (user_id, category_id, description, started_at, state)
         VALUES ($1, $2, $3, $4, 'running')
         RETURNING *`,
		entry.UserID, entry.CategoryID, entry.Description, entry.StartedAt,
	).StructScan(entry)
}

func (r *Repository) GetActive(ctx context.Context, userID uuid.UUID) (*domain.TimerEntry, error) {
	entry := &domain.TimerEntry{}
	err := r.db.QueryRowxContext(ctx,
		`SELECT * FROM timer_entries WHERE user_id = $1 AND state IN ('running', 'paused')`,
		userID,
	).StructScan(entry)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return entry, nil
}

func (r *Repository) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.TimerEntry, error) {
	entry := &domain.TimerEntry{}
	err := r.db.QueryRowxContext(ctx,
		`SELECT * FROM timer_entries WHERE id = $1 AND user_id = $2`, id, userID,
	).StructScan(entry)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return entry, nil
}

func (r *Repository) Update(ctx context.Context, entry *domain.TimerEntry) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE timer_entries SET
            category_id=$1, description=$2, started_at=$3, ended_at=$4,
            duration_sec=$5, state=$6, paused_at=$7, updated_at=NOW()
         WHERE id=$8 AND user_id=$9`,
		entry.CategoryID, entry.Description, entry.StartedAt, entry.EndedAt,
		entry.DurationSec, entry.State, entry.PausedAt, entry.ID, entry.UserID,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM timer_entries WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

type ListFilter struct {
	From       *time.Time
	To         *time.Time
	CategoryID *uuid.UUID
	Page       int
	Limit      int
}

func (r *Repository) List(ctx context.Context, userID uuid.UUID, f ListFilter) ([]domain.TimerEntry, int, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	offset := (f.Page - 1) * f.Limit

	query := `SELECT * FROM timer_entries WHERE user_id = :user_id`
	countQuery := `SELECT COUNT(*) FROM timer_entries WHERE user_id = :user_id`

	args := map[string]interface{}{"user_id": userID}

	if f.From != nil {
		query += ` AND started_at >= :from`
		countQuery += ` AND started_at >= :from`
		args["from"] = f.From
	}
	if f.To != nil {
		query += ` AND started_at <= :to`
		countQuery += ` AND started_at <= :to`
		args["to"] = f.To
	}
	if f.CategoryID != nil {
		query += ` AND category_id = :category_id`
		countQuery += ` AND category_id = :category_id`
		args["category_id"] = f.CategoryID
	}

	var total int
	countRows, err := r.db.NamedQueryContext(ctx, countQuery, args)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()
	if countRows.Next() {
		_ = countRows.Scan(&total)
	}

	query += ` ORDER BY started_at DESC LIMIT :limit OFFSET :offset`
	args["limit"] = f.Limit
	args["offset"] = offset

	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	entries := make([]domain.TimerEntry, 0)
	for rows.Next() {
		var e domain.TimerEntry
		if err := rows.StructScan(&e); err != nil {
			return nil, 0, err
		}
		entries = append(entries, e)
	}
	return entries, total, nil
}
