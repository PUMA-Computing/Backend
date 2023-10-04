package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"context"
)

func ListPermission() ([]*models.Permission, error) {
	var permissions []*models.Permission
	rows, err := database.DB.Query(context.Background(), `
		SELECT id, name
		FROM permissions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var permission models.Permission
		err := rows.Scan(&permission.ID, &permission.Name)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}
	return permissions, nil
}

func AssignPermissionsToRole(roleID int, permissionIDs []int) error {
	ctx := context.Background()
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				return
			}
		}
	}()

	for _, permissionID := range permissionIDs {
		_, err := tx.Exec(ctx, `
            INSERT INTO role_permissions (role_id, permission_id)
            VALUES ($1, $2)`,
			roleID, permissionID)
		if err != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
