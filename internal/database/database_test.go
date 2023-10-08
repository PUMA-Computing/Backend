package database

import (
	"Backend/configs"
	"testing"
)

func TestInit(t *testing.T) {
	mockDB := NewMockDatabase()

	config := &configs.Config{
		DBUser:     "pumadev",
		DBPassword: "pumadev2023",
		DBHost:     "139.59.116.226",
		DBPort:     "5432",
		DBName:     "testpuma",
	}

	mockDB.On("Init", config).Return(nil)
	mockDB.Init(config)
	mockDB.AssertExpectations(t)
}
