SHORT_ID := $(shell git rev-parse --short HEAD)

export V_TAG = ${SHORT_ID}

# This assumes there is a .env file in the root of the project
define setup_env
  $(eval ENV_FILE := .env)
  $(eval include .env)
  $(eval export)
endef

withEnv: 
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

docker-up: withEnv
	@docker compose up --build -d

docker-down: withEnv
	@docker compose down

docker-build: withEnv
	@docker build --target prod -t ${IMAGE_NAME}:${V_TAG} .

docker-run: withEnv
	@docker run --rm --name ${APP_NAME} -p 8080:8080 ${IMAGE_NAME}:${V_TAG}

db-sh: withEnv
	@docker compose exec postgres psql -U postgres -d app_dev

db-schema-dump: withEnv
	@docker compose exec postgres pg_dump -U postgres --schema-only -d app_dev > db/schema.sql

seed: withEnv
	@go run ./cmd/seed/main.go

migrate: withEnv
	@if command -v migrate > /dev/null; then \
	  migrate -database ${DATABASE_URL} -path ./db/migrations up;\
	else \
	  echo "Installing 'migrate' tool...";\
	  make install-tools;\
	  echo "Running migrate up...";\
	  migrate -database ${DATABASE_URL} -path ./db/migrations up;\
	fi

migrate-new:
	@if command -v migrate > /dev/null; then \
	  migrate create -ext sql -dir db/migrations -seq $(name);\
	else \
	  echo "Installing 'migrate' tool...";\
	  make install-tools; \
	  echo "Creating new migration '$(name)'...";\
	  migrate create -ext sql -dir db/migrations -seq $(name);\
	fi

migrate-down: withEnv
	@migrate -database ${DATABASE_URL} -path ./db/migrations down

migrate-drop: withEnv
	@migrate -database ${DATABASE_URL} -path ./db/migrations drop

install-tools:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Live Reload
watch-api:
	@if command -v air > /dev/null; then \
	  air; \
	  echo "Watching...";\
	else \
	  read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	  if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	    go install github.com/cosmtrek/air@latest; \
	    air; \
	    echo "Watching...";\
	  else \
	    echo "You chose not to install air. Exiting..."; \
	    exit 1; \
	  fi; \
	fi

# Generate Go code from SQL
sqlc:
	@if command -v sqlc > /dev/null; then \
	  echo "Gerating Go code...";\
	  sqlc generate; \
	else \
	  read -p "Go's 'sqlc' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	  if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest; \
	  	echo "Gerating Go code...";\
	  	sqlc generate; \
	  else \
	    echo "You chose not to install sqlc. Exiting..."; \
	    exit 1; \
	  fi; \
	fi
