package domain

type Streak struct {
	HabitID string `json:"habit_id"`
	Current int    `json:"current"`
	Max     int    `json:"max"`
}
