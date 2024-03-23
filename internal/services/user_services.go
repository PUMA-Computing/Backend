package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
	"errors"
	"github.com/google/uuid"
	"log"
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
		existingUser.Password = updatedUser.Password
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

func (us *UserService) GetUserByUsername(username string) (*models.User, error) {
	return app.GetUserByUsername(username)
}

func (us *UserService) GetUserByEmail(email string) (*models.User, error) {
	return app.GetUserByEmail(email)
}

func (us *UserService) CheckStudentIDExists(studentID string) (bool, error) {
	return app.CheckStudentIDExists(studentID)
}

func (us *UserService) ListUsers() ([]models.User, error) {
	log.Println("service list users")
	return app.ListUsers()
}
