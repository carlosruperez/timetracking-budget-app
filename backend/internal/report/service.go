package report

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

type CategorySummary struct {
	CategoryID   uuid.UUID `db:"category_id" json:"category_id"`
	CategoryName string    `db:"category_name" json:"category_name"`
	TotalSec     int64     `db:"total_sec" json:"total_sec"`
	EntryCount   int       `db:"entry_count" json:"entry_count"`
}

type DailyTotal struct {
	Date     string `db:"day" json:"date"`
	TotalSec int64  `db:"total_sec" json:"total_sec"`
}

func (s *Service) Summary(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]CategorySummary, error) {
	var result []CategorySummary
	err := s.db.SelectContext(ctx, &result, `
        SELECT
            te.category_id,
            c.name AS category_name,
            COALESCE(SUM(te.duration_sec), 0) AS total_sec,
            COUNT(*) AS entry_count
        FROM timer_entries te
        JOIN categories c ON c.id = te.category_id
        WHERE te.user_id = $1
          AND te.started_at >= $2
          AND te.started_at < $3
          AND te.state = 'stopped'
        GROUP BY te.category_id, c.name
        ORDER BY total_sec DESC
    `, userID, from, to)
	return result, err
}

func (s *Service) Daily(ctx context.Context, userID uuid.UUID, from, to time.Time, timezone string) ([]DailyTotal, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	_ = loc // timezone offset used in SQL via AT TIME ZONE

	var result []DailyTotal
	err = s.db.SelectContext(ctx, &result, `
        SELECT
            DATE(started_at AT TIME ZONE $4) AS day,
            COALESCE(SUM(duration_sec), 0) AS total_sec
        FROM timer_entries
        WHERE user_id = $1
          AND started_at >= $2
          AND started_at < $3
          AND state = 'stopped'
        GROUP BY day
        ORDER BY day ASC
    `, userID, from, to, timezone)
	return result, err
}

func (s *Service) Weekly(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]CategorySummary, error) {
	return s.Summary(ctx, userID, from, to)
}
