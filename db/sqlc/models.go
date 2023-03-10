// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"time"
)

type Company struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}

type User struct {
	ID                int64
	FirstName         string
	LastName          string
	Email             string
	EncryptedPassword string
	CreatedAt         time.Time
	CompanyID         int64
}
