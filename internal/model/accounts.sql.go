// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: accounts.sql

package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (
  name, balance, user_id, workspace_id
) VALUES (
  $1, $2, $3, $4
) returning id, name, balance, user_id, workspace_id, created_at, updated_at
`

type CreateAccountParams struct {
	Name        string         `json:"name"`
	Balance     pgtype.Numeric `json:"balance"`
	UserID      uuid.UUID      `json:"user_id"`
	WorkspaceID uuid.UUID      `json:"workspace_id"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (*Account, error) {
	row := q.db.QueryRow(ctx, createAccount,
		arg.Name,
		arg.Balance,
		arg.UserID,
		arg.WorkspaceID,
	)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Balance,
		&i.UserID,
		&i.WorkspaceID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteAccount, id)
	return err
}

const getAccountByID = `-- name: GetAccountByID :one
SELECT id, name, balance, user_id, workspace_id, created_at, updated_at FROM accounts WHERE id = $1
`

func (q *Queries) GetAccountByID(ctx context.Context, id uuid.UUID) (*Account, error) {
	row := q.db.QueryRow(ctx, getAccountByID, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Balance,
		&i.UserID,
		&i.WorkspaceID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getAccountsByWorkspace = `-- name: GetAccountsByWorkspace :many
SELECT id, name, balance, user_id, workspace_id, created_at, updated_at FROM accounts WHERE workspace_id = $1
`

func (q *Queries) GetAccountsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]*Account, error) {
	rows, err := q.db.Query(ctx, getAccountsByWorkspace, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Balance,
			&i.UserID,
			&i.WorkspaceID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccount = `-- name: UpdateAccount :exec
UPDATE accounts SET
  name = $1,
  balance = coalesce($2, balance)
WHERE id = $3
`

type UpdateAccountParams struct {
	Name    string         `json:"name"`
	Balance pgtype.Numeric `json:"balance"`
	ID      uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) error {
	_, err := q.db.Exec(ctx, updateAccount, arg.Name, arg.Balance, arg.ID)
	return err
}
