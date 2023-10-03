package database

import (
	"Backend/internal/models"
	"context"
)

func CreateRole(role *models.Roles) error {
	_, err := DB.Exec(context.Background(), `
		INSERT INTO roles (name) 
		VALUES ($1)`,
		role.Name)
	return err
}

func UpdateRole(roleID int, updatedRole *models.Roles) error {
	_, err := DB.Exec(context.Background(), `
		UPDATE roles SET name = $1
		WHERE id = $2`,
		updatedRole.Name, roleID)
	return err
}

func DeleteRole(roleID int) error {
	_, err := DB.Exec(context.Background(), `
		DELETE FROM roles WHERE id = $1`, roleID)
	return err
}

func GetRoleByID(roleID int) (*models.Roles, error) {
	var role models.Roles
	err := DB.QueryRow(context.Background(), `
		SELECT id, name
		FROM roles WHERE id = $1`, roleID).Scan(&role.ID, &role.Name)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func ListRoles() ([]*models.Roles, error) {
	var roles []*models.Roles
	rows, err := DB.Query(context.Background(), `
		SELECT id, name
		FROM roles`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var role models.Roles
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}
	return roles, nil
}

func AssignRoleToUser(userID int, roleID int) error {
	_, err := DB.Exec(context.Background(), `
		INSERT INTO user_roles (user_id, role_id) 
		VALUES ($1, $2)`,
		userID, roleID)
	return err
}
