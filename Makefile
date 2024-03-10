SHORT_ID := $(shell git rev-parse --short HEAD)

export V_TAG = ${SHORT_ID}

# This assumes there is a .env file in the root of the project
define setup_env
  $(eval ENV_FILE := .env)
  $(eval include .env)
  $(eval export)
endef

with-env: 
	$(call setup_env)

run:
	go run ./cmd/server/main.go

dev:
	make watch-api & make watch-web

watch-web:
	cd web && npm run dev

tidy:
	go mod tidy

build:
	@make build-web
	@make build-api
	@echo "Done!"

build-api:
	@echo "Building API..."
	@go build -o bin/server ./cmd/server/main.go
	@chmod +x bin/server

build-api-prod:
	@echo "Building production API..."
	@env GOOS=linux GOARCH=amd64 go build -o bin/server ./cmd/server/main.go
	@chmod +x bin/server

build-web:
	@echo "Building Web..."
	@cd ./web && npm install && npm run build

start:
	./bin/server

docker-up: with-env
	@docker compose up --build -d

docker-down: with-env
	@docker compose down

docker-build: with-env
	@docker build --target prod -t ${IMAGE_NAME}:${V_TAG} .

docker-run: with-env
	@docker run --rm --name ${APP_NAME} -p 8080:8080 ${IMAGE_NAME}:${V_TAG}

psql: with-env
	@docker compose exec postgres psql -U postgres -d app_dev

db-schema-dump: with-env
	@docker compose exec postgres pg_dump -U postgres --schema-only -d app_dev > db/schema.sql

seed: with-env
	@go run ./cmd/seed/main.go

migrate: with-env
	@make migrate-dev
	@make migrate-test

migrate-dev: with-env
	@migrate -database ${DATABASE_URL} -path ./db/migrations up;\

migrate-test: with-env
	@migrate -database ${DATABASE_URL_TEST} -path ./db/migrations up;\

migrate-new:
	migrate create -ext sql -dir db/migrations -seq $(name);\

migrate-down: with-env
	@migrate -database ${DATABASE_URL} -path ./db/migrations down

migrate-drop: with-env
	@migrate -database ${DATABASE_URL} -path ./db/migrations drop

install-tools:
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/cosmtrek/air@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Live Reload
watch-api:
	air

# Generate Go code from SQL
sqlc:
	sqlc generate
