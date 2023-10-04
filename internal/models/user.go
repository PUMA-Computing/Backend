package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID         uuid.UUID `pg:"type:uuid" json:"id"`
	Username   string    `json:"username" binding:"required"`
	Password   string    `json:"password" binding:"required"`
	FirstName  string    `json:"first_name" binding:"required"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name" binding:"required"`
	Email      string    `json:"email" binding:"required"`
	StudentID  string    `json:"student_id" binding:"required"`
	Major      string    `json:"major" binding:"required"`
	RoleID     int       `json:"role_id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
