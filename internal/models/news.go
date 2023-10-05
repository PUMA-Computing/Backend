package models

import (
	"github.com/google/uuid"
	"time"
)

type News struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	UserID      uuid.UUID `json:"user_id"`
	PublishDate time.Time `json:"publish_date"`
	Likes       int       `json:"likes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
