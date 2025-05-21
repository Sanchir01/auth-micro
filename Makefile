SILENT:
PHONY:


PHONY:
SILENT:
include .env.prod
export
MIGRATION_NAME ?= new_migration

DB_CONN_PROD = host=$(DB_HOST_PROD) user=$(DB_USER_PROD) password=$(DB_PASSWORD_PROD) port=$(DB_PORT_PROD) dbname=$(DB_NAME_PROD) sslmode=disable


generate-dataloaders:
	(cd internal/feature/color  && dataloaden  LoaderByID LoaderByIdColor \*github.com/Sanchir01/colors/pkg/feature/color.Color)

gql:
	go get github.com/99designs/gqlgen@latest && go run github.com/99designs/gqlgen generate
build:
	go build -o ./.bin/main ./cmd/main/main.go
run: build
	ENV_FILE=".env.prod" ./.bin/main

migrations-up:
	goose -dir migrations postgres "host=localhost user=postgres password=postgres port=5435 dbname=auth-micro sslmode=disable"  up

migrations-down:
	goose -dir migrations postgres  "host=localhost user=postgres password=postgres port=5435 dbname=auth-micro sslmode=disable"  down


migrations-status:
	goose -dir migrations postgres  "host=localhost user=postgres password=postgres port=5435 dbname=auth-micro sslmode=disable" status

migrations-new:
	goose -dir migrations create $(MIGRATION_NAME) sql

migrations-up-prod:
	goose -dir migrations postgres "$(DB_CONN_PROD)" up

migrations-down-prod:
	goose -dir migrations postgres "$(DB_CONN_PROD)" down

migrations-status-prod:
	goose -dir migrations postgres "$(DB_CONN_PROD)" status

docker-build:
	docker build -t candles .

docker:
	docker-compose  up -d

docker-app: docker-build docker

seed:
	go run cmd/seed/main.go

compose-prod:
	docker compose -f docker-compose.prod.yaml up --build -d

