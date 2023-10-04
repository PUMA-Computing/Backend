package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
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
		existingUserByUsername, err := app.GetUserByUsername(user.Username)
		if err != nil {
			return err
		}

		existingUserByEmail, err := app.GetUserByEmail(user.Email)
		if err != nil {
			return err
		}

		if existingUserByUsername != nil {
			return errors.New("username already exists")
		}

		if existingUserByEmail != nil {
			return errors.New("email already exists")
		}
	}

	user.ID = uuid.New()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.RoleID = 2
	user.CreatedAt = time.Time{}
	user.UpdatedAt = time.Time{}

	err = app.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) Login(username, password string) (*models.User, error) {
	var u models.User
	user, err := app.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (us *UserService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return app.GetUserByID(userID)
}

func (us *UserService) EditUser(userID uuid.UUID, updatedUser *models.User) error {
	existingUser, err := app.GetUserByUsernameOrEmail(updatedUser.Username, updatedUser.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("username or email already exists")
	}

	return app.UpdateUser(userID, updatedUser)
}

func (us *UserService) DeleteUser(userID uuid.UUID) error {
	return app.DeleteUser(userID)
}

func (us *UserService) ListUsers() ([]*models.User, error) {
	return app.ListUsers()
}
