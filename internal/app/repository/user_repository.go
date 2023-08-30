package repository

import "Backend/internal/app/domain"

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
}
