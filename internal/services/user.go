package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
	"Backend/pkg/utils"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) RegisterUser(user *models.User) error {
	hasRows, err := app.TableHasRows("users")
	if err != nil {
		return err
	}

	if hasRows {
		existingUserByUsername, err := app.IsUsernameExists(user.Username)
		if err != nil {
			return err
		}

		existingUserByEmail, err := app.IsEmailExists(user.Email)
		if err != nil {
			return err
		}

		if existingUserByUsername {
			return &utils.ConflictError{Message: "username already exists"}
		}

		if existingUserByEmail {
			return &utils.ConflictError{Message: "email already exists"}
		}
	}

	user.ID = uuid.New()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	err = app.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) Login(username, password string) (*models.User, error) {
	log.Println("Before calling GetUserByUsername")
	user, err := app.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	log.Println("After calling GetUserByUsername")

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	log.Println("After calling CompareHashAndPassword")

	return user, nil
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
