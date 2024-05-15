package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"strconv"
)

func CreateNews(news *models.News) error {
	_, err := database.DB.Exec(context.Background(), `
		INSERT INTO news (title, content, user_id, publish_date, likes, created_at, updated_at, thumbnail, slug, organization_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, news.Title, news.Content, news.UserID, news.PublishDate, news.Likes, news.CreatedAt, news.UpdatedAt, news.Thumbnail, news.Slug)
	return err
}

func UpdateNews(newsID int, news *models.News) error {
	_, err := database.DB.Exec(context.Background(), `
		UPDATE news SET title = $1, content = $2, publish_date = $3, updated_at = $4, thumbnail = $5, slug = $6, organization_id = $7
		WHERE id = $5`, news.Title, news.Content, news.PublishDate, news.UpdatedAt, news.Thumbnail, news.Slug, news.OrganizationID, newsID)
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
		SELECT id, title, content, user_id, publish_date, likes, created_at, updated_at, thumbnail, slug, organization_id
		FROM news WHERE id = $1`, newsID).Scan(&news.ID, &news.Title, &news.Content, &news.UserID, &news.PublishDate, &news.Likes, &news.CreatedAt, &news.UpdatedAt, &news.Thumbnail, &news.Slug, &news.OrganizationID)
	if err != nil {
		return nil, err
	}
	return &news, nil
}

func GetNewsBySlug(slug string) (*models.News, error) {
	var news models.News
	err := database.DB.QueryRow(context.Background(), `
		SELECT n.id, n.title, n.content, n.user_id, n.publish_date, n.likes, n.created_at, n.updated_at, n.thumbnail, n.slug, n.organization_id, o.name as organizations, CONCAT(u.first_name, ' ', u.last_name) AS author
		FROM news n
		LEFT JOIN organizations o ON n.organization_id = o.id
		LEFT JOIN users u ON n.user_id = u.id
		WHERE n.slug = $1`, slug).Scan(&news.ID, &news.Title, &news.Content, &news.UserID, &news.PublishDate, &news.Likes, &news.CreatedAt, &news.UpdatedAt, &news.Thumbnail, &news.Slug, &news.OrganizationID, &news.Organization, &news.Author)

	if err != nil {
		return nil, err
	}
	return &news, nil
}

// ListNews returns a list of news based on the query parameters
func ListNews(queryParams map[string]string) ([]*models.News, int, error) {
	limit := 10

	query := `
		SELECT n.id, n.title, n.content, n.user_id, n.publish_date, n.likes, n.created_at, n.updated_at, n.thumbnail, n.slug, n.organization_id, o.name as organizations, CONCAT(u.first_name, ' ', u.last_name) AS author
		FROM news n
		LEFT JOIN organizations o ON n.organization_id = o.id
		LEFT JOIN users u ON n.user_id = u.id
		WHERE 1 = 1`

	if queryParams["organization_id"] != "" {
		query += " AND n.organization_id = " + queryParams["organization_id"]
	}

	var totalRecords int
	err := database.DB.QueryRow(context.Background(), ` SELECT COUNT(*) FROM news n WHERE 1 = 1`).Scan(&totalRecords)
	if err != nil {
		return nil, 0, err
	}

	totalPages := (totalRecords + limit - 1) / limit
	if queryParams["page"] != "" {
		page, err := strconv.Atoi(queryParams["page"])
		if err != nil {
			return nil, totalPages, err
		}
		offset := (page - 1) * limit
		query += fmt.Sprintf(" ORDER BY n.created_at DESC LIMIT %d OFFSET %d", limit, offset)
	} else {
		query += fmt.Sprintf(" ORDER BY n.created_at DESC LIMIT %d", limit)
	}

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return nil, totalPages, err
	}

	defer rows.Close()

	var news []*models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.UserID, &n.PublishDate, &n.Likes, &n.CreatedAt, &n.UpdatedAt, &n.Thumbnail, &n.Slug, &n.OrganizationID, &n.Organization, &n.Author)
		if err != nil {
			return nil, totalPages, err
		}
		news = append(news, &n)
	}
	return news, totalPages, nil
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
