package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
	"github.com/google/uuid"
)

type NewsService struct {
}

func NewNewsService() *NewsService {
	return &NewsService{}
}

func (ns *NewsService) CreateNews(news *models.News) error {
	if err := app.CreateNews(news); err != nil {
		return err
	}

	return nil
}

func (ns *NewsService) EditNews(newsID int, updatedNews *models.News) error {
	existingNews, err := app.GetNewsByID(newsID)
	if err != nil {
		return err
	}

	if updatedNews.Title == "" {
		updatedNews.Title = existingNews.Title
	}

	if updatedNews.Content == "" {
		updatedNews.Content = existingNews.Content
	}

	if updatedNews.PublishDate.IsZero() {
		updatedNews.PublishDate = existingNews.PublishDate
	}

	if updatedNews.Likes == 0 {
		updatedNews.Likes = existingNews.Likes
	}

	if updatedNews.UserID == uuid.Nil {
		updatedNews.UserID = existingNews.UserID
	}

	if updatedNews.CreatedAt.IsZero() {
		updatedNews.CreatedAt = existingNews.CreatedAt
	}

	if updatedNews.UpdatedAt.IsZero() {
		updatedNews.UpdatedAt = existingNews.UpdatedAt
	}

	if err := app.UpdateNews(newsID, updatedNews); err != nil {
		return err
	}

	return nil
}

func (ns *NewsService) DeleteNews(newsID int) error {
	if err := app.DeleteNews(newsID); err != nil {
		return err
	}

	return nil
}

func (ns *NewsService) GetNewsByID(newsID int) (*models.News, error) {
	news, err := app.GetNewsByID(newsID)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (ns *NewsService) ListNews() ([]*models.News, error) {
	news, err := app.ListNews()
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (ns *NewsService) LikeNews(newsID int) error {
	if err := app.LikeNews(newsID); err != nil {
		return err
	}

	return nil
}
