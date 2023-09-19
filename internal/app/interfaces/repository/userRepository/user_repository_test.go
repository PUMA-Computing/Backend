package userRepository_test

import (
	"Backend/internal/app/interfaces/repository/postgresRepository"
	"Backend/internal/app/interfaces/repository/userRepository"
	"testing"
)

var db, _ = postgresRepository.NewPostgresRepository()
var userRepo userRepository.UserRepository

func TestMain(m *testing.M) {
	setupTestDB()
	defer cleanupTestDB()

	m.Run()
}

func setupTestDB() {

	// Migrate the database schema if needed.
	// Your migration code here...

	userRepo = userRepository.NewPostgresUserRepository(db.DB)
}

func cleanupTestDB() {
	// Clean up the database after running all tests.
	// Your clean up code here...

	db.Close()
}

func TestPostgresUserRepository_AuthenticateUser(t *testing.T) {
	t.Run("AuthenticateUser - Existing User", func(t *testing.T) {
		_, err := userRepo.AuthenticateUser("testuser@gmail.com", "12345")
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
	})
}

func TestPostgresUserRepository_RoleExists(t *testing.T) {
	t.Run("RoleExists - Existing Role", func(t *testing.T) {
		exists, err := userRepo.RoleExists(1) // Replace with an actual role ID that exists in your test database.
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if !exists {
			t.Errorf("Expected role to exist, but it does not")
		}
	})

	t.Run("RoleExists - Non-Existent Role", func(t *testing.T) {
		exists, err := userRepo.RoleExists(9) // Replace with a role ID that does not exist in your test database.
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if exists {
			t.Errorf("Expected role not to exist, but it does")
		}
	})
}

func TestPostgresUserRepository_HasPermission(t *testing.T) {
	t.Run("HasPermission - Existing Role and Permission", func(t *testing.T) {
		hasPermission, err := userRepo.HasPermission(1, 1) // Replace with an actual role ID and permission ID that exists in your test database.
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if !hasPermission {
			t.Errorf("Expected role to have permission, but it does not")
		}
	})
}
