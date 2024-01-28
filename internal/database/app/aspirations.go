package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"context"
)

func CreateAspiration(aspiration *models.Aspiration) error {
	_, err := database.DB.Exec(context.Background(), `
		INSERT INTO aspirations (user_id, subject, message, anonymous, organization_id, closed)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		aspiration.UserID, aspiration.Subject, aspiration.Message, aspiration.Anonymous, aspiration.OrganizationID, aspiration.Closed)
	return err
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
