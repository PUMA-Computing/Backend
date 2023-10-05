migration_up :
	migrate -path ./internal/migrations -database postgres://computing:computing2023@139.59.116.226:5432/puma?sslmode=disable up

migration_down :
	migrate -path ./internal/migrations -database postgres://computing:computing2023@139.59.116.226:5432/puma?sslmode=disable down


server:
	go run cmd/app/main.go