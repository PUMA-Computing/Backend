package migrations

import (
	"fmt"
	"gorm.io/gorm"
	"os"
)

func ExecuteMigrations(session *gorm.DB) error {
	migrationsDir := "./internal/migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := session.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(fmt.Sprintf("SELECT 1 FROM %s LIMIT 1", file.Name())).Error; err != nil {
				if err := tx.Exec(fmt.Sprintf("CREATE TABLE %s ()", file.Name())).Error; err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}
