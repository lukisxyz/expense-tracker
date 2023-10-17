package category

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func saveItem(
	ctx context.Context,
	tx pgx.Tx,
	item *Category,
) error {
	query := `
		INSERT INTO category(
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

func findItemById(
	ctx context.Context,
	tx pgx.Tx,
	id ulid.ULID,
) (
	Category,
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
			category
		WHERE
			id = $1
	`
	row := tx.QueryRow(
		ctx,
		query,
		id,
	)
	var item Category
	if err := row.Scan(
		&item.Id,
		&item.Name,
		&item.OwnerId,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
			return Category{}, ErrNotFound
		}
		return Category{}, err
	}
	return item, nil
}
