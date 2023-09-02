package userRepository

import (
	"Backend/internal/app/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(user *user.User) error
	GetUserByEmail(email string) (*user.User, error)
	GetUserByID(id uuid.UUID) (*user.User, error)
	UpdateUser(user *user.User) error
	DeleteUser(id uuid.UUID) error
}

type PostgresUserRepository struct {
	session *gorm.DB
}

func NewPostgresUserRepository(session *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{session: session}
}

func (r *PostgresUserRepository) RegisterUser(user *user.User) error {
	return r.session.Create(user).Error
}

func (r *PostgresUserRepository) GetUserByEmail(email string) (*user.User, error) {
	var user user.User
	if err := r.session.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetUserByID(id uuid.UUID) (*user.User, error) {
	var user user.User
	if err := r.session.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) UpdateUser(user *user.User) error {
	return r.session.Save(user).Error
}

func (r *PostgresUserRepository) DeleteUser(id uuid.UUID) error {
	return r.session.Delete(&user.User{}, id).Error
}
