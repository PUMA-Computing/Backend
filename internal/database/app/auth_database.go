package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"log"
)

func CreateUser(user *models.User) error {

	query := `
		INSERT INTO users (id, username, password, first_name, middle_name, last_name, email, student_id, major, year, role_id, email_verification_token, institution_name, gender)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	_, err := database.DB.Exec(
		context.Background(),
		query,
		user.ID, user.Username, user.Password, user.FirstName, user.MiddleName, user.LastName, user.Email,
		user.StudentID, user.Major, user.Year, user.RoleID, user.EmailVerificationToken, user.InstitutionName, user.Gender,
	)
	if err != nil {
		log.Printf("Error during query execution or scanning: %v", err)
		return err
	}
	return nil
}

func AuthenticateUser(usernameOrEmail string) (*models.User, error) {
	log.Printf("start AuthenticateUser")
	var user models.User
	var userID string
	var query string
	var err error

	log.Println("before query")

	query = `
		SELECT id, username, password, first_name, middle_name, last_name, email, student_id, major, year, role_id, email_verification_token, institution_name, gender
		FROM users
		WHERE username = $1 OR email = $1`

	log.Println("after query")

	err = database.DB.QueryRow(
		context.Background(),
		query,
		usernameOrEmail,
	).Scan(
		&userID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email,
		&user.StudentID, &user.Major, &user.Year, &user.RoleID, &user.EmailVerificationToken, &user.InstitutionName,
		&user.Gender,
	)

	log.Println("after scan")

	if errors.Is(err, sql.ErrNoRows) {
		log.Println("No user found with username or email:", usernameOrEmail)
		return nil, err
	} else if err != nil {
		log.Println("Error during query execution or scanning:", err)
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	user.ID, err = uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func IsEmailVerified(email string) (bool, error) {
	var verified bool
	query := `
		SELECT email_verified
		FROM users
		WHERE email = $1 OR username = $1`
	err := database.DB.QueryRow(
		context.Background(),
		query,
		email,
	).Scan(&verified)
	if err != nil {
		log.Printf("Error during query execution or scanning: %v", err)
		return false, err
	}
	return verified, nil
}

func IsTokenVerificationEmailExists(token string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email_verification_token = $1
		)`
	err := database.DB.QueryRow(
		context.Background(),
		query,
		token,
	).Scan(&exists)
	if err != nil {
		log.Printf("Error during query execution or scanning: %v", err)
		return false, err
	}
	return exists, nil
}

func UpdateEmailVerificationToken(email, token string) error {
	query := `
		UPDATE users
		SET email_verification_token = $1
		WHERE email = $2`
	_, err := database.DB.Exec(
		context.Background(),
		query,
		token,
		email,
	)
	if err != nil {
		log.Printf("Error during query execution or scanning: %v", err)
		return err
	}
	return nil

}

// VerifyEmail updates the email_verified field in the users table and return error if verification token is invalid
func VerifyEmail(token string) error {
	query := `
		UPDATE users
		SET email_verified = TRUE
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
