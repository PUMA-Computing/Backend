package repository

import (
	"github.com/gocql/gocql"
	"time"
)

type CassandraRepository struct {
	Session *gocql.Session
}

func NewCassandraRepository() (*CassandraRepository, error) {
	cluster := gocql.NewCluster("139.59.116.226")
	cluster.Keyspace = "puma"
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &CassandraRepository{Session: session}, nil
}

func (r *CassandraRepository) Close() {
	r.Session.Close()
}

func (r *CassandraRepository) StoreSessionData(userID gocql.UUID, sessionToken string, expirationTime time.Time) error {
	query := `
		INSERT INTO sessions (user_id, session_token, expired_at)
		VALUES (?, ?, ?)
	`
	if err := r.Session.Query(query, userID, sessionToken, expirationTime).Exec(); err != nil {
		return err
	}
	return nil
}
