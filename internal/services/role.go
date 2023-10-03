package services

import (
	"Backend/internal/database"
	"Backend/internal/models"
)

type RoleService struct {
}

func NewRoleService() *RoleService {
	return &RoleService{}
}

func (rs *RoleService) CreateRole(role *models.Roles) error {
	if err := database.CreateRole(role); err != nil {
		return err
	}

	return nil
}

func (rs *RoleService) EditRole(roleID int, updatedRole *models.Roles) error {
	if err := database.UpdateRole(roleID, updatedRole); err != nil {
		return err
	}

	return nil
}

func (rs *RoleService) DeleteRole(roleID int) error {
	if err := database.DeleteRole(roleID); err != nil {
		return err
	}

	return nil
}

func (rs *RoleService) GetRoleByID(roleID int) (*models.Roles, error) {
	role, err := database.GetRoleByID(roleID)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (rs *RoleService) ListRoles() ([]*models.Roles, error) {
	roles, err := database.ListRoles()
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (rs *RoleService) AssignRoleToUser(userID int, roleID int) error {
	if err := database.AssignRoleToUser(userID, roleID); err != nil {
		return err
	}

	return nil
}
