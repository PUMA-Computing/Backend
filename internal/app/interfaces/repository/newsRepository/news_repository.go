package newsRepository

import (
	"Backend/internal/app/domain/news"
	"github.com/gocql/gocql"
	"time"
)

type NewsRepository interface {
	CreateNews(news *news.News) error
	GetNews() ([]*news.News, error)
	GetNewsByID(id int64) (*news.News, error)
	UpdateNews(news *news.News) error
	DeleteNews(id int64) error
}

type CassandraForNewsRepository struct {
	session *gocql.Session
}

func NewCassandraForNewsRepository(session *gocql.Session) *CassandraForNewsRepository {
	return &CassandraForNewsRepository{session: session}
}

func (r *CassandraForNewsRepository) GetNews() ([]*news.News, error) {
	query := "SELECT id, author_id, title, content, categories, thumbnail, visible, published_date FROM news"
	iter := r.session.Query(query).Iter()

	var newsList []*news.News
	var id int64
	var authorID gocql.UUID
	var title string
	var content string
	var categories []string
	var thumbnail string
	var visible bool
	var publishedDate time.Time

	for iter.Scan(&id, &authorID, &title, &content, &categories, &thumbnail, &visible, &publishedDate) {
		newsTab := &news.News{
			ID:            id,
			AuthorID:      authorID,
			Title:         title,
			Content:       content,
			Categories:    categories,
			Thumbnail:     thumbnail,
			Visible:       visible,
			PublishedDate: publishedDate,
		}
		newsList = append(newsList, newsTab)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return newsList, nil
}

func (r *CassandraForNewsRepository) GetNewsByID(id int64) (*news.News, error) {
	var newsTab news.News
	query := "SELECT id, author_id, title, content, categories, thumbnail, visible, published_date FROM news WHERE id = ?"
	if err := r.session.Query(query, id).Scan(
		&newsTab.ID, &newsTab.AuthorID, &newsTab.Title, &newsTab.Content, &newsTab.Categories, &newsTab.Thumbnail, &newsTab.Visible, &newsTab.PublishedDate,
	); err != nil {
		return nil, err
	}
	return &newsTab, nil
}

func (r *CassandraForNewsRepository) CreateNews(news *news.News) error {
	query := "INSERT INTO news (id, author_id, title, content, categories, thumbnail, visible, published_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	if err := r.session.Query(query, news.ID, news.AuthorID, news.Title, news.Content, news.Categories, news.Thumbnail, news.Visible, news.PublishedDate).Exec(); err != nil {
		return err
	}
	return nil
}

func (r *CassandraForNewsRepository) UpdateNews(news *news.News) error {
	query := "UPDATE news SET title = ?, content = ?, categories = ?, thumbnail = ?, visible = ?, published_date = ? WHERE id = ?"
	if err := r.session.Query(query, news.Title, news.Content, news.Categories, news.Thumbnail, news.Visible, news.PublishedDate, news.ID).Exec(); err != nil {
		return err
	}
	return nil
}

func (r *CassandraForNewsRepository) DeleteNews(id int64) error {
	query := "DELETE FROM news WHERE id = ?"
	if err := r.session.Query(query, id).Exec(); err != nil {
		return err
	}
	return nil
}
