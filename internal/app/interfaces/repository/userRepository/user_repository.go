package userRepository

import (
	"Backend/internal/app/domain/roles"
	"Backend/internal/app/domain/user"
	"Backend/internal/app/interfaces/repository/postgresRepository"
	"Backend/pkg/bcrypt"
	"errors"
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
	GetUserRoleAndPermissions(id uuid.UUID) (*roles.Role, []*roles.Permission, []*roles.RolePermissions, error)
	Logout(id uuid.UUID) error
	CreateRole(role *roles.Role) error
	UpdateRole(role *roles.Role) error
	DeleteRole(id int) error
	AssignUserRole(userID uuid.UUID, roleID int) error
	AssignRolePermission(roleID int, permissionID int) error
	DeleteRolePermission(roleID int, permissionID int) error
	HasPermission(roleID int, permissionID int) (bool, error)
	RoleExists(roleID int) (bool, error)
}

type PostgresUserRepository struct {
	session *gorm.DB
}

func NewPostgresUserRepository(session *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{session: session}
}

/**
 * Authentication Management
 */

func (p *PostgresUserRepository) RegisterUser(user *user.User) error {
	existingUser := &user
	if err := p.session.Where("email = ?", user.Email).Or("nim = ?", user.NIM).First(&existingUser).Error; err == nil {
		return errors.New("user with the same email or nim already exists")
	}

	return p.session.Create(user).Error
}

func (p *PostgresUserRepository) AuthenticateUser(email, password string) (*user.User, error) {
	var u user.User
	if err := p.session.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.ComparePassword(u.Password, password); err != nil {
		return nil, err
	}

	return &u, nil
}

func (p *PostgresUserRepository) Logout(id uuid.UUID) error {
	if err := p.session.Delete(&postgresRepository.SessionData{}, id).Error; err != nil {
		return err
	}

	return nil
}

/**
 * End Authentication Management
 */

/**
 * Manage Profile Management
 */

func (p *PostgresUserRepository) UpdateUser(user *user.User) error {
	return p.session.Save(user).Error
}

func (p *PostgresUserRepository) DeleteUser(id uuid.UUID) error {
	return p.session.Delete(&user.User{}, id).Error
}

/**
 * End Manage Profile Management
 */

/**
 * Get User Management
 */

func (p *PostgresUserRepository) GetAllUsers() ([]*user.User, error) {
	var users []*user.User
	if err := p.session.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (p *PostgresUserRepository) GetUserByID(id uuid.UUID) (*user.User, error) {
	var u user.User
	if err := p.session.Where("id = ?", id).First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (p *PostgresUserRepository) GetUserRoleByID(id uuid.UUID) (int, error) {
	var u user.User
	if err := p.session.Where("id = ?", id).First(&u).Error; err != nil {
		return 0, err
	}

	return u.RoleID, nil
}

func (p *PostgresUserRepository) GetUserByEmail(email string) (*user.User, error) {
	var u user.User
	if err := p.session.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (p *PostgresUserRepository) GetUserByNIM(nim string) (*user.User, error) {
	var u user.User
	if err := p.session.Where("nim = ?", nim).First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (p *PostgresUserRepository) GetUserRoleByEmail(email string) (int, error) {
	var u user.User
	if err := p.session.Where("email = ?", email).First(&u).Error; err != nil {
		return 0, err
	}

	return u.RoleID, nil
}

func (p *PostgresUserRepository) GetUserRoleAndPermissions(id uuid.UUID) (*roles.Role, []*roles.Permission, []*roles.RolePermissions, error) {
	var userRole roles.Role
	var userPermissions []*roles.Permission
	var rolePermissions []*roles.RolePermissions

	if err := p.session.Where("id = ?", id).First(&userRole).Error; err != nil {
		return nil, nil, nil, err
	}

	if err := p.session.Model(&userRole).Association("Permissions").Find(&userPermissions); err != nil {
		return nil, nil, nil, err
	}

	if err := p.session.Where("role_id = ?", userRole.ID).Find(&rolePermissions).Error; err != nil {
		return nil, nil, nil, err // Return the error here
	}

	return &userRole, userPermissions, rolePermissions, nil
}

/**
 * End Get User Management
 */

/**
 * Role and Permission Management
 */

func (p *PostgresUserRepository) ManageRoleAndPermissions() ([]*roles.Role, []*roles.Permission, []*roles.RolePermissions, error) {
	var role []*roles.Role
	var permissions []*roles.Permission
	var rolePermissions []*roles.RolePermissions

	if err := p.session.Find(&role).Error; err != nil {
		return nil, nil, nil, err
	}

	if err := p.session.Find(&permissions).Error; err != nil {
		return nil, nil, nil, err
	}

	if err := p.session.Find(&rolePermissions).Error; err != nil {
		return nil, nil, nil, err
	}

	return role, permissions, rolePermissions, nil

}

func (p *PostgresUserRepository) CreateRole(role *roles.Role) error {
	return p.session.Create(role).Error
}

func (p *PostgresUserRepository) UpdateRole(role *roles.Role) error {
	return p.session.Save(role).Error
}

func (p *PostgresUserRepository) DeleteRole(id int) error {
	return p.session.Delete(&roles.Role{}, id).Error
}

func (p *PostgresUserRepository) AssignUserRole(userID uuid.UUID, roleID int) error {
	if err := p.session.Model(&user.User{}).Where("id = ?", userID).Update("role_id", roleID).Error; err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserRepository) AssignRolePermission(roleID int, permissionID int) error {
	var roleExists, permissionExists bool
	if err := p.session.Where("id = ?", roleID).First(&roles.Role{}).Error; err != nil {
		roleExists = true
	}
	if err := p.session.Where("id = ?", permissionID).First(&roles.Permission{}).Error; err != nil {
		permissionExists = true
	}

	if !roleExists || !permissionExists {
		return errors.New("role or permission doesn't exist")
	}

	rolePermissionMapping := roles.RolePermissions{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	if err := p.session.Create(rolePermissionMapping).Error; err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserRepository) DeleteRolePermission(roleID int, permissionID int) error {
	var roleExists, permissionExists bool
	if err := p.session.Where("id = ?", roleID).First(&roles.Role{}).Error; err != nil {
		roleExists = true
	}
	if err := p.session.Where("id = ?", permissionID).First(&roles.Permission{}).Error; err != nil {
		permissionExists = true
	}

	if !roleExists || !permissionExists {
		return errors.New("role or permission doesn't exist")
	}

	if err := p.session.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&roles.RolePermissions{}).Error; err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserRepository) RoleExists(roleID int) (bool, error) {
	var count int64
	if err := p.session.Model(&roles.Role{}).Where("id = ?", roleID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (p *PostgresUserRepository) permissionExists(permissionID int) (bool, error) {
	var count int64
	if err := p.session.Model(&roles.Permission{}).Where("id = ?", permissionID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (p *PostgresUserRepository) HasPermission(roleID int, permissionID int) (bool, error) {
	roleExists, err := p.RoleExists(roleID)
	if err != nil {
		return false, err
	}

	permissionExists, err := p.permissionExists(permissionID)
	if err != nil {
		return false, err
	}

	if !roleExists || !permissionExists {
		return false, errors.New("role or permission doesn't exist")
	}

	var rolePermissionMapping roles.RolePermissions
	if err := p.session.Where("role_id = ? AND permission_id = ?", roleID, permissionID).First(&rolePermissionMapping).Error; err != nil {
		return false, err
	}

	return true, nil
}

/**
 * End Role and Permission Management
 */
