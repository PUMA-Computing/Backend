package userService

import (
	"Backend/internal/app/domain/roles"
	"Backend/internal/app/domain/user"
	user2 "Backend/internal/app/interfaces/repository/userRepository"
	token2 "Backend/internal/utils/token"
	"Backend/pkg/bcrypt"
	"github.com/google/uuid"
	"time"
)

type AuthResponse struct {
	User        *user.User `json:"userService"`
	AccessToken string     `json:"access_token"`
}

type UserServices interface {
	RegisterUser(user *user.User) error
	AuthenticateUser(email, password string) (*AuthResponse, error)
	Logout(id uuid.UUID) error
	GetAllUsers() ([]*user.User, error)
	GetUserByID(id uuid.UUID) (*user.User, error)
	GetUserRoleByID(id uuid.UUID) (int, error)
	GetUserByEmail(email string) (*user.User, error)
	GetUserRoleByEmail(email string) (int, error)
	GetUserRoleAndPermissions(id uuid.UUID) (*roles.Role, []*roles.Permission, []*roles.RolePermissions, error)
	UpdateUser(user *user.User) error
	DeleteUser(id uuid.UUID) error
	CreateRole(role *roles.Role) error
	UpdateRole(role *roles.Role) error
	DeleteRole(id int) error
	AssignUserRole(id uuid.UUID, roleID int) error
	AssignRolePermission(roleID int, permissionID int) error
	DeleteRolePermission(roleID int, permissionID int) error
	HasPermission(roleID int, permissionID int) (bool, error)
}

type UserService struct {
	userRepository user2.UserRepository
}

func NewUserService(userRepository user2.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

/**
 * Authentication Management
 */

func (u *UserService) RegisterUser(user *user.User) error {
	hashedPassword, err := bcrypt.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	user.ID = uuid.New()
	user.RoleID = roles.RoleComputizen
	user.CreatedAt = time.Now()
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

func (u *UserService) Logout(id uuid.UUID) error {
	return u.userRepository.Logout(id)
}

/**
 * End Authentication Management
 */

/**
 * Manage Profile Management
 */

func (u *UserService) UpdateUser(user *user.User) error {
	return u.userRepository.UpdateUser(user)
}

func (u *UserService) DeleteUser(id uuid.UUID) error {
	return u.userRepository.DeleteUser(id)
}

/**
 * End Manage Profile Management
 */

/**
 * Get User Management
 */

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

func (u *UserService) GetUserRoleAndPermissions(id uuid.UUID) (*roles.Role, []*roles.Permission, []*roles.RolePermissions, error) {
	return u.userRepository.GetUserRoleAndPermissions(id)
}

/**
 * End Get User Management
 */

/**
 * Role and Permission Management
 */

func (u *UserService) CreateRole(role *roles.Role) error {
	return u.userRepository.CreateRole(role)
}

func (u *UserService) UpdateRole(role *roles.Role) error {
	return u.userRepository.UpdateRole(role)
}

func (u *UserService) DeleteRole(id int) error {
	return u.userRepository.DeleteRole(id)
}

func (u *UserService) AssignUserRole(id uuid.UUID, roleID int) error {
	return u.userRepository.AssignUserRole(id, roleID)
}

func (u *UserService) AssignRolePermission(roleID int, permissionID int) error {
	return u.userRepository.AssignRolePermission(roleID, permissionID)
}

func (u *UserService) DeleteRolePermission(roleID int, permissionID int) error {
	return u.userRepository.DeleteRolePermission(roleID, permissionID)
}

func (u *UserService) HasPermission(roleID int, permissionID int) (bool, error) {
	return u.userRepository.HasPermission(roleID, permissionID)
}

/**
 * End Role and Permission Management
 */
