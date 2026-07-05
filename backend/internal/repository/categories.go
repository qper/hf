package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/qper/hf/internal/domain"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, userID string, req domain.CreateCategoryRequest) (*domain.Category, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO categories (user_id, name, color, icon, sort_order)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, name, color, icon, sort_order, created_at, updated_at, deleted_at
	`, userID, req.Name, req.Color, req.Icon, req.SortOrder)

	var c domain.Category
	var deletedAt sql.NullTime
	if err := row.Scan(&c.ID, &c.UserID, &c.Name, &c.Color, &c.Icon, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt, &deletedAt); err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique") {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	if deletedAt.Valid {
		c.DeletedAt = &deletedAt.Time
	}
	return &c, nil
}

func (r *CategoryRepository) ListCategories(ctx context.Context, userID string) ([]domain.Category, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, name, color, icon, sort_order, created_at, updated_at, deleted_at
		FROM categories
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY sort_order, created_at
	`, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		var deletedAt sql.NullTime
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Color, &c.Icon, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt, &deletedAt); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			c.DeletedAt = &deletedAt.Time
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, userID string, categoryID string) (*domain.Category, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, name, color, icon, sort_order, created_at, updated_at, deleted_at
		FROM categories
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
		LIMIT 1
	`, categoryID, userID)

	var c domain.Category
	var deletedAt sql.NullTime
	if err := row.Scan(&c.ID, &c.UserID, &c.Name, &c.Color, &c.Icon, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if deletedAt.Valid {
		c.DeletedAt = &deletedAt.Time
	}
	return &c, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, userID string, categoryID string, req domain.UpdateCategoryRequest) (*domain.Category, error) {
	setParts := []string{"updated_at = NOW()"}
	args := []any{categoryID, userID}

	if req.Name != nil {
		placeholder := len(args) + 1
		setParts = append(setParts, fmt.Sprintf("name = $%d", placeholder))
		args = append(args, *req.Name)
	}
	if req.Color != nil {
		placeholder := len(args) + 1
		setParts = append(setParts, fmt.Sprintf("color = $%d", placeholder))
		args = append(args, *req.Color)
	}
	if req.Icon != nil {
		placeholder := len(args) + 1
		setParts = append(setParts, fmt.Sprintf("icon = $%d", placeholder))
		args = append(args, *req.Icon)
	}
	if req.SortOrder != nil {
		placeholder := len(args) + 1
		setParts = append(setParts, fmt.Sprintf("sort_order = $%d", placeholder))
		args = append(args, *req.SortOrder)
	}

	query := fmt.Sprintf(`
		UPDATE categories
		SET %s
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
		RETURNING id, user_id, name, color, icon, sort_order, created_at, updated_at, deleted_at
	`, strings.Join(setParts, ", "))

	row := r.db.QueryRowContext(ctx, query, args...)
	var c domain.Category
	var deletedAt sql.NullTime
	if err := row.Scan(&c.ID, &c.UserID, &c.Name, &c.Color, &c.Icon, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique") {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	if deletedAt.Valid {
		c.DeletedAt = &deletedAt.Time
	}
	return &c, nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, userID string, categoryID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	result, err := tx.ExecContext(ctx, `
		UPDATE categories
		SET deleted_at = $3, updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`, categoryID, userID, time.Now().UTC())
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

	if _, err := tx.ExecContext(ctx, `
		UPDATE habits
		SET category_id = NULL, updated_at = NOW()
		WHERE category_id = $1 AND user_id = $2 AND is_deleted = FALSE
	`, categoryID, userID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *CategoryRepository) ReorderCategories(ctx context.Context, userID string, ids []string) ([]domain.Category, error) {
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
				SELECT 1 FROM categories WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
			)
		`, id, userID).Scan(&exists)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, sql.ErrNoRows
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE categories
			SET sort_order = $1, updated_at = NOW()
			WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL
		`, idx+1, id, userID); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.ListCategories(ctx, userID)
}
