package models

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	UserID         uuid.UUID `json:"user_id"`
	Status         string    `json:"status"`
	Link           string    `json:"link"`
	Thumbnail      string    `json:"thumbnail"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt"`
	OrganizationID int       `json:"organization_id"`
}

type EventRegistration struct {
	ID               int       `json:"id"`
	EventID          int       `json:"event_id"`
	UserID           int       `json:"user_id"`
	RegistrationDate time.Time `json:"registration_date"`
}
