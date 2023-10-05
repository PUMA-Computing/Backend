package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID         uuid.UUID `pg:"type:uuid" json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	StudentID  string    `json:"student_id"`
	Major      string    `json:"major"`
	RoleID     int       `json:"role_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
