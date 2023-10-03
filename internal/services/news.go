package services

import (
	"Backend/internal/database"
	"Backend/internal/models"
)

type NewsService struct {
}

func NewNewsService() *NewsService {
	return &NewsService{}
}

func (ns *NewsService) CreateNews(news *models.News) error {
	if err := database.CreateNews(news); err != nil {
		return err
	}

	return nil
}

func (ns *NewsService) EditNews(newsID int, updatedNews *models.News) error {
	if err := database.UpdateNews(newsID, updatedNews); err != nil {
		return err
	}

	return nil
}

func (ns *NewsService) DeleteNews(newsID int) error {
	if err := database.DeleteNews(newsID); err != nil {
		return err
	}

	return nil
}

func (ns *NewsService) GetNewsByID(newsID int) (*models.News, error) {
	news, err := database.GetNewsByID(newsID)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (ns *NewsService) ListNews() ([]*models.News, error) {
	news, err := database.ListNews()
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (ns *NewsService) LikeNews(newsID int) error {
	if err := database.LikeNews(newsID); err != nil {
		return err
	}

	return nil
}
