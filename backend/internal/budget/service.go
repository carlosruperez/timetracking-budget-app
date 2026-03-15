package budget

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

type CreateRequest struct {
	CategoryID uuid.UUID         `json:"category_id"`
	PeriodType domain.PeriodType `json:"period_type"`
	BudgetSec  int64             `json:"budget_sec"`
}

type UpdateRequest struct {
	PeriodType domain.PeriodType `json:"period_type"`
	BudgetSec  int64             `json:"budget_sec"`
	Active     *bool             `json:"active"`
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, req CreateRequest) (*domain.BudgetRule, error) {
	if req.CategoryID == uuid.Nil || req.BudgetSec <= 0 {
		return nil, domain.ErrInvalidInput
	}
	if req.PeriodType == "" {
		req.PeriodType = domain.PeriodDaily
	}
	rule := &domain.BudgetRule{
		UserID:     userID,
		CategoryID: req.CategoryID,
		PeriodType: req.PeriodType,
		BudgetSec:  req.BudgetSec,
		Active:     true,
	}
	if err := s.repo.Create(ctx, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]domain.BudgetRule, error) {
	return s.repo.List(ctx, userID)
}

func (s *Service) Update(ctx context.Context, userID, ruleID uuid.UUID, req UpdateRequest) (*domain.BudgetRule, error) {
	rule, err := s.repo.GetByID(ctx, ruleID, userID)
	if err != nil {
		return nil, err
	}
	if req.PeriodType != "" {
		rule.PeriodType = req.PeriodType
	}
	if req.BudgetSec > 0 {
		rule.BudgetSec = req.BudgetSec
	}
	if req.Active != nil {
		rule.Active = *req.Active
	}
	if err := s.repo.Update(ctx, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *Service) Delete(ctx context.Context, userID, ruleID uuid.UUID) error {
	return s.repo.Delete(ctx, ruleID, userID)
}

func (s *Service) GetStatus(ctx context.Context, userID uuid.UUID, timezone string) ([]domain.BudgetStatus, error) {
	rules, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}

	now := time.Now().In(loc)
	var statuses []domain.BudgetStatus

	for _, rule := range rules {
		if !rule.Active {
			continue
		}
		start, end := periodBounds(now, rule.PeriodType)
		// Convert to UTC for DB query
		startUTC := start.UTC()
		endUTC := end.UTC()

		used, err := s.repo.GetUsedSecondsForPeriod(ctx, userID, rule.CategoryID, startUTC, endUTC)
		if err != nil {
			continue
		}

		remaining := rule.BudgetSec - used
		if remaining < 0 {
			remaining = 0
		}
		var pct float64
		if rule.BudgetSec > 0 {
			pct = float64(used) / float64(rule.BudgetSec) * 100
		}

		statuses = append(statuses, domain.BudgetStatus{
			Rule:         rule,
			UsedSec:      used,
			RemainingSec: remaining,
			PercentUsed:  pct,
			PeriodStart:  start,
			PeriodEnd:    end,
		})
	}

	return statuses, nil
}

// periodBounds returns [start, end) in local time for the current period
func periodBounds(now time.Time, pt domain.PeriodType) (time.Time, time.Time) {
	y, m, d := now.Date()
	loc := now.Location()

	switch pt {
	case domain.PeriodDaily:
		start := time.Date(y, m, d, 0, 0, 0, 0, loc)
		end := start.AddDate(0, 0, 1)
		return start, end
	case domain.PeriodWeekly:
		weekday := int(now.Weekday())
		if weekday == 0 { // Sunday → treat as 7
			weekday = 7
		}
		// Week starts on Monday
		start := time.Date(y, m, d-weekday+1, 0, 0, 0, 0, loc)
		end := start.AddDate(0, 0, 7)
		return start, end
	case domain.PeriodMonthly:
		start := time.Date(y, m, 1, 0, 0, 0, 0, loc)
		end := start.AddDate(0, 1, 0)
		return start, end
	default:
		start := time.Date(y, m, d, 0, 0, 0, 0, loc)
		end := start.AddDate(0, 0, 1)
		return start, end
	}
}
