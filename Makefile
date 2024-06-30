include .env

migration_up :
	migrate -path ./migrations -database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable up

migration_down :
	migrate -path ./migrations -database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable down

migration_create :
	migrate create -ext sql -dir ./migrations -seq $(name) && cat ./internal/models/${name}.go

server:
	go run cmd/app/main.go