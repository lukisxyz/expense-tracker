package book

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func findItemById(
	ctx context.Context,
	tx pgx.Tx,
	id ulid.ULID,
) (
	Book,
	error,
) {
	query := `
		SELECT
			id,
			name,
			owner_id,
			created_at,
			updated_at
		FROM
			book
		WHERE
			id = $1
	`
	row := tx.QueryRow(
		ctx,
		query,
		id,
	)
	var item Book
	if err := row.Scan(
		&item.Id,
		&item.Name,
		&item.OwnerId,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
			return Book{}, ErrNotFound
		}
		return Book{}, err
	}
	return item, nil
}

func findDefaultItemById(
	ctx context.Context,
	tx pgx.Tx,
	ownerId ulid.ULID,
) (
	Book,
	error,
) {
	query := `
		SELECT
			id,
			name,
			owner_id,
			created_at,
			updated_at
		FROM
			book
		WHERE
			owner_id = $1 AND is_default = true
	`
	row := tx.QueryRow(
		ctx,
		query,
		ownerId,
	)
	var item Book
	if err := row.Scan(
		&item.Id,
		&item.Name,
		&item.OwnerId,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
			return Book{}, ErrNotFound
		}
		return Book{}, err
	}
	return item, nil
}

func saveItem(
	ctx context.Context,
	tx pgx.Tx,
	item *Book,
) error {
	query := `
		INSERT INTO book(
			id,
			name,
			owner_id,
			created_at
		) VALUES (
			$1, $2, $3, $4
		) ON CONFLICT(id)
			DO UPDATE SET
			name=$2, owner_id=$3, updated_at=$5
	`

	_, err := tx.Exec(
		ctx,
		query,
		item.Id,
		item.Name,
		item.OwnerId,
		item.CreatedAt,
		item.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
