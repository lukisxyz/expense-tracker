package account

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
	Account,
	error,
) {
	query := `
		SELECT
			id,
			email,
			password,
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
			created_at
		) VALUES (
			$1, $2, $3, $4
		) ON CONFLICT(id)
			DO UPDATE SET
			email=$2, password=$3, updated_at=$5
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
