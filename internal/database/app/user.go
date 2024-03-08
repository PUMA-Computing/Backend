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

func GetUserByUsernameOrEmail(username, email string) (*models.User, error) {
	query := "SELECT * FROM users WHERE username = $1 OR email = $2"
	var user models.User
	err := database.DB.QueryRow(context.Background(), query, username, email).Scan(&user.ID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.StudentID, &user.Major, &user.Year, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found, return nil
		}
		return nil, err
	}
	return &user, nil
}

func IsUsernameOrEmailExists(username string) (bool, error) {
	log.Printf("username: %v", username)
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $1)"
	log.Printf("query: %v", query)
	var exists bool
	err := database.DB.QueryRow(context.Background(), query, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	log.Printf("exists: %v", exists)

	return exists, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	query := "SELECT * FROM users WHERE username = $1"
	var user models.User
	var userID string
	err := database.DB.QueryRow(context.Background(), query, username).Scan(&userID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.StudentID, &user.Major, &user.Year, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found, return nil
		}
		return nil, err // Return the actual error here
	}
	user.ID, err = uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	var user models.User
	var userID string
	err := database.DB.QueryRow(context.Background(), query, email).Scan(&userID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.StudentID, &user.Major, &user.Year, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found, return nil
		}
		return nil, err
	}
	user.ID, err = uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	var middleName sql.NullString
	var profilePicture sql.NullString
	var dateOfBirth sql.NullTime

	err := database.DB.QueryRow(context.Background(), "SELECT id, username, first_name, middle_name, last_name, email, student_id, major, year, profile_picture, date_of_birth, role_id, created_at, updated_at FROM users WHERE id = $1", userID).Scan(
		&user.ID, &user.Username, &user.FirstName, &middleName, &user.LastName, &user.Email,
		&user.StudentID, &user.Major, &user.Year, &profilePicture, &dateOfBirth, &user.RoleID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	user.MiddleName = ""
	if middleName.Valid {
		user.MiddleName = middleName.String
	}

	user.ProfilePicture = nil
	if profilePicture.Valid {
		user.ProfilePicture = &profilePicture.String
	}

	user.DateOfBirth = nil
	if dateOfBirth.Valid {
		user.DateOfBirth = &dateOfBirth.Time
	}

	return &user, nil
}

func GetUserByStudentID(studentID string) (*models.User, error) {
	var user models.User
	var userID string
	err := database.DB.QueryRow(context.Background(), "SELECT * FROM users WHERE student_id = $1", studentID).Scan(&userID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.StudentID, &user.Major, &user.Year, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	user.ID, err = uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func CheckStudentIDExists(studentID string) (bool, error) {
	var exists bool
	err := database.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE student_id = $1)", studentID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil

}

func UpdateUser(UserID uuid.UUID, updatedUser *models.User) error {
	log.Printf("updatedUser: %v", updatedUser)
	_, err := database.DB.Exec(context.Background(), "UPDATE users SET username = $1, password = $2, first_name = $3, middle_name = $4, last_name = $5, email = $6, student_id = $7, major = $8, year = $9, role_id = $10, WHERE id = $11",
		updatedUser.Username, updatedUser.Password, updatedUser.FirstName, updatedUser.MiddleName, updatedUser.LastName, updatedUser.Email, updatedUser.StudentID, updatedUser.Major, &updatedUser.Year, updatedUser.RoleID, UserID)
	log.Printf("err: %v", err)
	return err
}

func DeleteUser(userID uuid.UUID) error {
	_, err := database.DB.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
	return err
}

func ListUsers() ([]models.User, error) {
	var users []models.User

	log.Println("before query")

	rows, err := database.DB.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		// log the error
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		var middleName sql.NullString
		var profilePicture sql.NullString
		var dateOfBirth sql.NullTime

		log.Println("before scan")
		err := rows.Scan(
			&user.ID, &user.Username, &user.Password, &user.FirstName, &middleName, &user.LastName,
			&user.Email, &user.StudentID, &user.Major, &profilePicture, &dateOfBirth,
			&user.RoleID, &user.CreatedAt, &user.UpdatedAt, &user.Year,
		)
		if err != nil {
			return nil, err
		}

		log.Println("after scan")

		if middleName.Valid {
			user.MiddleName = middleName.String
		} else {
			user.MiddleName = ""
		}

		if profilePicture.Valid {
			user.ProfilePicture = &profilePicture.String
		} else {
			user.ProfilePicture = nil
		}

		if dateOfBirth.Valid {
			user.DateOfBirth = &dateOfBirth.Time
		} else {
			user.DateOfBirth = nil
		}

		users = append(users, user)
	}
	return users, nil
}
