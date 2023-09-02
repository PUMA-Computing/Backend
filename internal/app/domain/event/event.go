package event

import (
	"github.com/google/uuid"
	"time"
)

type Events struct {
	ID          int64     `json:"id"`
	Image       string    `json:"image"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}

type Registration struct {
	ID        int64     `json:"id"`
	EventID   int64     `json:"event_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
