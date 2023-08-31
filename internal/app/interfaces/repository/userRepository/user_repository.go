package userRepository

import (
	"Backend/internal/app/domain/user"
	"github.com/gocql/gocql"
)

type UserRepository interface {
	RegisterUser(user *user.User) error
	GetUserByEmail(email string) (*user.User, error)
}

type CassandraUserRepository struct {
	session *gocql.Session
}

func NewCassandraUserRepository(session *gocql.Session) *CassandraUserRepository {
	return &CassandraUserRepository{session: session}
}

func (r *CassandraUserRepository) RegisterUser(user *user.User) error {
	query := r.session.Query(
		"INSERT INTO users (id, first_name, last_name, email, password, nim, year, role) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.NIM, user.Year, user.Role,
	)

	return query.Exec()
}

func (r *CassandraUserRepository) GetUserByEmail(email string) (*user.User, error) {
	query := `SELECT id, first_name, last_name, email, password, nim, year, role FROM users WHERE email = ? ALLOW FILTERING`
	var user user.User
	if err := r.session.Query(query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.NIM, &user.Year, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}
