package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/qper/hf/internal/domain"
)

type BoardRepository struct {
	db *sql.DB
}

func NewBoardRepository(db *sql.DB) *BoardRepository {
	return &BoardRepository{db: db}
}

func (r *BoardRepository) ListBoardHabits(ctx context.Context, userID string, targetDate time.Time) ([]domain.BoardHabit, error) {
	query := `
		SELECT h.id, h.user_id, h.category_id, h.name, h.description, h.color, h.type, h.frequency, h.target_value, h.unit, h.sort_order,
		       EXISTS (
				SELECT 1 FROM habit_entries he
				WHERE he.habit_id = h.id AND he.entry_date = $2 AND he.completed = TRUE
			) AS is_completed,
		       calculate_streak_for_habit(h.id, $2) AS streak
		FROM habits h
		WHERE h.user_id = $1 AND h.is_deleted = FALSE AND h.is_archived = FALSE
		ORDER BY h.sort_order, h.created_at
	`

	rows, err := r.db.QueryContext(ctx, query, userID, targetDate)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var habits []domain.BoardHabit
	for rows.Next() {
		var h domain.BoardHabit
		var categoryID sql.NullString
		var description sql.NullString
		var color sql.NullString
		var targetValue sql.NullFloat64
		var unit sql.NullString
		if err := rows.Scan(&h.ID, &h.UserID, &categoryID, &h.Name, &description, &color, &h.Type, &h.Frequency, &targetValue, &unit, &h.SortOrder, &h.IsCompleted, &h.Streak); err != nil {
			return nil, err
		}
		if categoryID.Valid {
			value := categoryID.String
			h.CategoryID = &value
		}
		if description.Valid {
			value := description.String
			h.Description = &value
		}
		if color.Valid {
			value := color.String
			h.Color = &value
		}
		if targetValue.Valid {
			value := targetValue.Float64
			h.TargetValue = &value
		}
		if unit.Valid {
			value := unit.String
			h.Unit = &value
		}
		habits = append(habits, h)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return habits, nil
}
