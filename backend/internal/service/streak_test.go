package service

import (
	"context"
	"testing"
	"time"

	"github.com/qper/hf/internal/config"
	"github.com/qper/hf/internal/domain"
)

type stubStreakRepository struct {
	habits  []domain.Habit
	entries []domain.Entry
	err     error
}

func (s stubStreakRepository) GetHabitByID(ctx context.Context, userID string, habitID string) (*domain.Habit, error) {
	for _, habit := range s.habits {
		if habit.ID == habitID {
			return &habit, nil
		}
	}
	return nil, nil
}

func (s stubStreakRepository) ListEntriesByHabit(ctx context.Context, userID string, habitID string) ([]domain.Entry, error) {
	return s.entries, s.err
}

func TestStreakServiceCalculatesCurrentAndMaxStreaks(t *testing.T) {
	loc := time.FixedZone("UTC+2", 2*60*60)
	svc := NewStreakService(stubStreakRepository{habits: []domain.Habit{{ID: "h1", Type: domain.HabitTypeBoolean}}, entries: []domain.Entry{{Date: "2024-01-01", Completed: true}, {Date: "2024-01-02", Completed: true}, {Date: "2024-01-03", Completed: true}, {Date: "2024-01-04", Completed: true}, {Date: "2024-01-05", Completed: true}}}, config.Config{})
	svc.userTZ = loc

	streak, err := svc.GetStreak(context.Background(), "u1", "h1", "2024-01-05")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if streak.Current != 5 || streak.Max != 5 {
		t.Fatalf("expected current/max streak 5, got %+v", streak)
	}
}

func TestStreakServiceResetsCurrentAfterGap(t *testing.T) {
	loc := time.FixedZone("UTC+2", 2*60*60)
	svc := NewStreakService(stubStreakRepository{habits: []domain.Habit{{ID: "h1", Type: domain.HabitTypeBoolean}}, entries: []domain.Entry{{Date: "2024-01-01", Completed: true}, {Date: "2024-01-02", Completed: true}, {Date: "2024-01-04", Completed: true}}}, config.Config{})
	svc.userTZ = loc

	streak, err := svc.GetStreak(context.Background(), "u1", "h1", "2024-01-04")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if streak.Current != 0 {
		t.Fatalf("expected current streak 0 after a gap, got %d", streak.Current)
	}
}
