define setup_env
  $(eval ENV_FILE := .env)
  $(eval include .env)
  $(eval export)
endef

devEnv: 
	$(call setup_env, dev)

run:
	go run ./cmd/server/main.go

run-web:
	cd web && npm run dev

tidy:
	go mod tidy

build:
	go build -o bin/server ./cmd/server/main.go
	chmod +x bin/server

build-web:
	cd ./web && npm install && npm run build

start:
	./bin/server

docker-run:
	docker compose up --build -d

docker-down:
	docker compose down

db-sh:
	@docker compose exec postgres psql -U postgres -d app_dev

migrate: devEnv
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

migrate-down: devEnv
	@migrate -database ${DATABASE_URL} -path ./db/migrations down

install-tools:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Live Reload
dev:
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
			go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest \
	  	echo "Gerating Go code...";\
	  	sqlc generate; \
	  else \
	    echo "You chose not to install sqlc. Exiting..."; \
	    exit 1; \
	  fi; \
	fi
