package service

import (
	"context"
	"errors"
	"time"

	"github.com/qper/hf/internal/config"
	"github.com/qper/hf/internal/domain"
)

var (
	ErrBoardFutureDate = errors.New("future date is not allowed")
)

type BoardRepository interface {
	ListBoardHabits(ctx context.Context, userID string, targetDate time.Time) ([]domain.BoardHabit, error)
}

type BoardService struct {
	repo   BoardRepository
	config config.Config
}

func NewBoardService(repo BoardRepository, cfg config.Config) *BoardService {
	return &BoardService{repo: repo, config: cfg}
}

func (s *BoardService) GetBoard(ctx context.Context, userID string, date string, userTZ *time.Location) (*domain.Board, error) {
	if userTZ == nil {
		userTZ = time.UTC
	}

	parsed, err := time.ParseInLocation("2006-01-02", date, userTZ)
	if err != nil {
		return nil, err
	}

	today := time.Now().In(userTZ)
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, userTZ)
	targetDate := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, userTZ)
	if targetDate.After(todayDate) {
		return nil, ErrBoardFutureDate
	}

	habits, err := s.repo.ListBoardHabits(ctx, userID, targetDate)
	if err != nil {
		return nil, err
	}

	var done int
	for _, habit := range habits {
		if habit.IsCompleted {
			done++
		}
	}
	windowStart := todayDate.AddDate(0, 0, -s.config.EditWindowDays)
	isEditable := !targetDate.Before(windowStart) && !targetDate.After(todayDate)

	return &domain.Board{
		Date:       targetDate.Format("2006-01-02"),
		IsEditable: isEditable,
		Progress: domain.BoardProgress{
			Done:  done,
			Total: len(habits),
		},
		Habits: habits,
	}, nil
}
