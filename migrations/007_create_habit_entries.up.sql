CREATE TABLE habit_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    habit_id UUID NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    value INTEGER NOT NULL DEFAULT 1,
    note TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (habit_id, date)
);

CREATE INDEX idx_habit_entries_habit_date
    ON habit_entries (habit_id, date DESC);

CREATE INDEX idx_habit_entries_user_date
    ON habit_entries (user_id, date DESC);
