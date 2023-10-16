package book

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool

	ErrAssignPool = errors.New("cannot assign nil pool")

	ErrNotFound      = errors.New("book: item not found")
	ErrOwnerNotFound = errors.New("book: owner not found")
)

func SetPool(newPool *pgxpool.Pool) error {
	if newPool == nil {
		return ErrAssignPool
	}
	pool = newPool
	return nil
}