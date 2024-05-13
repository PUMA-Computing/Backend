package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"log"
)

func CreateUser(user *models.User) error {

	query := `
		INSERT INTO users (id, username, password, first_name, middle_name, last_name, email, student_id, major, year, role_id, email_verification_token)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := database.DB.Exec(
		context.Background(),
		query,
		user.ID, user.Username, user.Password, user.FirstName, user.MiddleName, user.LastName, user.Email, user.StudentID, user.Major, user.Year, user.RoleID, user.EmailVerificationToken,
	)
	if err != nil {
		log.Printf("Error during query execution or scanning: %v", err)
		return err
	}
	return nil
}

func AuthenticateUser(usernameOrEmail string) (*models.User, error) {
	var user models.User
	var userID string
	var query string
	var err error
	var username *string
	var middleName sql.NullString

	if usernameOrEmail != "" {
		query = `
			SELECT id, username, password, first_name, middle_name, last_name, email, student_id, major, year, role_id, created_at, updated_at, email_verified, student_id_verified
			FROM users
			WHERE username = $1 OR email = $1`

		var middleNamePtr *string

		err = database.DB.QueryRow(
			context.Background(),
			query,
			usernameOrEmail,
		).Scan(
			&userID, &username, &user.Password, &user.FirstName, &middleNamePtr, &user.LastName, &user.Email, &user.StudentID, &user.Major, &user.Year, &user.RoleID, &user.CreatedAt, &user.UpdatedAt, &user.EmailVerified, &user.StudentIDVerified,
		)
		if username != nil {
			user.Username = *username
		}
		if middleName != (sql.NullString{}) {
			user.MiddleName = *middleNamePtr
		}
		if err != nil {
			log.Println(err)
			return nil, err
		}
	} else {
		return nil, nil
	}
	user.ID, err = uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	user.MiddleName = middleName.String
	return &user, nil
}

func VerifyEmail(token string) error {
	query := `
		UPDATE users
		SET email_verified = true
		WHERE email_verification_token = $1`
	_, err := database.DB.Exec(
		context.Background(),
		query,
		token,
	)
	if err != nil {
		log.Printf("Error during query execution or scanning: %v", err)
		return err
	}
	return nil
}

func PasswordResetToken(token string) error {
	query := `
		UPDATE users
		SET password_reset_token = $1
		WHERE email = $2`
	_, err := database.DB.Exec(
		context.Background(),
		query,
		token,
	)
	if err != nil {
		log.Printf("Error during query execution or scanning: %v", err)
		return err
	}
	return nil
}
