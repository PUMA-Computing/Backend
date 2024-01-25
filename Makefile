migration_up :
	migrate -path ./migrations -database postgres://pufacomputing:qwertyuiop@139.59.116.226:5432/pufadbtest?sslmode=disable up

migration_down :
	migrate -path ./migrations -database postgres://pufacomputing:qwertyuiop@139.59.116.226:5432/pufadbtest?sslmode=disable down

migration_create :
	migrate create -ext sql -dir ./migrations -seq $(name)

server:
	go run cmd/app/main.go