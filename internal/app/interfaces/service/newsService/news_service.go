package newsService

import (
	"Backend/internal/app/domain/news"
	"Backend/internal/app/interfaces/repository/newsRepository"
)

type NewsService interface {
	GetNews() ([]*news.News, error)
	GetNewsByID(id int64) (*news.News, error)
	CreateNews(news *news.News) error
	UpdateNews(newsID int64, news *news.News) error
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

func (s *NewsServiceImpl) CreateNews(news *news.News) error {
	return s.newsRepository.CreateNews(news)
}

func (s *NewsServiceImpl) UpdateNews(newsID int64, UpdatedNews *news.News) error {
	existingNews, err := s.newsRepository.GetNewsByID(newsID)
	if err != nil {
		return err
	}

	existingNews.Title = UpdatedNews.Title
	existingNews.Content = UpdatedNews.Content
	existingNews.Categories = UpdatedNews.Categories
	existingNews.Thumbnail = UpdatedNews.Thumbnail
	existingNews.Visible = UpdatedNews.Visible
	existingNews.PublishedDate = UpdatedNews.PublishedDate

	return s.newsRepository.UpdateNews(existingNews)
}

func (s *NewsServiceImpl) DeleteNews(id int64) error {
	return s.newsRepository.DeleteNews(id)
}
