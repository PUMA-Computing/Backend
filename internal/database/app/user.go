package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
)

func GetUserByUsernameOrEmail(username string) (*models.User, error) {
	var user models.User
	var userID string
	var middleName sql.NullString
	err := database.DB.QueryRow(context.Background(), "SELECT * FROM users WHERE username = $1 OR email = $1", username).Scan(
		&userID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email,
		&user.StudentID, &user.Major, &user.ProfilePicture, &user.DateOfBirth, &user.RoleID, &user.CreatedAt,
		&user.UpdatedAt, &user.Year, &user.InstitutionName, &user.Gender)
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
	user.MiddleName = nil
	if middleName.Valid {
		user.MiddleName = &middleName.String
	}
	return &user, nil
}

func IsUsernameExists(username string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"
	var exists bool
	err := database.DB.QueryRow(context.Background(), query, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func IsEmailExists(email string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	var exists bool
	err := database.DB.QueryRow(context.Background(), query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	query := "SELECT * FROM users WHERE username = $1"
	var user models.User
	var userID string
	err := database.DB.QueryRow(context.Background(), query, username).Scan(
		&userID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email,
		&user.StudentID, &user.Major, &user.ProfilePicture, &user.DateOfBirth, &user.RoleID, &user.CreatedAt,
		&user.UpdatedAt, &user.Year, &user.InstitutionName, &user.Gender)
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
	err := database.DB.QueryRow(context.Background(), query, email).Scan(
		&userID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email,
		&user.StudentID, &user.Major, &user.ProfilePicture, &user.DateOfBirth, &user.RoleID, &user.CreatedAt,
		&user.UpdatedAt, &user.Year, &user.InstitutionName, &user.Gender, &user.TwoFAEnabled, &user.TwoFAImage, &user.TwoFASecret)
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

	err := database.DB.QueryRow(context.Background(), `
		SELECT id, username, password, first_name, middle_name, last_name, email, student_id, major, profile_picture, date_of_birth, role_id, created_at, updated_at, year, institution_name, gender, twofa_enabled, twofa_image, twofa_secret
		FROM users WHERE id = $1`, userID).Scan(
		&user.ID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email,
		&user.StudentID, &user.Major, &user.ProfilePicture, &user.DateOfBirth, &user.RoleID, &user.CreatedAt,
		&user.UpdatedAt, &user.Year, &user.InstitutionName, &user.Gender, &user.TwoFAEnabled, &user.TwoFAImage, &user.TwoFASecret)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func GetUserByStudentID(studentID string) (*models.User, error) {
	var user models.User
	var userID string
	err := database.DB.QueryRow(context.Background(), "SELECT * FROM users WHERE student_id = $1", studentID).Scan(
		&userID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email,
		&user.StudentID, &user.Major, &user.ProfilePicture, &user.DateOfBirth, &user.RoleID, &user.CreatedAt,
		&user.UpdatedAt, &user.Year, &user.InstitutionName, &user.Gender)
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

func GetRoleIDByUserID(userID uuid.UUID) (int, error) {
	var roleID int
	err := database.DB.QueryRow(context.Background(), "SELECT role_id FROM users WHERE id = $1", userID).Scan(&roleID)
	if err != nil {
		return 0, err
	}
	return roleID, nil
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
	query := "UPDATE users SET "
	args := []interface{}{}
	argID := 1

	if updatedUser.Username != "" {
		query += "username = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.Username)
		argID++
	}
	if updatedUser.Password != "" {
		query += "password = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.Password)
		argID++
	}
	if updatedUser.FirstName != "" {
		query += "first_name = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.FirstName)
		argID++
	}
	if updatedUser.MiddleName != nil && *updatedUser.MiddleName != "" {
		query += "middle_name = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.MiddleName)
		argID++
	}
	if updatedUser.LastName != "" {
		query += "last_name = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.LastName)
		argID++
	}
	if updatedUser.Email != "" {
		query += "email = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.Email)
		argID++
	}
	if updatedUser.StudentID != "" {
		query += "student_id = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.StudentID)
		argID++
	}
	if updatedUser.Major != "" {
		query += "major = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.Major)
		argID++
	}
	if updatedUser.Year != "" {
		query += "year = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.Year)
		argID++
	}
	if updatedUser.DateOfBirth != nil {
		query += "date_of_birth = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.DateOfBirth)
		argID++
	}
	if updatedUser.RoleID != 0 {
		query += "role_id = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.RoleID)
		argID++
	}
	if !updatedUser.UpdatedAt.IsZero() {
		query += "updated_at = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.UpdatedAt)
		argID++
	}
	if updatedUser.InstitutionName != nil && *updatedUser.InstitutionName != "" {
		query += "institution_name = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.InstitutionName)
		argID++
	}
	if updatedUser.Gender != "" {
		query += "gender = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.Gender)
		argID++
	}

	if updatedUser.TwoFASecret != nil && *updatedUser.TwoFASecret != "" {
		query += "twofa_secret = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.TwoFASecret)
		argID++
	}

	if updatedUser.TwoFAImage != nil && *updatedUser.TwoFAImage != "" {
		query += "twofa_image = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.TwoFAImage)
		argID++
	}

	if updatedUser.TwoFAEnabled != false {
		query += "twofa_enabled = $" + strconv.Itoa(argID) + ", "
		args = append(args, updatedUser.TwoFAEnabled)
		argID++

	}

	// Remove the last comma and space
	query = strings.TrimSuffix(query, ", ")

	// Add the WHERE clause
	query += " WHERE id = $" + strconv.Itoa(argID)
	args = append(args, UserID)

	_, err := database.DB.Exec(context.Background(), query, args...)
	return err
}

func ToggleTwoFA(userID uuid.UUID, enable bool) error {
	var err error
	if enable {
		_, err = database.DB.Exec(context.Background(),
			"UPDATE users SET twofa_enabled = $1 WHERE id = $2",
			enable, userID)
	} else {
		_, err = database.DB.Exec(context.Background(),
			"UPDATE users SET twofa_enabled = $1, twofa_image = NULL, twofa_secret = NULL WHERE id = $2",
			enable, userID)
	}
	return err
}

func AdminUpdateRoleAndStudentIDVerified(userID uuid.UUID, roleID int, studentIDVerified bool) error {
	_, err := database.DB.Exec(context.Background(), "UPDATE users SET role_id = $1, student_id_verified = $2 WHERE id = $3", roleID, studentIDVerified, userID)
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
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User

		log.Println("before scan")
		err := rows.Scan(
			&user.ID, &user.Username, &user.Password, &user.FirstName, &user.MiddleName, &user.LastName,
			&user.Email, &user.StudentID, &user.Major, &user.ProfilePicture, &user.DateOfBirth,
			&user.RoleID, &user.CreatedAt, &user.UpdatedAt, &user.Year, &user.EmailVerified,
			&user.EmailVerificationToken, &user.PasswordResetToken, &user.PasswordResetExpires,
			&user.StudentIDVerified, &user.StudentIDVerification, &user.InstitutionName, &user.Gender,
			&user.TwoFAEnabled, &user.TwoFAImage, &user.TwoFASecret,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		log.Println("after scan")

		users = append(users, user)

		log.Println("after append")
	}

	if err = rows.Err(); err != nil {
		log.Println("Error after row iteration:", err)
		return nil, err
	}

	log.Println("after loop")

	return users, nil
}

func UploadProfilePicture(userID uuid.UUID, profilePicture string) error {
	_, err := database.DB.Exec(context.Background(), "UPDATE users SET profile_picture = $1 WHERE id = $2", profilePicture, userID)
	return err
}
