package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/qper/hf/internal/domain"
)

type HabitRepository struct {
	db *sql.DB
}

func NewHabitRepository(db *sql.DB) *HabitRepository {
	return &HabitRepository{db: db}
}

func (r *HabitRepository) CreateHabit(ctx context.Context, userID string, req domain.CreateHabitRequest) (*domain.Habit, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO habits (
			user_id,
			category_id,
			name,
			description,
			color,
			type,
			frequency,
			target_value,
			unit,
			sort_order,
			is_archived,
			is_deleted
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, user_id, category_id, name, description, color, type, frequency, target_value, unit, sort_order, is_archived, is_deleted, created_at, updated_at, deleted_at
	`, userID, req.CategoryID, req.Name, req.Description, req.Color, req.Type, req.Frequency, req.TargetValue, req.Unit, req.SortOrder, false, false)

	var h domain.Habit
	var categoryID sql.NullString
	var description sql.NullString
	var color sql.NullString
	var targetValue sql.NullFloat64
	var unit sql.NullString
	var deletedAt sql.NullTime
	if err := row.Scan(&h.ID, &h.UserID, &categoryID, &h.Name, &description, &color, &h.Type, &h.Frequency, &targetValue, &unit, &h.SortOrder, &h.IsArchived, &h.IsDeleted, &h.CreatedAt, &h.UpdatedAt, &deletedAt); err != nil {
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
	if deletedAt.Valid {
		h.DeletedAt = &deletedAt.Time
	}
	return &h, nil
}

func (r *HabitRepository) ListHabits(ctx context.Context, userID string, categoryID *string, archived *bool) ([]domain.Habit, error) {
	query := `
		SELECT id, user_id, category_id, name, description, color, type, frequency, target_value, unit, sort_order, is_archived, is_deleted, created_at, updated_at, deleted_at
		FROM habits
		WHERE user_id = $1 AND is_deleted = FALSE`
	args := []any{userID}
	argCount := 2

	if categoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argCount)
		args = append(args, *categoryID)
		argCount++
	}
	if archived != nil {
		query += fmt.Sprintf(" AND is_archived = $%d", argCount)
		args = append(args, *archived)
		argCount++
	}
	query += " ORDER BY sort_order, created_at"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var habits []domain.Habit
	for rows.Next() {
		var h domain.Habit
		var categoryID sql.NullString
		var description sql.NullString
		var color sql.NullString
		var targetValue sql.NullFloat64
		var unit sql.NullString
		var deletedAt sql.NullTime
		if err := rows.Scan(&h.ID, &h.UserID, &categoryID, &h.Name, &description, &color, &h.Type, &h.Frequency, &targetValue, &unit, &h.SortOrder, &h.IsArchived, &h.IsDeleted, &h.CreatedAt, &h.UpdatedAt, &deletedAt); err != nil {
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
		if deletedAt.Valid {
			h.DeletedAt = &deletedAt.Time
		}
		habits = append(habits, h)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return habits, nil
}

func (r *HabitRepository) GetHabitByID(ctx context.Context, userID string, habitID string) (*domain.Habit, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, category_id, name, description, color, type, frequency, target_value, unit, sort_order, is_archived, is_deleted, created_at, updated_at, deleted_at
		FROM habits
		WHERE id = $1 AND user_id = $2 AND is_deleted = FALSE
		LIMIT 1
	`, habitID, userID)

	var h domain.Habit
	var categoryID sql.NullString
	var description sql.NullString
	var color sql.NullString
	var targetValue sql.NullFloat64
	var unit sql.NullString
	var deletedAt sql.NullTime
	if err := row.Scan(&h.ID, &h.UserID, &categoryID, &h.Name, &description, &color, &h.Type, &h.Frequency, &targetValue, &unit, &h.SortOrder, &h.IsArchived, &h.IsDeleted, &h.CreatedAt, &h.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
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
	if deletedAt.Valid {
		h.DeletedAt = &deletedAt.Time
	}
	return &h, nil
}

func (r *HabitRepository) UpdateHabit(ctx context.Context, userID string, habitID string, req domain.UpdateHabitRequest) (*domain.Habit, error) {
	setParts := []string{"updated_at = NOW()"}
	args := []any{habitID, userID}
	argCount := 3

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argCount))
		args = append(args, *req.Name)
		argCount++
	}
	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argCount))
		args = append(args, *req.Description)
		argCount++
	}
	if req.Color != nil {
		setParts = append(setParts, fmt.Sprintf("color = $%d", argCount))
		args = append(args, *req.Color)
		argCount++
	}
	if req.CategoryID != nil {
		setParts = append(setParts, fmt.Sprintf("category_id = $%d", argCount))
		args = append(args, *req.CategoryID)
		argCount++
	}
	if req.Frequency != nil {
		setParts = append(setParts, fmt.Sprintf("frequency = $%d", argCount))
		args = append(args, *req.Frequency)
		argCount++
	}
	if req.TargetValue != nil {
		setParts = append(setParts, fmt.Sprintf("target_value = $%d", argCount))
		args = append(args, *req.TargetValue)
		argCount++
	}
	if req.Unit != nil {
		setParts = append(setParts, fmt.Sprintf("unit = $%d", argCount))
		args = append(args, *req.Unit)
		argCount++
	}
	if req.SortOrder != nil {
		setParts = append(setParts, fmt.Sprintf("sort_order = $%d", argCount))
		args = append(args, *req.SortOrder)
		argCount++
	}
	if req.Archived != nil {
		setParts = append(setParts, fmt.Sprintf("is_archived = $%d", argCount))
		args = append(args, *req.Archived)
		argCount++
	}

	query := fmt.Sprintf(`
		UPDATE habits
		SET %s
		WHERE id = $1 AND user_id = $2 AND is_deleted = FALSE
		RETURNING id, user_id, category_id, name, description, color, type, frequency, target_value, unit, sort_order, is_archived, is_deleted, created_at, updated_at, deleted_at
	`, strings.Join(setParts, ", "))

	row := r.db.QueryRowContext(ctx, query, args...)
	var h domain.Habit
	var categoryID sql.NullString
	var description sql.NullString
	var color sql.NullString
	var targetValue sql.NullFloat64
	var unit sql.NullString
	var deletedAt sql.NullTime
	if err := row.Scan(&h.ID, &h.UserID, &categoryID, &h.Name, &description, &color, &h.Type, &h.Frequency, &targetValue, &unit, &h.SortOrder, &h.IsArchived, &h.IsDeleted, &h.CreatedAt, &h.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
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
	if deletedAt.Valid {
		h.DeletedAt = &deletedAt.Time
	}
	return &h, nil
}

func (r *HabitRepository) DeleteHabit(ctx context.Context, userID string, habitID string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE habits
		SET is_deleted = TRUE, deleted_at = $3, updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_deleted = FALSE
	`, habitID, userID, time.Now().UTC())
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
