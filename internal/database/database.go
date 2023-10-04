package database

import (
	"Backend/configs"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var DB *pgxpool.Pool

func Init(config *configs.Config) {
	var err error

	connStr := "user=" + config.DBUser +
		" password=" + config.DBPassword +
		" host=" + config.DBHost +
		" port=" + config.DBPort +
		" dbname=" + config.DBName +
		" sslmode=disable"

	DB, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Connected to the database")
}

func Close() {
	DB.Close()
}

func GetDB() *pgxpool.Pool {
	return DB
}
