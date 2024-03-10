package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/go-kom/config"
	"github.com/opchaves/go-kom/internal/model"
	"github.com/opchaves/go-kom/pkg/password"
)

func main() {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, config.DatabaseURL)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer db.Close()
	query := model.New(db)

	tablesToDelete := []string{
		"transactions",
		"categories",
		"accounts",
		"workspaces_users",
		"workspaces",
		"tokens",
		"users",
	}

	for _, t := range tablesToDelete {
		_, err = db.Exec(ctx, fmt.Sprintf("DELETE FROM %s", t))
		if err != nil {
			log.Fatalf("Unable to delete table %v. err: %v\n", t, err)
		}
	}

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

	workspace1, err := query.CreateWorkspace(ctx, workspaceParams)
	if err != nil {
		log.Fatalf("Unable to create workspace: %v\n", err)
	}

	workspace2, err := query.CreateWorkspace(ctx, workspace2Params)
	if err != nil {
		log.Fatalf("Unable to create workspace: %v\n", err)
	}

	workspaceUser1Params := model.CreateWorkspaceUserParams{
		UserID:      user.ID,
		WorkspaceID: workspace1.ID,
		Role:        "admin",
	}

	workspaceUser2Params := model.CreateWorkspaceUserParams{
		UserID:      user2.ID,
		WorkspaceID: workspace2.ID,
		Role:        "admin",
	}

	_, err = query.CreateWorkspaceUser(ctx, workspaceUser1Params)
	if err != nil {
		log.Fatalf("Unable to create workspace user: %v\n", err)
	}

	_, err = query.CreateWorkspaceUser(ctx, workspaceUser2Params)
	if err != nil {
		log.Fatalf("Unable to create workspace user 2: %v\n", err)
	}

	log.Println("Seeding completed successfully!")
}
