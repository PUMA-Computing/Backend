package userService

import (
	"Backend/internal/app/domain/user"
	user2 "Backend/internal/app/interfaces/repository/userRepository"
	token2 "Backend/internal/utils/token"
	"Backend/pkg/bcrypt"
	"github.com/google/uuid"
)

type AuthResponse struct {
	User        *user.User `json:"userService"`
	AccessToken string     `json:"access_token"`
}

type UserServices interface {
	RegisterUser(user *user.User) error
	AuthenticateUser(email, password string) (*AuthResponse, error)
	GetAllUsers() ([]*user.User, error)
	GetUserByID(id uuid.UUID) (*user.User, error)
	GetUserRoleByID(id uuid.UUID) (int, error)
	GetUserByEmail(email string) (*user.User, error)
	GetUserRoleByEmail(email string) (int, error)
	UpdateUser(user *user.User) error
	DeleteUser(id uuid.UUID) error
}

type UserService struct {
	userRepository user2.UserRepository
}

func NewUserService(userRepository user2.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (u *UserService) RegisterUser(user *user.User) error {
	hashedPassword, err := bcrypt.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	user.ID = uuid.New()
	return u.userRepository.RegisterUser(user)
}

func (u *UserService) AuthenticateUser(email, password string) (*AuthResponse, error) {
	user, err := u.userRepository.AuthenticateUser(email, password)
	if err != nil {
		return nil, err
	}

	token, err := token2.GenerateJWTToken(user.ID, user.RoleID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{User: user, AccessToken: token}, nil
}

func (u *UserService) GetAllUsers() ([]*user.User, error) {
	return u.userRepository.GetAllUsers()
}

func (u *UserService) GetUserByEmail(email string) (*user.User, error) {
	return u.userRepository.GetUserByEmail(email)
}

func (u *UserService) GetUserRoleByEmail(email string) (int, error) {
	return u.userRepository.GetUserRoleByEmail(email)
}

func (u *UserService) GetUserByID(id uuid.UUID) (*user.User, error) {
	return u.userRepository.GetUserByID(id)
}

func (u *UserService) GetUserRoleByID(id uuid.UUID) (int, error) {
	return u.userRepository.GetUserRoleByID(id)
}

func (u *UserService) UpdateUser(user *user.User) error {
	return u.userRepository.UpdateUser(user)
}

func (u *UserService) DeleteUser(id uuid.UUID) error {
	return u.userRepository.DeleteUser(id)
}
