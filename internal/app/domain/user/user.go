package user

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	RoleID    int       `json:"role_id"`
	NIM       string    `json:"nim"`
	Year      string    `json:"year"`
	CreatedAt time.Time `json:"created_at"`
}

type SessionData struct {
	UserID       uuid.UUID `json:"user_id" gorm:"PRIMARY_KEY"`
	SessionToken string    `json:"session_token"gorm:"unique"`
	ExpiredAt    time.Time `json:"expired_at"`
}
