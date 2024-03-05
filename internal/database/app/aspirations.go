package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"context"
	"github.com/google/uuid"
)

func CreateAspiration(aspiration *models.Aspiration) (*models.Aspiration, error) {
	_, err := database.DB.Exec(context.Background(), `
		INSERT INTO aspirations (user_id, subject, message, anonymous, organization_id, closed)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		aspiration.UserID, aspiration.Subject, aspiration.Message, aspiration.Anonymous, aspiration.OrganizationID, aspiration.Closed)

	if err != nil {
		return nil, err
	}

	row := database.DB.QueryRow(context.Background(), `
	SELECT id, created_at FROM aspirations WHERE user_id = $1 AND subject = $2 AND message = $3 AND anonymous = $4 AND organization_id = $5 AND closed = $6`,
		aspiration.UserID, aspiration.Subject, aspiration.Message, aspiration.Anonymous, aspiration.OrganizationID, aspiration.Closed)

	err = row.Scan(&aspiration.ID, &aspiration.CreatedAt)
	if err != nil {
		return nil, err
	}

	return aspiration, nil
}

func CloseAspirationByID(id int) error {
	_, err := database.DB.Exec(context.Background(), `
		UPDATE aspirations SET closed = true WHERE id = $1`, id)
	return err
}

func DeleteAspirationByID(id int) error {
	_, err := database.DB.Exec(context.Background(), `
		DELETE FROM aspirations WHERE id = $1`, id)
	return err
}

func GetAspirations(queryParams map[string]string) ([]models.Aspiration, error) {
	var aspirations []models.Aspiration

	query := `
		SELECT * FROM aspirations
		WHERE 1=1`

	if queryParams["organization_id"] != "" {
		query += " AND organization_id = " + queryParams["organization_id"]
	}

	if queryParams["user_id"] != "" {
		query += " AND user_id = " + queryParams["user_id"]
	}

	if queryParams["closed"] != "" {
		query += " AND closed = " + queryParams["closed"]
	}

	if queryParams["anonymous"] != "" {
		query += " AND anonymous = " + queryParams["anonymous"]
	}

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var aspiration models.Aspiration
		err := rows.Scan(
			&aspiration.ID,
			&aspiration.UserID,
			&aspiration.Subject,
			&aspiration.Message,
			&aspiration.Anonymous,
			&aspiration.OrganizationID,
			&aspiration.Closed,
			&aspiration.CreatedAt,
			&aspiration.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		aspirations = append(aspirations, aspiration)
	}

	return aspirations, nil
}

func GetAspirationByID(id int) (*models.Aspiration, error) {
	var aspiration models.Aspiration

	row := database.DB.QueryRow(context.Background(), `
		SELECT * FROM aspirations WHERE id = $1`, id)

	err := row.Scan(
		&aspiration.ID,
		&aspiration.UserID,
		&aspiration.Subject,
		&aspiration.Message,
		&aspiration.Anonymous,
		&aspiration.OrganizationID,
		&aspiration.CreatedAt,
		&aspiration.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &aspiration, nil
}

func UpvoteExists(userID uuid.UUID, aspirationID int) (bool, error) {
	var exists bool

	row := database.DB.QueryRow(context.Background(), `
		SELECT EXISTS(SELECT 1 FROM aspirations_upvote WHERE user_id = $1 AND aspiration_id = $2)`, userID, aspirationID)

	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func AddUpvote(userID uuid.UUID, aspirationID int) error {
	_, err := database.DB.Exec(context.Background(), `
		INSERT INTO aspirations_upvote (user_id, aspiration_id) VALUES ($1, $2)`, userID, aspirationID)
	return err
}

func RemoveUpvote(userID uuid.UUID, aspirationID int) error {
	_, err := database.DB.Exec(context.Background(), `
		DELETE FROM aspirations_upvote WHERE user_id = $1 AND aspiration_id = $2`, userID, aspirationID)
	return err
}

func GetUpvotesByAspirationID(aspirationID int) (int, error) {
	var upvotes int

	row := database.DB.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM aspirations_upvote WHERE aspiration_id = $1`, aspirationID)

	err := row.Scan(&upvotes)
	if err != nil {
		return 0, err
	}

	return upvotes, nil
}

func AddAdminReply(aspirationID int, reply string) error {
	_, err := database.DB.Exec(context.Background(), `
		UPDATE aspirations SET admin_reply = $1 WHERE id = $2`, reply, aspirationID)
	return err
}
