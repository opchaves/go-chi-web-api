package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/go-chi-web-api/internal/app"
	"github.com/opchaves/go-chi-web-api/internal/config"
	"github.com/opchaves/go-chi-web-api/internal/model"
	"github.com/opchaves/go-chi-web-api/internal/server"
)

func main() {
	ctx := context.TODO()

	db, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		log.Fatalln("Could not connect to database", err)
	}
	defer db.Close()

	log.Println("Database successfully connected.")

	query := model.New(db)

	// TODO init services with `query`
	// TODO add services, `query` and `db` to `server` to be available in all handlers

	users, err := query.GetUsers(ctx)
	if err != nil {
		log.Fatalln("Could not get users", err)
	}
	log.Println("Users:", users)

	s := server.New(
		config.Name,
		server.UseHost(config.Host),
		server.UsePort(config.Port),
	)

	if err := app.AddRoutes(s); err != nil {
		log.Fatalf("error adding routes: %v", err)
	}

	if err := s.Run(); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
