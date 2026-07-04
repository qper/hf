CREATE OR REPLACE FUNCTION current_streak(p_habit_id UUID, p_today DATE)
RETURNS INTEGER
LANGUAGE plpgsql
AS $$
DECLARE
    v_streak INTEGER := 0;
    v_date DATE;
    v_exists BOOLEAN;
BEGIN
    IF p_habit_id IS NULL THEN
        RETURN 0;
    END IF;

    FOR v_date IN SELECT date
                  FROM habit_entries
                  WHERE habit_id = p_habit_id
                    AND date <= p_today
                  ORDER BY date DESC
    LOOP
        SELECT EXISTS (
            SELECT 1
            FROM habit_entries
            WHERE habit_id = p_habit_id
              AND date = v_date
        ) INTO v_exists;

        IF v_exists THEN
            v_streak := v_streak + 1;
        ELSE
            EXIT;
        END IF;

        IF v_date < p_today THEN
            EXIT;
        END IF;
    END LOOP;

    RETURN v_streak;
END;
$$;
