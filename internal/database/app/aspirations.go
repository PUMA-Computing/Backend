package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"context"
)

func CreateAspiration(aspiration *models.Aspiration) error {
	_, err := database.DB.Exec(context.Background(), `
		INSERT INTO aspirations (user_id, subject, message, anonymous, organization_id, closed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		aspiration.UserID, aspiration.Subject, aspiration.Message, aspiration.Anonymous, aspiration.OrganizationID, aspiration.Closed)
	return err
}

func CloseAspirationByID(id int) error {
	_, err := database.DB.Exec(context.Background(), `
		UPDATE aspirations SET close = true WHERE id = $1`, id)
	return err
}

func DeleteAspirationByID(id int) error {
	_, err := database.DB.Exec(context.Background(), `
		DELETE FROM aspirations WHERE id = $1`, id)
	return err
}

func GetAspirationsByOrganizationID(organizationID int) ([]models.Aspiration, error) {
	rows, err := database.DB.Query(context.Background(), `
		SELECT * FROM aspirations WHERE organization_id = $1`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aspirations []models.Aspiration
	for rows.Next() {
		var aspiration models.Aspiration
		err := rows.Scan(
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
		aspirations = append(aspirations, aspiration)
	}
	return aspirations, nil
}

func GetAspirationsByUserID(userID int) ([]models.Aspiration, error) {
	rows, err := database.DB.Query(context.Background(), `
		SELECT * FROM aspirations WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aspirations []models.Aspiration
	for rows.Next() {
		var aspiration models.Aspiration
		err := rows.Scan(
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
		aspirations = append(aspirations, aspiration)
	}
	return aspirations, nil
}

func GetAspirationByID(id int) (*models.Aspiration, error) {
	var aspiration models.Aspiration
	err := database.DB.QueryRow(context.Background(), `
		SELECT * FROM aspirations WHERE id = $1`, id).Scan(
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
