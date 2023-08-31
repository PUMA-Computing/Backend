package news

import (
	"github.com/gocql/gocql"
	"time"
)

type News struct {
	ID            int64      `json:"id"`
	AuthorID      gocql.UUID `json:"author_id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	Categories    []string   `json:"categories"`
	Thumbnail     string     `json:"thumbnail"`
	Visible       bool       `json:"visible"`
	PublishedDate time.Time  `json:"published_date"`
}

type NewsImage struct {
	ID        int64  `json:"id"`
	NewsID    int64  `json:"news_id"`
	ImageData []byte `json:"image_data"`
}
