package models

import "time"

type News struct {
	ID          int
	Title       string
	Content     string
	UserID      int
	PublishDate time.Time
	Likes       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
