package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kanatsanan6/go-test/utils"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type SignUpTxParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CompanyID int64     `json:"company_id"`
	CreatedAt time.Time `json:"created_at"`
}

type SignUpTxResult struct {
	User UserResponse `json:"user"`
}

func (store *Store) SignUpTx(ctx context.Context, args SignUpTxParams) (SignUpTxResult, error) {
	var result SignUpTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		company, err := q.CreateCompany(ctx, sql.NullString{})
		if err != nil {
			return err
		}

		hashPassword, err := utils.GeneratePassword([]byte(args.Password))
		if err != nil {
			return err
		}

		user, err := q.CreateUser(ctx, CreateUserParams{
			FirstName:         args.FirstName,
			LastName:          args.LastName,
			Email:             args.Email,
			EncryptedPassword: hashPassword,
			CompanyID:         company.ID,
		})
		if err != nil {
			return err
		}

		result.User = UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CompanyID: company.ID,
			CreatedAt: user.CreatedAt,
		}

		return nil
	})

	return result, err
}
