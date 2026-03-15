package category

import (
	"context"

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

func (r *Repository) Create(ctx context.Context, cat *domain.Category) error {
	return r.db.QueryRowxContext(ctx,
		`INSERT INTO categories (user_id, name, color, icon) VALUES ($1, $2, $3, $4) RETURNING *`,
		cat.UserID, cat.Name, cat.Color, cat.Icon,
	).StructScan(cat)
}

func (r *Repository) List(ctx context.Context, userID uuid.UUID) ([]domain.Category, error) {
	cats := make([]domain.Category, 0)
	err := r.db.SelectContext(ctx, &cats,
		`SELECT * FROM categories WHERE user_id = $1 ORDER BY name`, userID)
	return cats, err
}

func (r *Repository) Get(ctx context.Context, id, userID uuid.UUID) (*domain.Category, error) {
	cat := &domain.Category{}
	err := r.db.QueryRowxContext(ctx,
		`SELECT * FROM categories WHERE id = $1 AND user_id = $2`, id, userID,
	).StructScan(cat)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return cat, nil
}

func (r *Repository) Update(ctx context.Context, cat *domain.Category) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE categories SET name=$1, color=$2, icon=$3, updated_at=NOW() WHERE id=$4 AND user_id=$5`,
		cat.Name, cat.Color, cat.Icon, cat.ID, cat.UserID,
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
		`DELETE FROM categories WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
