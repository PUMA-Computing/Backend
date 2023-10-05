package models

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type EventRegistration struct {
	ID               int       `json:"id"`
	EventID          int       `json:"event_id"`
	UserID           int       `json:"user_id"`
	RegistrationDate time.Time `json:"registration_date"`
}
