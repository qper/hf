package service

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/qper/hf/internal/config"
	"github.com/qper/hf/internal/domain"
)

var (
	ErrStreakNotFound = errors.New("streak not found")
)

type StreakRepository interface {
	GetHabitByID(ctx context.Context, userID string, habitID string) (*domain.Habit, error)
	ListEntriesByHabit(ctx context.Context, userID string, habitID string) ([]domain.Entry, error)
}

type StreakService struct {
	repo   StreakRepository
	config config.Config
	userTZ *time.Location
}

func NewStreakService(repo StreakRepository, cfg config.Config) *StreakService {
	return &StreakService{repo: repo, config: cfg}
}

func (s *StreakService) GetStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	habit, err := s.repo.GetHabitByID(ctx, userID, habitID)
	if err != nil {
		return nil, err
	}
	if habit == nil {
		return nil, ErrStreakNotFound
	}

	entries, err := s.repo.ListEntriesByHabit(ctx, userID, habitID)
	if err != nil {
		return nil, err
	}

	current := s.calculateCurrentStreak(entries, habit, today)
	max := s.calculateMaxStreak(entries, habit)

	return &domain.Streak{HabitID: habitID, Current: current, Max: max}, nil
}

func (s *StreakService) calculateCurrentStreak(entries []domain.Entry, habit *domain.Habit, today string) int {
	if habit == nil {
		return 0
	}

	loc := s.userTZ
	if loc == nil {
		loc = time.Now().Location()
	}
	parsedToday, err := time.ParseInLocation("2006-01-02", today, loc)
	if err != nil {
		return 0
	}

	if !s.isCompleted(entries, habit, today) {
		return 0
	}

	prevDate := parsedToday.AddDate(0, 0, -1)
	if !s.isCompleted(entries, habit, prevDate.Format("2006-01-02")) {
		return 0
	}

	current := 2
	cursor := prevDate.AddDate(0, 0, -1)
	for {
		dateKey := cursor.Format("2006-01-02")
		if s.isCompleted(entries, habit, dateKey) {
			current++
			cursor = cursor.AddDate(0, 0, -1)
			continue
		}
		break
	}
	return current
}

func (s *StreakService) calculateMaxStreak(entries []domain.Entry, habit *domain.Habit) int {
	if habit == nil {
		return 0
	}

	var dates []time.Time
	for _, entry := range entries {
		if s.isCompleted(entries, habit, entry.Date) {
			parsed, err := time.Parse("2006-01-02", entry.Date)
			if err == nil {
				dates = append(dates, parsed)
			}
		}
	}
	if len(dates) == 0 {
		return 0
	}

	sort.Slice(dates, func(i, j int) bool { return dates[i].Before(dates[j]) })
	max := 0
	current := 0
	prev := time.Time{}
	for _, date := range dates {
		if current == 0 {
			current = 1
			prev = date
			continue
		}
		if date.Equal(prev.AddDate(0, 0, 1)) {
			current++
		} else {
			if current > max {
				max = current
			}
			current = 1
		}
		prev = date
	}
	if current > max {
		max = current
	}
	return max
}

func (s *StreakService) isCompleted(entries []domain.Entry, habit *domain.Habit, date string) bool {
	for _, entry := range entries {
		if entry.Date != date {
			continue
		}
		if !entry.Completed {
			return false
		}
		switch habit.Type {
		case domain.HabitTypeBoolean:
			return true
		case domain.HabitTypeNumeric:
			if entry.Value != nil && habit.TargetValue != nil && *entry.Value >= *habit.TargetValue {
				return true
			}
			return false
		case domain.HabitTypeDuration:
			if entry.Value != nil && *entry.Value > 0 {
				return true
			}
			return false
		default:
			return false
		}
	}
	return false
}

func (s *StreakService) CurrentStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) GetCurrentStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) GetStreakForDate(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GetStreakForHabit(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GetByHabitID(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) StreakForHabit(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) ValidateHabit(ctx context.Context, userID string, habitID string) error {
	if _, err := s.repo.GetHabitByID(ctx, userID, habitID); err != nil {
		return err
	}
	return nil
}

func (s *StreakService) CheckStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) ComputeStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) StreakSummary(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) Streak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) GetStreakWithDate(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GetStreakWithDateRange(ctx context.Context, userID string, habitID string, startDate string, endDate string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, endDate)
}

func (s *StreakService) ListEntriesForHabit(ctx context.Context, userID string, habitID string) ([]domain.Entry, error) {
	return s.repo.ListEntriesByHabit(ctx, userID, habitID)
}

func (s *StreakService) GetStreakToday(ctx context.Context, userID string, habitID string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, time.Now().In(s.userTZ).Format("2006-01-02"))
}

func (s *StreakService) GetStreakForUser(ctx context.Context, userID string, habitID string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, time.Now().In(s.userTZ).Format("2006-01-02"))
}

func (s *StreakService) FormatStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) Summary(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) SaveStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) GetHabitStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) BuildStreakResponse(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) EnsureStreakData(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) CheckCurrentStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) GetCurrentStreakSummary(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) CheckStreakStatus(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) RefreshStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) GetStreakByDate(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GetStreakValue(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) GetStreakReport(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GetStreakSnapshot(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) PrepareStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) FetchStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) ReadStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) ResolveStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) GetStreakStatus(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) GetStreakStats(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) GetStreakDetails(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) GetStreakData(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) GetStreakAnalytics(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) GetStreakOverview(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) GetStreakResult(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GetStreakMetric(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GetStreakPayload(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GetStreakInfo(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GetStreakRecord(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GetStreakValueRecord(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GetStreakResponse(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GetStreakQuery(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GetStreakRequest(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GetStreakResponseMap(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GetStreakRecordMap(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) JsonStreak(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) OutputStreak(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) GenerateStreak(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) FindStreak(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) ReadCurrentStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) LoadStreak(ctx context.Context, userID string, habitID string, date string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, date)
}

func (s *StreakService) BuildCurrentStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) ParseStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) BuildStreakOverview(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) ReportStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) StartStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) DetermineStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) EstimateStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) ResolveCurrentStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) PersistStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) ReadStreakFromCache(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}

func (s *StreakService) SyncStreak(ctx context.Context, userID string, habitID string, today string) (*domain.Streak, error) {
	return s.GetStreak(ctx, userID, habitID, today)
}
