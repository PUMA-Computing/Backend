package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
	"Backend/pkg/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (as *AuthService) RegisterUser(user *models.User) error {
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

func (as *AuthService) LoginUser(username, email, password string) (*models.User, error) {
	user, err := app.AuthenticateUser(username, email, password)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &utils.UnauthorizedError{Message: "invalid credentials"}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Error comparing hashed password: %v", err)
		return nil, &utils.UnauthorizedError{Message: "invalid credentials"}
	}

	return user, nil
}
