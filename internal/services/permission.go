package services

import (
	"Backend/internal/database"
	"Backend/internal/models"
)

type PermissionService struct {
}

func NewPermissionService() *PermissionService {
	return &PermissionService{}
}

func (ps *PermissionService) ListPermission() ([]*models.Permission, error) {
	permissions, err := database.ListPermission()
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (ps *PermissionService) AssignPermissionToRole(roleID int, permissionIDs []int) error {
	if err := database.AssignPermissionsToRole(roleID, permissionIDs); err != nil {
		return err
	}

	return nil
}
