package service

import (
	"Backend/internal/app/domain"
	"Backend/internal/app/repository"
	"Backend/internal/utils"
)

type UserServices interface {
	CreateUser(user *domain.User) error
	AuthenticateUser(email, password string) (*domain.User, error)
}

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) CreateUser(user *domain.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.userRepository.CreateUser(user)
}

func (s *UserService) AuthenticateUser(email, password string) (*domain.User, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := utils.ComparePassword(user.Password, password); err != nil {
		return nil, err
	}

	return user, nil
}
