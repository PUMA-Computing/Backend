package repository

import (
	"Backend/internal/app/domain"
	"github.com/gocql/gocql"
)

type CassandraUserRepository struct {
	session *gocql.Session
}

func NewCassandraUserRepository(session *gocql.Session) *CassandraUserRepository {
	return &CassandraUserRepository{session: session}
}

func (r *CassandraUserRepository) CreateUser(user *domain.User) error {
	query := `INSERT INTO users (id, first_name, last_name, email, password, nim, year, role) VALUES (?, ?, ?, ?, ?, ?)`
	err := r.session.Query(query, user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.NIM, user.Year, user.Role).Exec()
	if err != nil {
		return err
	}
	return err
}

func (r *CassandraUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	query := `SELECT id, first_name, last_name, email, password, nim, year, role FROM users WHERE email = ?`
	var user domain.User
	if err := r.session.Query(query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.NIM, &user.Year, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}
