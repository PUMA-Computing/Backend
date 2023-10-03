package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
)

func CheckPermission(ctx context.Context, userID uuid.UUID, requiredPermission string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM users u
        JOIN user_roles ur ON u.id = ur.user_id
        JOIN roles r ON ur.role_id = r.id
        JOIN role_permissions rp ON r.id = rp.role_id
        JOIN permissions p ON rp.permission_id = p.id
        WHERE u.id = $1 AND p.name = $2`

	var count int
	err := DB.QueryRow(ctx, query, userID, requiredPermission).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return count > 0, nil
}
