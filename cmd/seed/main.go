package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/go-chi-web-api/internal/config"
	"github.com/opchaves/go-chi-web-api/internal/model"
	"github.com/opchaves/go-chi-web-api/pkg/password"
)

// TODO maybe use faker to generate random data

func main() {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, config.DatabaseURL)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer db.Close()
	query := model.New(db)

	password, _ := password.Hash("password12")

	userParams := model.CreateUserParams{
		FirstName: "Paulo",
		LastName:  "Chaves",
		Email:     "paulo@example.com",
		Password:  password,
		Verified:  true,
		Avatar:    "https://example.com/avatar.jpg",
	}

	user2Params := model.CreateUserParams{
		FirstName: "Jonh",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  password,
		Verified:  true,
		Avatar:    "https://example.com/avatar.jpg",
	}

	_, err = db.Exec(ctx, "DELETE FROM workspaces")
	if err != nil {
		log.Fatalf("Unable to delete workspaces: %v\n", err)
	}

	_, err = db.Exec(ctx, "DELETE FROM tokens")
	if err != nil {
		log.Fatalf("Unable to delete tokens: %v\n", err)
	}

	_, err = db.Exec(ctx, "DELETE FROM users")
	if err != nil {
		log.Fatalf("Unable to delete users: %v\n", err)
	}

	user, err := query.CreateUser(ctx, userParams)
	if err != nil {
		log.Fatalf("Unable to create user: %v\n", err)
	}

	user2, err := query.CreateUser(ctx, user2Params)
	if err != nil {
		log.Fatalf("Unable to create user 2: %v\n", err)
	}

	description := "My personal workspace"
	workspaceParams := model.CreateWorkspaceParams{
		Name:        "My Workspace",
		UserID:      user.ID,
		Currency:    "USD",
		Language:    "en",
		Description: &description,
	}

	workspace2Params := model.CreateWorkspaceParams{
		Name:        "My Workspace",
		UserID:      user2.ID,
		Currency:    "USD",
		Language:    "en",
		Description: &description,
	}

	_, err = query.CreateWorkspace(ctx, workspaceParams)
	if err != nil {
		log.Fatalf("Unable to create workspace: %v\n", err)
	}

	_, err = query.CreateWorkspace(ctx, workspace2Params)
	if err != nil {
		log.Fatalf("Unable to create workspace: %v\n", err)
	}

	log.Println("Seeding completed successfully!")
}
