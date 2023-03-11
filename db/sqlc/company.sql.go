// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: company.sql

package db

import (
	"context"
	"database/sql"
)

const createCompany = `-- name: CreateCompany :one
INSERT INTO companies (
  name
) VALUES (
  $1
)
RETURNING id, name, created_at
`

func (q *Queries) CreateCompany(ctx context.Context, name sql.NullString) (Company, error) {
	row := q.db.QueryRowContext(ctx, createCompany, name)
	var i Company
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return i, err
}
