package record

import (
	"context"
	"time"

	bookModel "github.com/flukis/expt/service/internals/book"
	categoryModel "github.com/flukis/expt/service/internals/category"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
)

type RecordList struct {
	Records []Record `json:"data"`
	Count   int      `json:"count"`
}

var emptyList = RecordList{
	Records: []Record{},
	Count:   0,
}

func findItemsByBookId(
	ctx context.Context,
	tx pgx.Tx,
	id ulid.ULID,
) (
	RecordList,
	error,
) {
	var itemCount int

	row := tx.QueryRow(
		ctx,
		`SELECT COUNT(id) as c FROM record WHERE book_id = $1`,
		id,
	)

	if err := row.Scan(&itemCount); err != nil {
		log.Warn().Err(err).Msg("cannot find a count in category")
		return emptyList, err
	}

	if itemCount == 0 {
		return emptyList, nil
	}

	log.Debug().Int("count", itemCount).Msg("found record items")
	items := make([]Record, itemCount)
	rows, err := tx.Query(
		ctx,
		`
			SELECT
				id,
				note,
				amount,
				is_expense,
				book_id,
				category_id,
				created_at,
				updated_at
			FROM
				record
			WHERE
				book_id = $1
			ORDER by id
		`,
		id,
	)
	if err != nil {
		return emptyList, err
	}
	defer rows.Close()

	var i int
	for i = range items {
		var id ulid.ULID
		var note string
		var amount float64
		var isExpense bool
		var categoryId ulid.ULID
		var category categoryModel.Category
		var bookId ulid.ULID
		var book bookModel.Book
		var createdAt time.Time
		var updatedAt null.Time
		if !rows.Next() {
			break
		}
		if err := rows.Scan(
			&id,
			&note,
			&amount,
			&isExpense,
			&categoryId,
			&bookId,
			&createdAt,
			&updatedAt,
		); err != nil {
			log.Warn().Err(err).Msg("cannot scan an item")
			return emptyList, err
		}
		category.Id = categoryId
		book.Id = bookId
		items[i] = Record{
			Id:        id,
			Amount:    amount,
			IsExpense: isExpense,
			Category:  category,
			Book:      book,
			Note:      note,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}

	list := RecordList{
		Records: items,
		Count:   itemCount,
	}

	return list, nil
}
