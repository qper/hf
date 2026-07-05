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

	if categoryID != nil {
		placeholder := len(args) + 1
		query += fmt.Sprintf(" AND category_id = $%d", placeholder)
		args = append(args, *categoryID)
	}
	if archived != nil {
		placeholder := len(args) + 1
		query += fmt.Sprintf(" AND is_archived = $%d", placeholder)
		args = append(args, *archived)
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

	if req.Name != nil {
		placeholder := len(args) + 1
		setParts = append(setParts, fmt.Sprintf("name = $%d", placeholder))
		args = append(args, *req.Name)
	}
	if req.Description != nil {
		placeholder := len(args) + 1
		setParts = append(setParts, fmt.Sprintf("description = $%d", placeholder))
		args = append(args, *req.Description)
	}
	if req.Color != nil {
		placeholder := len(args) + 1
		setParts = append(setParts, fmt.Sprintf("color = $%d", placeholder))
		args = append(args, *req.Color)
	}
	if req.CategoryID != nil {
		placeholder := len(args) + 1
		setParts = append(setParts, fmt.Sprintf("category_id = $%d", placeholder))
		args = append(args, *req.CategoryID)
	}
	if req.Frequency != nil {
		placeholder := len(args) + 1
		setParts = append(setParts, fmt.Sprintf("frequency = $%d", placeholder))
		args = append(args, *req.Frequency)
	}
	if req.TargetValue != nil {
		placeholder := len(args) + 1
		setParts = append(setParts, fmt.Sprintf("target_value = $%d", placeholder))
		args = append(args, *req.TargetValue)
	}
	if req.Unit != nil {
		placeholder := len(args) + 1
		setParts = append(setParts, fmt.Sprintf("unit = $%d", placeholder))
		args = append(args, *req.Unit)
	}
	if req.SortOrder != nil {
		placeholder := len(args) + 1
		setParts = append(setParts, fmt.Sprintf("sort_order = $%d", placeholder))
		args = append(args, *req.SortOrder)
	}
	if req.Archived != nil {
		placeholder := len(args) + 1
		setParts = append(setParts, fmt.Sprintf("is_archived = $%d", placeholder))
		args = append(args, *req.Archived)
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

func (r *HabitRepository) ArchiveHabit(ctx context.Context, userID string, habitID string, archived bool) (*domain.Habit, error) {
	var archivedAt interface{}
	if archived {
		archivedAt = time.Now().UTC()
	} else {
		archivedAt = nil
	}

	row := r.db.QueryRowContext(ctx, `
		UPDATE habits
		SET is_archived = $3, updated_at = NOW(), archived_at = $4
		WHERE id = $1 AND user_id = $2 AND is_deleted = FALSE
		RETURNING id, user_id, category_id, name, description, color, type, frequency, target_value, unit, sort_order, is_archived, is_deleted, created_at, updated_at, deleted_at
	`, habitID, userID, archived, archivedAt)

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

func (r *HabitRepository) ReorderHabits(ctx context.Context, userID string, ids []string) ([]domain.Habit, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for idx, id := range ids {
		var exists bool
		err := tx.QueryRowContext(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM habits WHERE id = $1 AND user_id = $2 AND is_deleted = FALSE
			)
		`, id, userID).Scan(&exists)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, sql.ErrNoRows
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE habits
			SET sort_order = $1, updated_at = NOW()
			WHERE id = $2 AND user_id = $3 AND is_deleted = FALSE
		`, idx+1, id, userID); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.ListHabits(ctx, userID, nil, nil)
}
