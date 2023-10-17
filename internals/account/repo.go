package account

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func findItemByEmail(
	ctx context.Context,
	tx pgx.Tx,
	email string,
) (
	Account,
	error,
) {
	query := `
		SELECT
			id,
			email,
			password,
			book_default,
			created_at,
			updated_at
		FROM
			Account
		WHERE
			email = $1
	`
	row := tx.QueryRow(
		ctx,
		query,
		email,
	)
	var item Account
	if err := row.Scan(
		&item.Id,
		&item.Email,
		&item.Password,
		&item.DefaultBookId,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
			return Account{}, ErrNotFound
		}
		return Account{}, err
	}
	return item, nil
}

func findItemById(
	ctx context.Context,
	tx pgx.Tx,
	id ulid.ULID,
) (
	Account,
	error,
) {
	query := `
		SELECT
			id,
			email,
			password,
			book_default,
			created_at,
			updated_at
		FROM
			Account
		WHERE
			id = $1
	`
	row := tx.QueryRow(
		ctx,
		query,
		id,
	)
	var item Account
	if err := row.Scan(
		&item.Id,
		&item.Email,
		&item.Password,
		&item.DefaultBookId,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
			return Account{}, ErrNotFound
		}
		return Account{}, err
	}
	return item, nil
}

func saveItem(
	ctx context.Context,
	tx pgx.Tx,
	item *Account,
) error {
	query := `
		INSERT INTO Account(
			id,
			email,
			password,
			book_default,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5
		) ON CONFLICT(id)
			DO UPDATE SET
			email=$2, password=$3, book_default=$4, updated_at=$6
	`

	_, err := tx.Exec(
		ctx,
		query,
		item.Id,
		item.Email,
		item.Password,
		item.CreatedAt,
		item.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
