package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/qper/hf/internal/config"
	"github.com/qper/hf/internal/domain"
)

type stubBoardRepository struct {
	habits []domain.BoardHabit
	err    error
}

func (s stubBoardRepository) ListBoardHabits(ctx context.Context, userID string, targetDate time.Time) ([]domain.BoardHabit, error) {
	return s.habits, s.err
}

func TestBoardServiceComputesEditableWindow(t *testing.T) {
	loc := time.FixedZone("UTC+2", 2*60*60)
	svc := NewBoardService(stubBoardRepository{habits: []domain.BoardHabit{{ID: "h1"}}}, config.Config{EditWindowDays: 1})

	today := time.Now().In(loc)
	cases := []struct {
		name         string
		target       time.Time
		wantEditable bool
	}{
		{name: "today", target: today, wantEditable: true},
		{name: "yesterday", target: today.AddDate(0, 0, -1), wantEditable: true},
		{name: "two days ago", target: today.AddDate(0, 0, -2), wantEditable: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			board, err := svc.GetBoard(context.Background(), "u1", tc.target.Format("2006-01-02"), loc)
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if board.IsEditable != tc.wantEditable {
				t.Fatalf("expected is_editable=%v, got %v", tc.wantEditable, board.IsEditable)
			}
		})
	}
}

func TestBoardServiceRejectsFutureDates(t *testing.T) {
	loc := time.FixedZone("UTC+2", 2*60*60)
	svc := NewBoardService(stubBoardRepository{}, config.Config{EditWindowDays: 1})

	future := time.Now().In(loc).AddDate(0, 0, 1)
	_, err := svc.GetBoard(context.Background(), "u1", future.Format("2006-01-02"), loc)
	if !errors.Is(err, ErrBoardFutureDate) {
		t.Fatalf("expected ErrBoardFutureDate, got %v", err)
	}
}

func TestBoardServiceCalculatesProgressTotals(t *testing.T) {
	loc := time.FixedZone("UTC+2", 2*60*60)
	svc := NewBoardService(stubBoardRepository{habits: []domain.BoardHabit{{ID: "h1", IsCompleted: true}, {ID: "h2", IsCompleted: false}}}, config.Config{EditWindowDays: 1})

	today := time.Now().In(loc)
	board, err := svc.GetBoard(context.Background(), "u1", today.Format("2006-01-02"), loc)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if board.Progress.Done != 1 || board.Progress.Total != 2 {
		t.Fatalf("expected progress {done:1,total:2}, got %+v", board.Progress)
	}
}
