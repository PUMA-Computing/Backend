package services

import (
	"Backend/internal/database"
	"Backend/internal/database/app"
	"Backend/internal/models"
	"context"
	"errors"
	"github.com/google/uuid"
)

type PermissionService struct {
}

func NewPermissionService() *PermissionService {
	return &PermissionService{}
}

func (ps *PermissionService) ListPermission() ([]*models.Permission, error) {
	permissions, err := app.ListPermission()
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (ps *PermissionService) AssignPermissionToRole(roleID int, permissionIDs []int) error {
	existingRole, err := app.GetRoleByID(roleID)
	if err != nil {
		return err
	}

	if existingRole == nil {
		return errors.New("role not found")
	}

	if err := app.AssignPermissionsToRole(roleID, permissionIDs); err != nil {
		return err
	}

	return nil
}

func (ps *PermissionService) CheckPermission(ctx context.Context, userID uuid.UUID, requiredPermission string) (bool, error) {
	return database.CheckPermission(ctx, userID, requiredPermission)
}
