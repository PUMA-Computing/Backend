package userService

import (
	"Backend/internal/app/domain/user"
	user2 "Backend/internal/app/interfaces/repository/userRepository"
	"Backend/internal/utils/token"
	"Backend/pkg/bcrypt"
	"github.com/google/uuid"
)

type AuthResponse struct {
	User  *user.User `json:"userService"`
	Token string     `json:"token"`
}

type UserServices interface {
	RegisterUser(user *user.User) error
	AuthenticateUser(email, password string) (*AuthResponse, error)
	GetUserByID(id uuid.UUID) (*user.User, error)
	UpdateUser(user *user.User) error
	DeleteUser(id uuid.UUID) error
}

type UserService struct {
	userRepository user2.UserRepository
}

func NewUserService(userRepository user2.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) RegisterUser(user *user.User) error {
	hashedPassword, err := bcrypt.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	user.ID = uuid.New()
	return s.userRepository.RegisterUser(user)
}

func (s *UserService) AuthenticateUser(email, password string) (*AuthResponse, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.ComparePassword(user.Password, password); err != nil {
		return nil, err
	}

	token, err := token.GenerateJWTToken(user.ID, user.RoleID)
	if err != nil {
		return nil, err
	}

	response := &AuthResponse{
		User:  user,
		Token: token,
	}

	return response, nil
}

func (s *UserService) GetUserByEmail(email string) (*user.User, error) {
	return s.userRepository.GetUserByEmail(email)
}

func (s *UserService) GetUserByID(id uuid.UUID) (*user.User, error) {
	return s.userRepository.GetUserByID(id)
}

func (s *UserService) UpdateUser(user *user.User) error {
	return s.userRepository.UpdateUser(user)
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.userRepository.DeleteUser(id)
}
