package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
	"Backend/pkg/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
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
		var checkUsernameOrEmail string
		if user.Username != "" {
			checkUsernameOrEmail = user.Username
		} else {
			checkUsernameOrEmail = user.Email
		}
		existingUsernameOrEmail, err := app.IsUsernameOrEmailExists(checkUsernameOrEmail)
		if err != nil {
			return err
		}

		if existingUsernameOrEmail {
			return &utils.ConflictError{Message: "Username or email already exists"}
		}

	}

	user.ID = uuid.New()
	user.RoleID = 2
	user.CreatedAt = time.Time{}
	user.UpdatedAt = time.Time{}

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

func (as *AuthService) LoginUser(username string, password string) (*models.User, error) {
	user, err := app.AuthenticateUser(username)
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

func (as *AuthService) IsUsernameOrEmailExists(username string, email string) (bool, error) {
	if username != "" {
		return app.IsUsernameOrEmailExists(username)
	} else {
		return app.IsUsernameOrEmailExists(email)
	}
}

func (as *AuthService) CheckStudentIDExists(studentID string) (bool, error) {
	return app.CheckStudentIDExists(studentID)
}
