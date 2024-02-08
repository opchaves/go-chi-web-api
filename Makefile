dev:
	air

run:
	go run ./cmd/server/main.go

tidy:
	go mod tidy

build:
	go build -o bin/server ./cmd/server/main.go
	chmod +x bin/server

start:
	./bin/server

db-start:
	docker compose up -d

db-remove:
	docker compose down --volumes --remove-orphans

db-sh:
	docker compose exec postgres psql -U opchaves -d app_dev

install-tools:
	go install github.com/cosmtrek/air@latest
