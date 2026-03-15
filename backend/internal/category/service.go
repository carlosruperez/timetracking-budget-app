package category

import (
	"context"

	"github.com/google/uuid"
	"github.com/rupi/timetracking/internal/domain"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

type UpdateRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, req CreateRequest) (*domain.Category, error) {
	if req.Name == "" {
		return nil, domain.ErrInvalidInput
	}
	if req.Color == "" {
		req.Color = "#6366f1"
	}
	if req.Icon == "" {
		req.Icon = "clock"
	}
	cat := &domain.Category{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
		Icon:   req.Icon,
	}
	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]domain.Category, error) {
	return s.repo.List(ctx, userID)
}

func (s *Service) Get(ctx context.Context, userID, catID uuid.UUID) (*domain.Category, error) {
	return s.repo.Get(ctx, catID, userID)
}

func (s *Service) Update(ctx context.Context, userID, catID uuid.UUID, req UpdateRequest) (*domain.Category, error) {
	cat, err := s.repo.Get(ctx, catID, userID)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		cat.Name = req.Name
	}
	if req.Color != "" {
		cat.Color = req.Color
	}
	if req.Icon != "" {
		cat.Icon = req.Icon
	}
	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *Service) Delete(ctx context.Context, userID, catID uuid.UUID) error {
	return s.repo.Delete(ctx, catID, userID)
}
