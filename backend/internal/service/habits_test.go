package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/qper/hf/internal/domain"
)

type stubHabitRepository struct {
	created *domain.Habit
	got     *domain.Habit
}

func (s stubHabitRepository) CreateHabit(ctx context.Context, userID string, req domain.CreateHabitRequest) (*domain.Habit, error) {
	return s.created, nil
}

func (s stubHabitRepository) ListHabits(ctx context.Context, userID string, categoryID *string, archived *bool) ([]domain.Habit, error) {
	return nil, nil
}

func (s stubHabitRepository) GetHabitByID(ctx context.Context, userID string, habitID string) (*domain.Habit, error) {
	return s.got, nil
}

func (s stubHabitRepository) UpdateHabit(ctx context.Context, userID string, habitID string, req domain.UpdateHabitRequest) (*domain.Habit, error) {
	return s.got, nil
}

func (s stubHabitRepository) DeleteHabit(ctx context.Context, userID string, habitID string) error {
	return nil
}

func (s stubHabitRepository) ArchiveHabit(ctx context.Context, userID string, habitID string, archived bool) (*domain.Habit, error) {
	return nil, nil
}

func (s stubHabitRepository) ReorderHabits(ctx context.Context, userID string, ids []string) ([]domain.Habit, error) {
	return nil, nil
}

func TestCreateHabitRejectsInvalidPayload(t *testing.T) {
	svc := NewHabitService(stubHabitRepository{})

	_, err := svc.Create(context.Background(), "u1", domain.CreateHabitRequest{Name: "", Type: "invalid", Frequency: "daily"})
	if !errors.Is(err, ErrHabitValidation) {
		t.Fatalf("expected ErrHabitValidation, got %v", err)
	}
}

func TestCreateHabitRequiresNumericTargetData(t *testing.T) {
	svc := NewHabitService(stubHabitRepository{})

	_, err := svc.Create(context.Background(), "u1", domain.CreateHabitRequest{
		Name:      "Water",
		Type:      string(domain.HabitTypeNumeric),
		Frequency: string(domain.HabitFrequencyDaily),
		Unit:      nil,
	})
	if !errors.Is(err, ErrHabitValidation) {
		t.Fatalf("expected ErrHabitValidation, got %v", err)
	}
}

func TestUpdateHabitRejectsTypeChange(t *testing.T) {
	svc := NewHabitService(stubHabitRepository{})

	_, err := svc.Update(context.Background(), "u1", "h1", domain.UpdateHabitRequest{Type: "numeric"})
	if !errors.Is(err, ErrHabitValidation) {
		t.Fatalf("expected ErrHabitValidation, got %v", err)
	}
}

func TestGetHabitReturnsErrNotFoundForMissingHabit(t *testing.T) {
	svc := NewHabitService(stubHabitRepository{})

	_, err := svc.GetByID(context.Background(), "u1", "missing")
	if !errors.Is(err, ErrHabitNotFound) {
		t.Fatalf("expected ErrHabitNotFound, got %v", err)
	}
}

func TestDeleteHabitMarksHabitDeleted(t *testing.T) {
	svc := NewHabitService(stubHabitRepository{got: &domain.Habit{ID: "h1", UserID: "u1", Name: "Read", Type: domain.HabitTypeBoolean, Frequency: domain.HabitFrequencyDaily, CreatedAt: time.Now(), UpdatedAt: time.Now()}})

	err := svc.Delete(context.Background(), "u1", "h1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
