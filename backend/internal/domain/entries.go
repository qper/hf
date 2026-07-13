package domain

import "time"

type Entry struct {
	ID        string    `json:"id"`
	HabitID   string    `json:"habit_id"`
	Date      string    `json:"date"`
	Completed bool      `json:"completed"`
	Value     *float64  `json:"value,omitempty"`
	Note      *string   `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateEntryRequest struct {
	HabitID   string   `json:"habit_id"`
	Date      string   `json:"date"`
	Completed *bool    `json:"completed,omitempty"`
	Value     *float64 `json:"value,omitempty"`
	Note      *string  `json:"note,omitempty"`
}

type UpdateEntryRequest struct {
	Value *float64 `json:"value,omitempty"`
	Note  *string  `json:"note,omitempty"`
}
