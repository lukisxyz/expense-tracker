package record

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func saveItem(
	ctx context.Context,
	tx pgx.Tx,
	item *Record,
) error {
	query := `
		INSERT INTO record(
			id,
			note,
			amount,
			is_expense,
			book_id,
			category_id,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) ON CONFLICT(id)
			DO UPDATE SET
			note=$2, amount=$3, is_expense=$4, book_id=$5, category_id=$6, updated_at=$8
	`

	_, err := tx.Exec(
		ctx,
		query,
		item.Id,
		item.Note,
		item.Amount,
		item.IsExpense,
		item.Book.Id,
		item.Category.Id,
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
	Record,
	error,
) {
	query := `
		SELECT
			r.id AS record_id,
			r.note AS record_note,
			r.amount AS record_amount,
			r.is_expense AS record_is_expense,
			r.created_at AS record_created_at,
			r.updated_at AS record_updated_at,
			c.id AS category_id,
			c.name AS category_name,
			b.id AS book_id,
			b.name AS book_name
		FROM
			public.record r
		JOIN
			public.category c ON r.category_id = c.id
		JOIN
			public.book b ON r.book_id = b.id
		WHERE
			r.id = $1;
	`
	row := tx.QueryRow(
		ctx,
		query,
		id,
	)
	var item Record
	if err := row.Scan(
		&item.Id,
		&item.Note,
		&item.Amount,
		&item.IsExpense,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.Book.Id,
		&item.Book.Name,
		&item.Category.Id,
		&item.Category.Name,
	); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
			return Record{}, ErrNotFound
		}
		return Record{}, err
	}
	return item, nil
}
