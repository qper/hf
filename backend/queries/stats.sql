-- name: GetCompletionRate :one
SELECT COALESCE(AVG(CASE WHEN value >= target_value THEN 1.0 ELSE 0.0 END), 0.0) AS completion_rate
FROM habit_entries he
JOIN habits h ON h.id = he.habit_id
WHERE h.user_id = $1
  AND he.date >= $2
  AND he.date <= $3;

-- name: GetHeatmapData :many
SELECT he.date, SUM(he.value) AS total_value
FROM habit_entries he
WHERE he.user_id = $1
  AND he.date >= $2
  AND he.date <= $3
GROUP BY he.date
ORDER BY he.date;

-- name: GetCurrentStreak :one
SELECT COALESCE(current_streak($1, $2), 0) AS current_streak;

-- name: GetMaxStreak :one
SELECT 0 AS max_streak;
