package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return app.GetUserByID(userID)
}

func (us *UserService) EditUser(userID uuid.UUID, updatedUser *models.User) error {
	existingUser, err := app.GetUserByID(userID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	if updatedUser.Username != "" {
		existingUser.Username = updatedUser.Username
	}

	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		existingUser.Password = string(hashedPassword)
	}

	if updatedUser.FirstName != "" {
		existingUser.FirstName = updatedUser.FirstName
	}

	if updatedUser.MiddleName != "" {
		existingUser.MiddleName = updatedUser.MiddleName
	}

	if updatedUser.LastName != "" {
		existingUser.LastName = updatedUser.LastName
	}

	if updatedUser.Email != "" {
		existingUser.Email = updatedUser.Email
	}

	if updatedUser.StudentID != "" {
		existingUser.StudentID = updatedUser.StudentID
	}

	if updatedUser.Major != "" {
		existingUser.Major = updatedUser.Major
	}

	return app.UpdateUser(userID, existingUser)
}

func (us *UserService) DeleteUser(userID uuid.UUID) error {
	return app.DeleteUser(userID)
}

func (us *UserService) ListUsers() ([]*models.User, error) {
	return app.ListUsers()
}
