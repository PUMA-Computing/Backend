package eventRepository_test

import (
	"Backend/internal/app/interfaces/repository/eventRepository"
	"Backend/internal/app/interfaces/repository/postgresRepository"
	"testing"
)

var db, _ = postgresRepository.NewPostgresRepository()
var eventRepo eventRepository.EventRepository

func TestMain(m *testing.M) {
	setupTestDB()
	defer cleanupTestDB()

	m.Run()
}

func setupTestDB() {

	// Migrate the database schema if needed.
	// Your migration code here...

	eventRepo = eventRepository.NewPostgresForEventRepository(db.DB)
}

func cleanupTestDB() {
	// Clean up the database after running all tests.
	// Your clean up code here...

	db.Close()
}
