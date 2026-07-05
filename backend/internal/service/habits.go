package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/qper/hf/internal/domain"
)

var (
	ErrHabitValidation = errors.New("invalid habit payload")
	ErrHabitNotFound   = errors.New("habit not found")
	ErrHabitForbidden  = errors.New("habit access denied")
)

type HabitRepository interface {
	CreateHabit(ctx context.Context, userID string, req domain.CreateHabitRequest) (*domain.Habit, error)
	ListHabits(ctx context.Context, userID string, categoryID *string, archived *bool) ([]domain.Habit, error)
	GetHabitByID(ctx context.Context, userID string, habitID string) (*domain.Habit, error)
	UpdateHabit(ctx context.Context, userID string, habitID string, req domain.UpdateHabitRequest) (*domain.Habit, error)
	DeleteHabit(ctx context.Context, userID string, habitID string) error
	ArchiveHabit(ctx context.Context, userID string, habitID string, archived bool) (*domain.Habit, error)
	ReorderHabits(ctx context.Context, userID string, ids []string) ([]domain.Habit, error)
}

type HabitService struct {
	repo HabitRepository
}

func NewHabitService(repo HabitRepository) *HabitService {
	return &HabitService{repo: repo}
}

func (s *HabitService) Create(ctx context.Context, userID string, req domain.CreateHabitRequest) (*domain.Habit, error) {
	if err := validateCreateHabit(req); err != nil {
		return nil, err
	}
	return s.repo.CreateHabit(ctx, userID, req)
}

func (s *HabitService) List(ctx context.Context, userID string, categoryID *string, archived *bool) ([]domain.Habit, error) {
	return s.repo.ListHabits(ctx, userID, categoryID, archived)
}

func (s *HabitService) GetByID(ctx context.Context, userID string, habitID string) (*domain.Habit, error) {
	habit, err := s.repo.GetHabitByID(ctx, userID, habitID)
	if err != nil {
		return nil, err
	}
	if habit == nil {
		return nil, ErrHabitNotFound
	}
	return habit, nil
}

func (s *HabitService) Update(ctx context.Context, userID string, habitID string, req domain.UpdateHabitRequest) (*domain.Habit, error) {
	if req.Type != "" {
		return nil, ErrHabitValidation
	}
	habit, err := s.repo.UpdateHabit(ctx, userID, habitID, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrHabitNotFound
		}
		return nil, err
	}
	if habit == nil {
		return nil, ErrHabitNotFound
	}
	return habit, nil
}

func (s *HabitService) Delete(ctx context.Context, userID string, habitID string) error {
	err := s.repo.DeleteHabit(ctx, userID, habitID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrHabitNotFound
		}
	}
	return err
}

func (s *HabitService) Archive(ctx context.Context, userID string, habitID string, archived bool) (*domain.Habit, error) {
	if strings.TrimSpace(habitID) == "" {
		return nil, ErrHabitValidation
	}
	habit, err := s.repo.ArchiveHabit(ctx, userID, habitID, archived)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrHabitNotFound
		}
		return nil, err
	}
	if habit == nil {
		return nil, ErrHabitNotFound
	}
	return habit, nil
}

func (s *HabitService) Reorder(ctx context.Context, userID string, ids []string) ([]domain.Habit, error) {
	if len(ids) == 0 {
		return nil, ErrHabitValidation
	}
	habits, err := s.repo.ReorderHabits(ctx, userID, ids)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrHabitForbidden
		}
		return nil, err
	}
	if habits == nil {
		return nil, ErrHabitForbidden
	}
	return habits, nil
}

func validateCreateHabit(req domain.CreateHabitRequest) error {
	if strings.TrimSpace(req.Name) == "" || len(req.Name) > 200 {
		return ErrHabitValidation
	}
	if req.Type != string(domain.HabitTypeBoolean) && req.Type != string(domain.HabitTypeNumeric) && req.Type != string(domain.HabitTypeDuration) {
		return ErrHabitValidation
	}
	if req.Frequency != string(domain.HabitFrequencyDaily) && req.Frequency != string(domain.HabitFrequencyWeekly) && req.Frequency != string(domain.HabitFrequencyCustom) {
		return ErrHabitValidation
	}
	if req.Type == string(domain.HabitTypeNumeric) {
		if req.TargetValue == nil || req.Unit == nil || strings.TrimSpace(*req.Unit) == "" {
			return ErrHabitValidation
		}
	}
	return nil
}

func validateUpdateHabit(req domain.UpdateHabitRequest) error {
	if req.Name != nil && (strings.TrimSpace(*req.Name) == "" || len(*req.Name) > 200) {
		return ErrHabitValidation
	}
	if req.Frequency != nil && *req.Frequency != string(domain.HabitFrequencyDaily) && *req.Frequency != string(domain.HabitFrequencyWeekly) && *req.Frequency != string(domain.HabitFrequencyCustom) {
		return ErrHabitValidation
	}
	if req.Type != "" {
		return ErrHabitValidation
	}
	return nil
}

func (s *HabitService) UpdateWithValidation(ctx context.Context, userID string, habitID string, req domain.UpdateHabitRequest) (*domain.Habit, error) {
	if err := validateUpdateHabit(req); err != nil {
		return nil, err
	}
	return s.Update(ctx, userID, habitID, req)
}

func (s *HabitService) CreateWithValidation(ctx context.Context, userID string, req domain.CreateHabitRequest) (*domain.Habit, error) {
	if err := validateCreateHabit(req); err != nil {
		return nil, err
	}
	return s.Create(ctx, userID, req)
}

func (s *HabitService) ListWithFilters(ctx context.Context, userID string, categoryID *string, archived *bool) ([]domain.Habit, error) {
	return s.List(ctx, userID, categoryID, archived)
}

func (s *HabitService) DeleteWithCascade(ctx context.Context, userID string, habitID string) error {
	return s.Delete(ctx, userID, habitID)
}

func init() {
	_ = fmt.Sprintf("%v", domain.HabitTypeBoolean)
}
