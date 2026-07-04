-- name: GetCompletionRate :one
SELECT COALESCE(AVG(CASE WHEN he.completed = TRUE THEN 1.0 ELSE 0.0 END), 0.0) AS completion_rate
FROM habit_entries he
JOIN habits h ON h.id = he.habit_id
WHERE h.user_id = $1
  AND he.entry_date >= $2
  AND he.entry_date <= $3;

-- name: GetHeatmapData :many
SELECT he.entry_date AS entry_date, COUNT(*) FILTER (WHERE he.completed = TRUE) AS total_value
FROM habit_entries he
JOIN habits h ON h.id = he.habit_id
WHERE h.user_id = $1
  AND he.entry_date >= $2
  AND he.entry_date <= $3
GROUP BY he.entry_date
ORDER BY he.entry_date;

-- name: GetCurrentStreak :one
SELECT COALESCE(calculate_streak_for_habit($1, $2), 0) AS current_streak;

-- name: GetMaxStreak :one
SELECT 0 AS max_streak;
