package account

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool

	ErrAssignPool    = errors.New("cannot assign nil pool")
	ErrNotFound      = errors.New("account: item not found")
	ErrWrongPassword = errors.New("account: password wrong")
)

func SetPool(newPool *pgxpool.Pool) error {
	if newPool == nil {
		return ErrAssignPool
	}
	pool = newPool
	return nil
}
