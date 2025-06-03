APP_NAME=go-hrms
MIGRATIONS_DIR=./database/migrations
DB_URL=postgres://:@localhost:5432/hrms?sslmode=disable

.PHONY: run migrate-up migrate-down migrate-create

run:
	go run .

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)