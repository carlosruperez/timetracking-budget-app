package timer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rupi/timetracking/internal/domain"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type StartRequest struct {
	CategoryID  uuid.UUID `json:"category_id"`
	Description string    `json:"description"`
}

func (s *Service) Start(ctx context.Context, userID uuid.UUID, req StartRequest) (*domain.TimerEntry, error) {
	if req.CategoryID == uuid.Nil {
		return nil, domain.ErrInvalidInput
	}
	entry := &domain.TimerEntry{
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		StartedAt:   time.Now().UTC(),
		State:       domain.TimerStateRunning,
	}
	if err := s.repo.Create(ctx, entry); err != nil {
		// DB unique constraint violation -> active timer exists
		return nil, domain.ErrActiveTimerExists
	}
	return entry, nil
}

func (s *Service) GetActive(ctx context.Context, userID uuid.UUID) (*domain.TimerEntry, error) {
	return s.repo.GetActive(ctx, userID)
}

func (s *Service) Pause(ctx context.Context, userID uuid.UUID) (*domain.TimerEntry, error) {
	entry, err := s.repo.GetActive(ctx, userID)
	if err != nil {
		return nil, err
	}
	_, err = ValidateTransition(entry.State, "pause")
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	// Accumulate elapsed into duration_sec
	entry.DurationSec += int64(now.Sub(entry.StartedAt).Seconds())
	entry.StartedAt = now // reset -- will be used on resume
	entry.State = domain.TimerStatePaused
	entry.PausedAt = &now

	if err := s.repo.Update(ctx, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *Service) Resume(ctx context.Context, userID uuid.UUID) (*domain.TimerEntry, error) {
	entry, err := s.repo.GetActive(ctx, userID)
	if err != nil {
		return nil, err
	}
	_, err = ValidateTransition(entry.State, "resume")
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	entry.StartedAt = now
	entry.State = domain.TimerStateRunning
	entry.PausedAt = nil

	if err := s.repo.Update(ctx, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *Service) Stop(ctx context.Context, userID uuid.UUID) (*domain.TimerEntry, error) {
	entry, err := s.repo.GetActive(ctx, userID)
	if err != nil {
		return nil, err
	}
	_, err = ValidateTransition(entry.State, "stop")
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	if entry.State == domain.TimerStateRunning {
		entry.DurationSec += int64(now.Sub(entry.StartedAt).Seconds())
	}
	entry.State = domain.TimerStateStopped
	entry.EndedAt = &now
	entry.PausedAt = nil

	if err := s.repo.Update(ctx, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *Service) GetEntry(ctx context.Context, userID, entryID uuid.UUID) (*domain.TimerEntry, error) {
	return s.repo.GetByID(ctx, entryID, userID)
}

func (s *Service) UpdateEntry(ctx context.Context, userID, entryID uuid.UUID, description string, categoryID uuid.UUID) (*domain.TimerEntry, error) {
	entry, err := s.repo.GetByID(ctx, entryID, userID)
	if err != nil {
		return nil, err
	}
	if description != "" {
		entry.Description = description
	}
	if categoryID != uuid.Nil {
		entry.CategoryID = categoryID
	}
	if err := s.repo.Update(ctx, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *Service) DeleteEntry(ctx context.Context, userID, entryID uuid.UUID) error {
	return s.repo.Delete(ctx, entryID, userID)
}

func (s *Service) ListEntries(ctx context.Context, userID uuid.UUID, f ListFilter) ([]domain.TimerEntry, int, error) {
	return s.repo.List(ctx, userID, f)
}

func (s *Service) GetElapsed(ctx context.Context, userID uuid.UUID) (int64, *domain.TimerEntry, error) {
	entry, err := s.repo.GetActive(ctx, userID)
	if err != nil {
		return 0, nil, err
	}
	return entry.ElapsedSec(), entry, nil
}
