package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/qper/hf/internal/config"
	"github.com/qper/hf/internal/domain"
)

var (
	ErrEntryForbidden = errors.New("entry access denied")
)

type EntryRepository interface {
	CreateEntry(ctx context.Context, userID string, req domain.CreateEntryRequest) (*domain.Entry, error)
	UpdateEntry(ctx context.Context, userID string, entryID string, req domain.UpdateEntryRequest) (*domain.Entry, error)
	DeleteEntry(ctx context.Context, userID string, entryID string) (*domain.Entry, error)
	GetHabitByID(ctx context.Context, userID string, habitID string) (*domain.Habit, error)
}

type EntryService struct {
	repo   EntryRepository
	config config.Config
	userTZ *time.Location
}

func NewEntryService(repo EntryRepository, cfg config.Config) *EntryService {
	return &EntryService{repo: repo, config: cfg}
}

func (s *EntryService) Create(ctx context.Context, userID string, req domain.CreateEntryRequest) (*domain.Entry, error) {
	if err := s.validateDate(req.Date); err != nil {
		return nil, err
	}

	habit, err := s.repo.GetHabitByID(ctx, userID, req.HabitID)
	if err != nil {
		return nil, err
	}
	if habit == nil {
		return nil, ErrEntryForbidden
	}

	completed := false
	switch habit.Type {
	case domain.HabitTypeBoolean:
		if req.Completed != nil {
			completed = *req.Completed
		}
	case domain.HabitTypeNumeric:
		if req.Value != nil && habit.TargetValue != nil && *req.Value >= *habit.TargetValue {
			completed = true
		}
	case domain.HabitTypeDuration:
		if req.Value != nil && *req.Value > 0 {
			completed = true
		}
	}

	request := req
	request.Completed = &completed
	return s.repo.CreateEntry(ctx, userID, request)
}

func (s *EntryService) Update(ctx context.Context, userID string, entryID string, req domain.UpdateEntryRequest) (*domain.Entry, error) {
	return s.repo.UpdateEntry(ctx, userID, entryID, req)
}

func (s *EntryService) Delete(ctx context.Context, userID string, entryID string) (*domain.Entry, error) {
	return s.repo.DeleteEntry(ctx, userID, entryID)
}

func (s *EntryService) validateDate(date string) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return fmt.Errorf("invalid date")
	}

	loc := s.userTZ
	if loc == nil {
		loc = time.Now().Location()
	}
	parsed, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		return fmt.Errorf("invalid date")
	}
	parsed = parsed.In(loc)
	dayStart := time.Now().In(loc).AddDate(0, 0, -s.config.EditWindowDays)
	dayStart = time.Date(dayStart.Year(), dayStart.Month(), dayStart.Day(), 0, 0, 0, 0, loc)
	todayStart := time.Now().In(loc)
	todayStart = time.Date(todayStart.Year(), todayStart.Month(), todayStart.Day(), 0, 0, 0, 0, loc)
	if parsed.Before(dayStart) || parsed.After(todayStart) {
		return ErrEntryForbidden
	}
	return nil
}
