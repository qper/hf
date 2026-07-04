-- name: CreateHabit :one
INSERT INTO habits (
    user_id,
    category_id,
    name,
    description,
    color,
    sort_order,
    is_archived,
    is_deleted
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetHabits :many
SELECT h.*
FROM habits h
WHERE h.user_id = $1 AND h.is_deleted = FALSE
ORDER BY h.sort_order, h.created_at;

-- name: GetHabitByID :one
SELECT *
FROM habits
WHERE id = $1 AND is_deleted = FALSE
LIMIT 1;

-- name: UpdateHabit :one
UPDATE habits
SET name = $2,
    description = $3,
    color = $4,
    category_id = $5,
    sort_order = $6,
    updated_at = NOW()
WHERE id = $1 AND is_deleted = FALSE
RETURNING *;

-- name: SoftDeleteHabit :exec
UPDATE habits
SET is_deleted = TRUE, deleted_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: ArchiveHabit :exec
UPDATE habits
SET is_archived = TRUE, updated_at = NOW()
WHERE id = $1;

-- name: ReorderHabits :exec
UPDATE habits
SET sort_order = $2, updated_at = NOW()
WHERE id = $1;
