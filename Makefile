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

install-tools:
	go install github.com/cosmtrek/air@latest
