package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/qper/hf/internal/config"
	"github.com/qper/hf/internal/domain"
)

type stubEntryRepository struct {
	entry *domain.Entry
	err   error
}

func (s stubEntryRepository) CreateEntry(ctx context.Context, userID string, req domain.CreateEntryRequest) (*domain.Entry, error) {
	return s.entry, s.err
}

func (s stubEntryRepository) UpdateEntry(ctx context.Context, userID string, entryID string, req domain.UpdateEntryRequest) (*domain.Entry, error) {
	return s.entry, s.err
}

func (s stubEntryRepository) DeleteEntry(ctx context.Context, userID string, entryID string) (*domain.Entry, error) {
	return s.entry, s.err
}

func TestEntryServiceRejectsOutOfWindowDates(t *testing.T) {
	loc := time.FixedZone("UTC+2", 2*60*60)
	svc := NewEntryService(stubEntryRepository{}, config.Config{EditWindowDays: 1})

	today := time.Now().In(loc)
	cases := []struct {
		name      string
		date      string
		wantError error
	}{
		{name: "today", date: today.Format("2006-01-02"), wantError: nil},
		{name: "yesterday", date: today.AddDate(0, 0, -1).Format("2006-01-02"), wantError: nil},
		{name: "two days ago", date: today.AddDate(0, 0, -2).Format("2006-01-02"), wantError: ErrEntryForbidden},
		{name: "tomorrow", date: today.AddDate(0, 0, 1).Format("2006-01-02"), wantError: ErrEntryForbidden},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.Create(context.Background(), "u1", domain.CreateEntryRequest{HabitID: "h1", Date: tc.date})
			if tc.wantError == nil {
				if err != nil {
					t.Fatalf("expected nil error, got %v", err)
				}
				return
			}
			if !errors.Is(err, tc.wantError) {
				t.Fatalf("expected %v, got %v", tc.wantError, err)
			}
		})
	}
}

func TestEntryServiceUsesHabitCompletionRules(t *testing.T) {
	loc := time.FixedZone("UTC+2", 2*60*60)
	svc := NewEntryService(stubEntryRepository{}, config.Config{EditWindowDays: 1})

	today := time.Now().In(loc).Format("2006-01-02")
	_, err := svc.Create(context.Background(), "u1", domain.CreateEntryRequest{HabitID: "h1", Date: today, Completed: true})
	if !errors.Is(err, ErrEntryForbidden) {
		t.Fatalf("expected ErrEntryForbidden for missing habit, got %v", err)
	}
}
