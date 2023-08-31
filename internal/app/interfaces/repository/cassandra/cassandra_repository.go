package cassandra

import (
	"github.com/gocql/gocql"
	"time"
)

type SessionData struct {
	UserID       gocql.UUID
	SessionToken string
	ExpiredAt    time.Time
}

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

func (r *CassandraRepository) GetSessionData(userID gocql.UUID) (*SessionData, error) {
	query := `
		SELECT user_id, session_token, expired_at
		FROM sessions
		WHERE user_id = ?
		LIMIT 1
	`
	var sessionData SessionData
	if err := r.Session.Query(query, userID).Scan(&sessionData.UserID, &sessionData.SessionToken, &sessionData.ExpiredAt); err != nil {
		return nil, err
	}
	return &sessionData, nil
}
