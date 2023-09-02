package newsService

import (
	"Backend/internal/app/domain/news"
	"Backend/internal/app/interfaces/repository/newsRepository"
)

type NewsService interface {
	GetNews() ([]*news.News, error)
	GetNewsByID(id int64) (*news.News, error)
	CreateNews(news *news.News) error
	UpdateNews(news *news.News) error
	DeleteNews(id int64) error
}

type NewsServiceImpl struct {
	newsRepository newsRepository.NewsRepository
}

func NewNewsService(NewsRepository newsRepository.NewsRepository) *NewsServiceImpl {
	return &NewsServiceImpl{newsRepository: NewsRepository}
}

func (s *NewsServiceImpl) GetNews() ([]*news.News, error) {
	newsGet, err := s.newsRepository.GetNews()
	if err != nil {
		return nil, err
	}
	return newsGet, nil
}

func (s *NewsServiceImpl) GetNewsByID(id int64) (*news.News, error) {
	newsGetID, err := s.newsRepository.GetNewsByID(id)
	if err != nil {
		return nil, err
	}
	return newsGetID, nil
}

func (r *NewsServiceImpl) GetNewsByCategory(category string) ([]*news.News, error) {
	newsGetCategory, err := r.newsRepository.GetNewsByCategory(category)
	if err != nil {
		return nil, err
	}
	return newsGetCategory, nil
}

func (r *NewsServiceImpl) GetNewsByStatus(status string) ([]*news.News, error) {
	newsGetStatus, err := r.newsRepository.GetNewsByStatus(status)
	if err != nil {
		return nil, err
	}
	return newsGetStatus, nil
}

func (s *NewsServiceImpl) CreateNews(news *news.News) error {
	return s.newsRepository.CreateNews(news)
}

func (s *NewsServiceImpl) UpdateNews(news *news.News) error {
	existingNews, err := s.newsRepository.GetNewsByID(news.ID)
	if err != nil {
		return err
	}

	existingNews.Title = news.Title
	existingNews.Content = news.Content
	existingNews.CategoryID = news.CategoryID
	existingNews.Thumbnail = news.Thumbnail
	existingNews.Status = news.Status
	existingNews.PublishedDate = news.PublishedDate

	return s.newsRepository.UpdateNews(existingNews)
}

func (s *NewsServiceImpl) DeleteNews(id int64) error {
	return s.newsRepository.DeleteNews(id)
}
