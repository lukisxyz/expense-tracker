package category

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool

	ErrAssignPool = errors.New("cannot assign nil pool")
	ErrNotFound   = errors.New("category: item not found")
)

func SetPool(newPool *pgxpool.Pool) error {
	if newPool == nil {
		return ErrAssignPool
	}
	pool = newPool
	return nil
}
