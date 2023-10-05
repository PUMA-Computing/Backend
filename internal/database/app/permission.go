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
	_, err := database.DB.Exec(context.Background(), `
		DELETE FROM role_permissions
		WHERE role_id = $1`, roleID)
	if err != nil {
		return err
	}

	for _, permissionID := range permissionIDs {
		_, err := database.DB.Exec(context.Background(), `
			INSERT INTO role_permissions (role_id, permission_id)
			VALUES ($1, $2)`, roleID, permissionID)
		if err != nil {
			return err
		}
	}

	return nil
}
