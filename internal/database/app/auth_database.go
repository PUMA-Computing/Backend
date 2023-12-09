package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"context"
	"github.com/google/uuid"
	"log"
)

func TableHasRows(tableName string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM " + tableName + ")"
	var exists bool
	err := database.DB.QueryRow(context.Background(), query).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, username, password, first_name, middle_name, last_name, email, student_id, major, role_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	log.Printf("SQL Query: %v", query)

	_, err := database.DB.Exec(context.Background(), query, user.ID, user.Username, user.Password, user.FirstName, user.MiddleName, user.LastName, user.Email, user.StudentID, user.Major, user.RoleID)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}
	return err
}

func AuthenticateUser(username, email string, password string) (*models.User, error) {
	var user models.User
	var userID string
	var query string
	var err error

	if username != "" {
		query = "SELECT * FROM users WHERE username = $1"
		err = database.DB.QueryRow(context.Background(), query, username).Scan(&userID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.StudentID, &user.Major, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)
	} else if email != "" {
		query = "SELECT * FROM users WHERE email = $1"
		err = database.DB.QueryRow(context.Background(), query, email).Scan(&userID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.StudentID, &user.Major, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)
	} else {
		return nil, nil
	}
	user.ID, err = uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
