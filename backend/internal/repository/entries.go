package repository

import (
	"context"
	"database/sql"

	"github.com/qper/hf/internal/domain"
)

type EntryRepository struct {
	db *sql.DB
}

func NewEntryRepository(db *sql.DB) *EntryRepository {
	return &EntryRepository{db: db}
}

func (r *EntryRepository) CreateEntry(ctx context.Context, userID string, req domain.CreateEntryRequest) (*domain.Entry, error) {
	var completed bool
	if req.Completed != nil {
		completed = *req.Completed
	}

	var value sql.NullFloat64
	if req.Value != nil {
		value.Valid = true
		value.Float64 = *req.Value
	}

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO habit_entries (habit_id, entry_date, completed, value, note, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (habit_id, entry_date) DO UPDATE SET
			completed = EXCLUDED.completed,
			value = EXCLUDED.value,
			note = EXCLUDED.note,
			updated_at = NOW()
		RETURNING id, habit_id, entry_date, completed, value, note, created_at, updated_at
	`, req.HabitID, req.Date, completed, value, req.Note)

	var entry domain.Entry
	var note sql.NullString
	if err := row.Scan(&entry.ID, &entry.HabitID, &entry.Date, &entry.Completed, &value, &note, &entry.CreatedAt, &entry.UpdatedAt); err != nil {
		return nil, err
	}
	if value.Valid {
		entryValue := value.Float64
		entry.Value = &entryValue
	}
	if note.Valid {
		noteValue := note.String
		entry.Note = &noteValue
	}
	entry.Date = req.Date
	return &entry, nil
}

func (r *EntryRepository) UpdateEntry(ctx context.Context, userID string, entryID string, req domain.UpdateEntryRequest) (*domain.Entry, error) {
	var value sql.NullFloat64
	if req.Value != nil {
		value.Valid = true
		value.Float64 = *req.Value
	}

	row := r.db.QueryRowContext(ctx, `
		UPDATE habit_entries
		SET updated_at = NOW(), value = COALESCE($3, value), note = COALESCE($4, note)
		WHERE id = $1
		RETURNING id, habit_id, entry_date, completed, value, note, created_at, updated_at
	`, entryID, userID, value, req.Note)

	var entry domain.Entry
	var note sql.NullString
	if err := row.Scan(&entry.ID, &entry.HabitID, &entry.Date, &entry.Completed, &value, &note, &entry.CreatedAt, &entry.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if value.Valid {
		entryValue := value.Float64
		entry.Value = &entryValue
	}
	if note.Valid {
		noteValue := note.String
		entry.Note = &noteValue
	}
	return &entry, nil
}

func (r *EntryRepository) DeleteEntry(ctx context.Context, userID string, entryID string) (*domain.Entry, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE habit_entries
		SET completed = FALSE, value = NULL, note = NULL, updated_at = NOW()
		WHERE id = $1
		RETURNING id, habit_id, entry_date, completed, value, note, created_at, updated_at
	`, entryID)

	var entry domain.Entry
	var value sql.NullFloat64
	var note sql.NullString
	if err := row.Scan(&entry.ID, &entry.HabitID, &entry.Date, &entry.Completed, &value, &note, &entry.CreatedAt, &entry.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if value.Valid {
		entryValue := value.Float64
		entry.Value = &entryValue
	}
	if note.Valid {
		noteValue := note.String
		entry.Note = &noteValue
	}
	return &entry, nil
}

func (r *EntryRepository) GetHabitByID(ctx context.Context, userID string, habitID string) (*domain.Habit, error) {
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
