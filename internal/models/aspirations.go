package models

import "github.com/google/uuid"

type Aspiration struct {
	ID             int       `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Subject        string    `json:"subject"`
	Message        string    `json:"message"`
	Anonymous      bool      `json:"anonymous"`
	OrganizationID int       `json:"organization_id"`
	Closed         bool      `json:"close"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}
