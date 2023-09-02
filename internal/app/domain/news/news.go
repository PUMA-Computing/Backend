package news

import (
	"github.com/google/uuid"
	"time"
)

type News struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement:true"`
	AuthorID      uuid.UUID `json:"author_id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	CategoryID    int       `json:"category_id"`
	Thumbnail     string    `json:"thumbnail"`
	Status        string    `json:"status"`
	PublishedDate time.Time `json:"published_date"`
	CreatedAt     time.Time `json:"created_at"`
}

type Categories struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
