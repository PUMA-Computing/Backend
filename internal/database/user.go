package database

import (
	"Backend/internal/models"
	"context"
	"github.com/google/uuid"
)

func CreateUser(user *models.User) error {
	_, err := DB.Exec(context.Background(), "INSERT INTO users (username, password, first_name, middle_name, last_name, email, student_id, major, role_id) VALUES ($1, $2, $3, $4, $5, $6, $7, &8, $9)",
		user.Username, user.Password, user.FirstName, user.MiddleName, user.LastName, user.Email, user.StudentID, user.Major, user.RoleID)
	return err
}

func GetUserByUsernameOrEmail(username, email string) (*models.User, error) {
	var user models.User
	err := DB.QueryRow(context.Background(), "SELECT * FROM users WHERE username = $1 OR email = $2", username, email).Scan(&user.ID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.StudentID, &user.Major, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := DB.QueryRow(context.Background(), "SELECT * FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.StudentID, &user.Major, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := DB.QueryRow(context.Background(), "SELECT * FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.StudentID, &user.Major, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(UserID uuid.UUID, updatedUser *models.User) error {
	_, err := DB.Exec(context.Background(), "UPDATE users SET username = $1, password = $2, first_name = $3, middle_name = $4, last_name = $5, email = $6, student_id = &7, major = &8, role_id = $9 WHERE id = $10",
		updatedUser.Username, updatedUser.Password, updatedUser.FirstName, updatedUser.MiddleName, updatedUser.LastName, updatedUser.Email, updatedUser.StudentID, updatedUser.Major, updatedUser.RoleID, UserID)
	return err
}

func DeleteUser(userID uuid.UUID) error {
	_, err := DB.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
	return err
}

func ListUsers() ([]*models.User, error) {
	var users []*models.User
	rows, err := DB.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.StudentID, &user.Major, &user.RoleID, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
