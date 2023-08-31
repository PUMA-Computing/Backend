package event

import "time"

type Event struct {
	ID              string    `json:"id"`
	Image           string    `json:"image"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Date            time.Time `json:"date"`
	RegisteredUsers []string  `json:"registered_users"`
}
