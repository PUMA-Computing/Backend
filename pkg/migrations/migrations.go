package migrations

import (
	"fmt"
	"github.com/gocql/gocql"
	"os"
	"path/filepath"
)

func ExecuteMigrations(session *gocql.Session) error {
	migrationsDir := "./internal/migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".cql" {
			cqlBytes, err := os.ReadFile(filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				return err
			}

			cqlQuery := string(cqlBytes)
			if err := session.Query(cqlQuery).Exec(); err != nil {
				return err
			}

			fmt.Printf("Executed migration: %s\n", file.Name())
		}
	}

	return nil
}
