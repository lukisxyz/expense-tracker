package record

import (
	"context"

	"github.com/jackc/pgx/v5"
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

// func findItemById(
// 	ctx context.Context,
// 	tx pgx.Tx,
// 	id ulid.ULID,
// ) (
// 	Record,
// 	error,
// ) {
// 	query := `
// 		SELECT
// 			id,
// 			note,
// 			amount,
// 			is_expense,
// 			book_id,
// 			category_id,
// 			created_at,
// 			updated_at
// 		FROM
// 			record
// 		WHERE
// 			id = $1
// 	`
// 	row := tx.QueryRow(
// 		ctx,
// 		query,
// 		id,
// 	)
// 	var item Record
// 	if err := row.Scan(
// 		&item.Id,
// 		&item.Note,
// 		&item.Amount,
// 		&item.IsExpense,
// 		&item.Book.Id,
// 		&item.Category.Id,
// 		&item.CreatedAt,
// 		&item.UpdatedAt,
// 	); err != nil {
// 		if err == pgx.ErrNoRows {
// 			log.Debug().Err(err).Msg("can't find any item")
// 			return Record{}, ErrNotFound
// 		}
// 		return Record{}, err
// 	}
// 	return item, nil
// }
