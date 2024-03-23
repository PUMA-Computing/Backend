package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"context"
	"github.com/google/uuid"
)

func CreateNews(news *models.News) error {
	_, err := database.DB.Exec(context.Background(), `
		INSERT INTO news (title, content, user_id, publish_date)
		VALUES ($1, $2, $3, $4)`,
		news.Title, news.Content, news.UserID, news.PublishDate)
	return err
}

func UpdateNews(newsID int, news *models.News) error {
	_, err := database.DB.Exec(context.Background(), `
		UPDATE news SET title = $1, content = $2, publish_date = $3, updated_at = $4
		WHERE id = $5`, news.Title, news.Content, news.PublishDate, news.UpdatedAt, newsID)
	return err
}

func DeleteNews(newsID int) error {
	_, err := database.DB.Exec(context.Background(), `
		DELETE FROM news WHERE id = $1`, newsID)
	return err
}

func GetNewsByID(newsID int) (*models.News, error) {
	var news models.News
	err := database.DB.QueryRow(context.Background(), `
		SELECT id, title, content, user_id, publish_date, likes, created_at, updated_at
		FROM news WHERE id = $1`, newsID).Scan(&news.ID, &news.Title, &news.Content, &news.UserID, &news.PublishDate, &news.Likes, &news.CreatedAt, &news.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &news, nil
}

func ListNews() ([]*models.News, error) {
	var news []*models.News
	rows, err := database.DB.Query(context.Background(), `
		SELECT id, title, content, user_id, publish_date, created_at, updated_at
		FROM news`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var n models.News
		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.UserID, &n.PublishDate, &n.Likes, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			return nil, err
		}
		news = append(news, &n)
	}
	return news, nil
}

func LikeNews(userID uuid.UUID, newsID int) error {
	_, err := database.DB.Exec(context.Background(), `
		INSERT INTO news_likes (user_id, news_id) VALUES ($1, $2)`, userID, newsID)
	return err
}

func UnlikeNews(userID uuid.UUID, newsID int) error {
	_, err := database.DB.Exec(context.Background(), `
		DELETE FROM news_likes WHERE user_id = $1 AND news_id = $2`, userID, newsID)
	return err
}
