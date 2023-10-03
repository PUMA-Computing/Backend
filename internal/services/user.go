package services

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) RegisterUser(user *models.User) error {
	existingUser, err := database.GetUserByUsernameOrEmail(user.Username, user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("username or email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	err = database.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) Login(username, password string) (*models.User, error) {
	var u models.User
	user, err := database.GetUserByUsername(username)
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
	return database.GetUserByID(userID)
}

func (us *UserService) EditUser(userID uuid.UUID, updatedUser *models.User) error {
	existingUser, err := database.GetUserByUsernameOrEmail(updatedUser.Username, updatedUser.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("username or email already exists")
	}

	return database.UpdateUser(userID, updatedUser)
}

func (us *UserService) DeleteUser(userID uuid.UUID) error {
	return database.DeleteUser(userID)
}

func (us *UserService) ListUsers() ([]*models.User, error) {
	return database.ListUsers()
}
