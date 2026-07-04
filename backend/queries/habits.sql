-- name: CreateHabit :one
INSERT INTO habits (
    user_id,
    category_id,
    name,
    description,
    habit_type,
    habit_freq,
    target_value,
    unit,
    is_active,
    start_date,
    sort_order
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: GetHabits :many
SELECT h.*
FROM habits h
WHERE h.user_id = $1 AND h.deleted_at IS NULL
ORDER BY h.sort_order, h.created_at;

-- name: GetHabitByID :one
SELECT *
FROM habits
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateHabit :one
UPDATE habits
SET name = $2,
    description = $3,
    habit_type = $4,
    habit_freq = $5,
    target_value = $6,
    unit = $7,
    is_active = $8,
    start_date = $9,
    category_id = $10,
    sort_order = $11,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteHabit :exec
UPDATE habits
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: ArchiveHabit :exec
UPDATE habits
SET is_active = FALSE, updated_at = NOW()
WHERE id = $1;

-- name: ReorderHabits :exec
UPDATE habits
SET sort_order = $2, updated_at = NOW()
WHERE id = $1;
