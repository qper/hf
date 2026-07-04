-- name: UpsertEntry :one
INSERT INTO habit_entries (user_id, habit_id, date, value, note)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (habit_id, date) DO UPDATE SET
    value = EXCLUDED.value,
    note = EXCLUDED.note,
    updated_at = NOW()
RETURNING *;

-- name: GetEntryByHabitDate :one
SELECT *
FROM habit_entries
WHERE habit_id = $1 AND date = $2
LIMIT 1;

-- name: GetEntriesForDateRange :many
SELECT he.*
FROM habit_entries he
WHERE he.user_id = $1
  AND he.date >= $2
  AND he.date <= $3
ORDER BY he.date DESC;
