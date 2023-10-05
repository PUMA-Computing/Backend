package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"context"
	"github.com/google/uuid"
)

func CreateRole(role *models.Roles) error {
	_, err := database.DB.Exec(context.Background(), `
		INSERT INTO roles (name, created_at, updated_at) 
		VALUES ($1, $2, $3)`,
		role.Name, role.CreatedAt, role.UpdatedAt)
	return err
}

func UpdateRole(roleID int, updatedRole *models.Roles) error {
	_, err := database.DB.Exec(context.Background(), `
		UPDATE roles SET name = $1, updated_at = $2
		WHERE id = $3`,
		updatedRole.Name, updatedRole.UpdatedAt, roleID)
	return err
}

func DeleteRole(roleID int) error {
	_, err := database.DB.Exec(context.Background(), `
		DELETE FROM roles WHERE id = $1`, roleID)
	return err
}

func GetRoleByID(roleID int) (*models.Roles, error) {
	var role models.Roles
	err := database.DB.QueryRow(context.Background(), `
		SELECT id, name, updated_at, created_at
		FROM roles WHERE id = $1`, roleID).Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func ListRoles() ([]*models.Roles, error) {
	var roles []*models.Roles
	rows, err := database.DB.Query(context.Background(), `
		SELECT id, name, updated_at, created_at
		FROM roles`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var role models.Roles
		err := rows.Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}
	return roles, nil
}

func AssignRoleToUser(userID uuid.UUID, roleID int) error {
	// Update the RoleID of the user on table users
	_, err := database.DB.Exec(context.Background(), `
		UPDATE users SET role_id = $1
		WHERE id = $2`,
		roleID, userID)
	return err
}
