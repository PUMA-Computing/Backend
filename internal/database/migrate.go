package database

import (
	"Backend/configs"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func Migrate() {

	config := configs.LoadConfig()

	connStr := "user=" + config.DBUser +
		" password=" + config.DBPassword +
		" host=" + config.DBHost +
		" port=" + config.DBPort +
		" dbname=" + config.DBName +
		" sslmode=disable"

	fmt.Println("Connection String:", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error creating temporary database connection: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing temporary database connection: %v", err)
		}
	}(db) // Close the temporary *sql.DB when done with migrations

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Error creating migration driver: %v", err)
	}

	fileSrc, err := (&file.File{}).Open("migrations")
	if err != nil {
		log.Fatalf("Error opening migration source: %v", err)
	}

	m, err := migrate.NewWithInstance("file", fileSrc, config.DBName, driver)
	if err != nil {
		log.Fatalf("Error creating migration instance: %v", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No Migration to Apply")
			return
		}
		log.Fatalf("Error applying migration: %v", err)
	}
	log.Println("Migration Applied Successfully")

}
