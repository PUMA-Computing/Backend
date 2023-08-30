package service

import (
	"Backend/internal/app/domain"
	"Backend/internal/app/repository"
	"Backend/internal/utils"
	"strconv"
)

type AuthResponse struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}

type UserServices interface {
	RegisterUser(user *domain.User) error
	AuthenticateUser(email, password string) (*AuthResponse, error)
}

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) RegisterUser(user *domain.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.userRepository.RegisterUser(user)
}

func (s *UserService) AuthenticateUser(email, password string) (*AuthResponse, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := utils.ComparePassword(user.Password, password); err != nil {
		return nil, err
	}

	token, err := utils.GenerateJWTToken(strconv.FormatInt(user.ID, 10), domain.Role(user.Role))
	if err != nil {
		return nil, err
	}

	response := &AuthResponse{
		User:  user,
		Token: token,
	}

	return response, nil
}
