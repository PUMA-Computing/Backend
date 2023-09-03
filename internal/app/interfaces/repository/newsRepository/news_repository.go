package newsRepository

import (
	"Backend/internal/app/domain/news"
	"gorm.io/gorm"
)

type NewsRepository interface {
	GetNews() ([]*news.News, error)
	GetNewsByID(id int64) (*news.News, error)
	GetNewsByCategory(category string) ([]*news.News, error)
	GetNewsByStatus(status string) ([]*news.News, error)
	CreateNews(news *news.News) error
	UpdateNews(news *news.News) error
	DeleteNews(id int64) error
}

type PostgresForNewsRepository struct {
	DB *gorm.DB
}

func NewPostgresForNewsRepository(DB *gorm.DB) *PostgresForNewsRepository {
	return &PostgresForNewsRepository{DB: DB}
}

func (r *PostgresForNewsRepository) GetNews() ([]*news.News, error) {
	var newsList []*news.News
	if err := r.DB.Find(&newsList).Error; err != nil {
		return nil, err
	}
	return newsList, nil
}

func (r *PostgresForNewsRepository) GetNewsByID(id int64) (*news.News, error) {
	var newsTab news.News
	if err := r.DB.Where("id = ?", id).First(&newsTab).Error; err != nil {
		return nil, err
	}
	return &newsTab, nil
}

func (r *PostgresForNewsRepository) GetNewsByCategory(category string) ([]*news.News, error) {
	var newsList []*news.News
	if err := r.DB.Where("category = ?", category).Find(&newsList).Error; err != nil {
		return nil, err
	}
	return newsList, nil
}

func (r *PostgresForNewsRepository) GetNewsByStatus(status string) ([]*news.News, error) {
	var newsList []*news.News
	if err := r.DB.Where("status = ?", status).Find(&newsList).Error; err != nil {
		return nil, err
	}
	return newsList, nil
}

func (r *PostgresForNewsRepository) GetNewsByAuthorID(authorID int64) ([]*news.News, error) {
	var newsList []*news.News
	if err := r.DB.Where("author_id = ?", authorID).Find(&newsList).Error; err != nil {
		return nil, err
	}
	return newsList, nil
}

func (r *PostgresForNewsRepository) CreateNews(news *news.News) error {
	return r.DB.Create(news).Error
}

func (r *PostgresForNewsRepository) UpdateNews(news *news.News) error {
	return r.DB.Save(news).Error
}

func (r *PostgresForNewsRepository) DeleteNews(id int64) error {
	return r.DB.Delete(&news.News{}, id).Error
}
