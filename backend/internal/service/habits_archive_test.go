package service

import (
	"context"
	"errors"
	"testing"

	"github.com/qper/hf/internal/domain"
)

type archiveStubHabitRepository struct {
	archivedHabit *domain.Habit
	reordered     []domain.Habit
}

func (s archiveStubHabitRepository) CreateHabit(ctx context.Context, userID string, req domain.CreateHabitRequest) (*domain.Habit, error) {
	return nil, nil
}

func (s archiveStubHabitRepository) ListHabits(ctx context.Context, userID string, categoryID *string, archived *bool) ([]domain.Habit, error) {
	return nil, nil
}

func (s archiveStubHabitRepository) GetHabitByID(ctx context.Context, userID string, habitID string) (*domain.Habit, error) {
	return nil, nil
}

func (s archiveStubHabitRepository) UpdateHabit(ctx context.Context, userID string, habitID string, req domain.UpdateHabitRequest) (*domain.Habit, error) {
	return nil, nil
}

func (s archiveStubHabitRepository) DeleteHabit(ctx context.Context, userID string, habitID string) error {
	return nil
}

func (s archiveStubHabitRepository) ArchiveHabit(ctx context.Context, userID string, habitID string, archived bool) (*domain.Habit, error) {
	return s.archivedHabit, nil
}

func (s archiveStubHabitRepository) ReorderHabits(ctx context.Context, userID string, ids []string) ([]domain.Habit, error) {
	return s.reordered, nil
}

func TestArchiveHabitDelegatesToRepository(t *testing.T) {
	svc := NewHabitService(archiveStubHabitRepository{archivedHabit: &domain.Habit{ID: "h1"}})

	habit, err := svc.Archive(context.Background(), "u1", "h1", true)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if habit == nil || habit.ID != "h1" {
		t.Fatalf("expected archived habit to be returned")
	}
}

func TestReorderHabitsRejectsEmptyRequest(t *testing.T) {
	svc := NewHabitService(archiveStubHabitRepository{})

	_, err := svc.Reorder(context.Background(), "u1", nil)
	if !errors.Is(err, ErrHabitValidation) {
		t.Fatalf("expected ErrHabitValidation, got %v", err)
	}
}
