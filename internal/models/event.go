package models

import "time"

type Event struct {
	ID          int
	Title       string
	Description string
	Date        time.Time
	UserID      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type EventRegistration struct {
	ID               int
	EventID          int
	UserID           int
	RegistrationDate time.Time
}
