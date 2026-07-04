DROP INDEX IF EXISTS idx_habits_name_trgm;
DROP INDEX IF EXISTS idx_habits_sort_order;
DROP INDEX IF EXISTS idx_habits_category_id;
DROP INDEX IF EXISTS idx_habits_user_id;
DROP TABLE IF EXISTS habits;
DROP TYPE IF EXISTS habit_freq;
DROP TYPE IF EXISTS habit_type;
