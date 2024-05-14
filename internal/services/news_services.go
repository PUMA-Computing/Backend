package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
	"Backend/pkg/utils"
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

	utils.ReflectiveUpdate(existingNews, updatedNews)

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

func (ns *NewsService) ListNews(queryParams map[string]string) ([]*models.News, int, error) {
	news, totalRecords, err := app.ListNews(queryParams)
	if err != nil {
		return nil, 0, err
	}
	return news, totalRecords, nil
}

func (ns *NewsService) LikeNews(userID uuid.UUID, newsID int) error {
	if err := app.LikeNews(userID, newsID); err != nil {
		return err
	}
	return nil
}
