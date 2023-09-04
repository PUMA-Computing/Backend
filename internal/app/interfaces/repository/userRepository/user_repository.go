package userRepository

import (
	"Backend/internal/app/domain/user"
	"Backend/pkg/bcrypt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(user *user.User) error
	AuthenticateUser(email, password string) (*user.User, error)
	GetAllUsers() ([]*user.User, error)
	GetUserByEmail(email string) (*user.User, error)
	GetUserByID(id uuid.UUID) (*user.User, error)
	UpdateUser(user *user.User) error
	DeleteUser(id uuid.UUID) error
	GetUserRoleByID(id uuid.UUID) (int, error)
	GetUserRoleByEmail(email string) (int, error)
}

type PostgresUserRepository struct {
	session *gorm.DB
}

func NewPostgresUserRepository(session *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{session: session}
}

func (p *PostgresUserRepository) RegisterUser(user *user.User) error {
	return p.session.Create(user).Error
}

func (p *PostgresUserRepository) AuthenticateUser(email, password string) (*user.User, error) {
	var user user.User
	if err := p.session.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.ComparePassword(user.Password, password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *PostgresUserRepository) GetAllUsers() ([]*user.User, error) {
	var users []*user.User
	if err := p.session.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (p *PostgresUserRepository) GetUserByID(id uuid.UUID) (*user.User, error) {
	var user user.User
	if err := p.session.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *PostgresUserRepository) GetUserRoleByID(id uuid.UUID) (int, error) {
	var user user.User
	if err := p.session.Where("id = ?", id).First(&user).Error; err != nil {
		return 0, err
	}

	return user.RoleID, nil
}
func (p *PostgresUserRepository) GetUserByEmail(email string) (*user.User, error) {
	var user user.User
	if err := p.session.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *PostgresUserRepository) GetUserRoleByEmail(email string) (int, error) {
	var user user.User
	if err := p.session.Where("email = ?", email).First(&user).Error; err != nil {
		return 0, err
	}

	return user.RoleID, nil
}

func (p *PostgresUserRepository) UpdateUser(user *user.User) error {
	return p.session.Save(user).Error
}

func (p *PostgresUserRepository) DeleteUser(id uuid.UUID) error {
	return p.session.Delete(&user.User{}, id).Error
}
