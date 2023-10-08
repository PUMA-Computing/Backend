package database

import (
	"Backend/configs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Init(config *configs.Config) {
	m.Called(config)
}

func (m *MockDatabase) Close() {
	m.Called()
}

func (m *MockDatabase) GetDB() *pgxpool.Pool {
	args := m.Called()
	return args.Get(0).(*pgxpool.Pool)
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{}
}
