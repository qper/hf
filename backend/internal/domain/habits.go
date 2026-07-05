package domain

import "time"

type HabitType string

type HabitFrequency string

const (
	HabitTypeBoolean  HabitType = "boolean"
	HabitTypeNumeric  HabitType = "numeric"
	HabitTypeDuration HabitType = "duration"
)

const (
	HabitFrequencyDaily  HabitFrequency = "daily"
	HabitFrequencyWeekly HabitFrequency = "weekly"
	HabitFrequencyCustom HabitFrequency = "custom"
)

type Habit struct {
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
	IsArchived  bool           `json:"is_archived"`
	IsDeleted   bool           `json:"is_deleted"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty"`
}

type CreateHabitRequest struct {
	CategoryID  *string  `json:"category_id,omitempty"`
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Color       *string  `json:"color,omitempty"`
	Type        string   `json:"type"`
	Frequency   string   `json:"frequency"`
	TargetValue *float64 `json:"target_value,omitempty"`
	Unit        *string  `json:"unit,omitempty"`
	SortOrder   int      `json:"sort_order,omitempty"`
}

type UpdateHabitRequest struct {
	CategoryID  *string  `json:"category_id,omitempty"`
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Color       *string  `json:"color,omitempty"`
	Type        string   `json:"type,omitempty"`
	Frequency   *string  `json:"frequency,omitempty"`
	TargetValue *float64 `json:"target_value,omitempty"`
	Unit        *string  `json:"unit,omitempty"`
	SortOrder   *int     `json:"sort_order,omitempty"`
	Archived    *bool    `json:"archived,omitempty"`
}

type ArchiveHabitRequest struct {
	Archived bool `json:"archived"`
}

type ReorderHabitsRequest struct {
	IDs []string `json:"ids"`
}
