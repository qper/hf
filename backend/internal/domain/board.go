package domain

type BoardProgress struct {
	Done  int `json:"done"`
	Total int `json:"total"`
}

type BoardHabit struct {
	ID          string         `json:"id"`
	UserID      string         `json:"user_id"`
	CategoryID  *string        `json:"category_id,omitempty"`
	Name        string         `json:"name"`
	Description *string        `json:"description,omitempty"`
	Color       *string        `json:"color,omitempty"`
	Type        HabitType      `json:"type"`
	Frequency   HabitFrequency `json:"frequency"`
	TargetValue *float64       `json:"target_value,omitempty"`
	Unit        *string        `json:"unit,omitempty"`
	SortOrder   int            `json:"sort_order"`
	IsCompleted bool           `json:"is_completed"`
	Streak      int            `json:"streak"`
}

type Board struct {
	Date       string        `json:"date"`
	IsEditable bool          `json:"is_editable"`
	Progress   BoardProgress `json:"progress"`
	Habits     []BoardHabit  `json:"habits"`
}
