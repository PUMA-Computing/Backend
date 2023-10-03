package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID         uuid.UUID
	Username   string
	Password   string
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	StudentID  string
	Major      string
	RoleID     int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
