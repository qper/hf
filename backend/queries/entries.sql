-- name: UpsertEntry :one
INSERT INTO habit_entries (habit_id, entry_date, completed, note)
VALUES ($1, $2, $3, $4)
ON CONFLICT (habit_id, entry_date) DO UPDATE SET
    completed = EXCLUDED.completed,
    note = EXCLUDED.note,
    updated_at = NOW()
RETURNING *;

-- name: GetEntryByHabitDate :one
SELECT *
FROM habit_entries
WHERE habit_id = $1 AND entry_date = $2
LIMIT 1;

-- name: GetEntriesForDateRange :many
SELECT he.*
FROM habit_entries he
JOIN habits h ON h.id = he.habit_id
WHERE h.user_id = $1
  AND he.entry_date >= $2
  AND he.entry_date <= $3
ORDER BY he.entry_date DESC;
