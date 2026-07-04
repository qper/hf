CREATE OR REPLACE FUNCTION calculate_streak_for_habit(p_habit_id UUID, p_date DATE)
RETURNS INT AS $$
DECLARE
    streak INT := 0;
    current_date DATE := p_date;
BEGIN
    WHILE EXISTS (
        SELECT 1
        FROM habit_entries
        WHERE habit_id = p_habit_id
          AND entry_date = current_date
          AND completed = TRUE
    ) LOOP
        streak := streak + 1;
        current_date := current_date - INTERVAL '1 day';
    END LOOP;

    RETURN streak;
END;
$$ LANGUAGE plpgsql;
