DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'habit_type') THEN
        CREATE TYPE habit_type AS ENUM ('good', 'bad');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'habit_freq') THEN
        CREATE TYPE habit_freq AS ENUM ('daily', 'weekly', 'monthly');
    END IF;
END
$$;

CREATE TABLE habits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NULL REFERENCES categories(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT NULL,
    habit_type habit_type NOT NULL,
    habit_freq habit_freq NOT NULL,
    target_value INTEGER NOT NULL DEFAULT 1,
    unit VARCHAR(50) NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    start_date DATE NOT NULL DEFAULT CURRENT_DATE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_habits_user_id
    ON habits (user_id);

CREATE INDEX idx_habits_category_id
    ON habits (category_id);

CREATE INDEX idx_habits_sort_order
    ON habits (sort_order);

CREATE INDEX idx_habits_name_trgm
    ON habits USING GIN (name gin_trgm_ops);
