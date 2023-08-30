package repository

import "github.com/gocql/gocql"

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
