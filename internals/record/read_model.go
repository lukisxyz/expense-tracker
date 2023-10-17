package record

import (
	"context"
	"time"

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
				r.book_id = $1
			ORDER BY r.id;
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
		var category RecordCategory
		var book RecordBook
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
			&createdAt,
			&updatedAt,
			&category.Id,
			&category.Name,
			&book.Id,
			&book.Name,
		); err != nil {
			log.Warn().Err(err).Msg("cannot scan an item")
			return emptyList, err
		}
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

func findItemsByCategoryId(
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
		`SELECT COUNT(id) as c FROM record WHERE category_id = $1`,
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
				r.category_id = $1
			ORDER BY r.id;
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
		var category RecordCategory
		var book RecordBook
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
			&createdAt,
			&updatedAt,
			&category.Id,
			&category.Name,
			&book.Id,
			&book.Name,
		); err != nil {
			log.Warn().Err(err).Msg("cannot scan an item")
			return emptyList, err
		}
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
