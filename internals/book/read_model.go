package book

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
)

type BookList struct {
	Books []Book `json:"data"`
	Count int    `json:"count"`
}

var emptyList = BookList{
	Books: []Book{},
	Count: 0,
}

func findItemsByOwnerId(
	ctx context.Context,
	tx pgx.Tx,
	id ulid.ULID,
) (
	BookList,
	error,
) {
	var itemCount int

	row := tx.QueryRow(
		ctx,
		`SELECT COUNT(id) as c FROM book WHERE owner_id = $1`,
		id,
	)

	if err := row.Scan(&itemCount); err != nil {
		log.Warn().Err(err).Msg("cannot find a count in book")
		return emptyList, err
	}

	if itemCount == 0 {
		return emptyList, nil
	}

	log.Debug().Int("count", itemCount).Msg("found book items")
	items := make([]Book, itemCount)
	rows, err := tx.Query(
		ctx,
		`
			SELECT
				id,
				name,
				owner_id,
				created_at,
				updated_at
			FROM
				book
			WHERE
				owner_id = $1
			ORDER BY id
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
		var name string
		var ownerId ulid.ULID
		var createdAt time.Time
		var updatedAt null.Time
		if !rows.Next() {
			break
		}
		if err := rows.Scan(
			&id,
			&name,
			&ownerId,
			&createdAt,
			&updatedAt,
		); err != nil {
			log.Warn().Err(err).Msg("cannot scan an item")
			return emptyList, err
		}
		items[i] = Book{
			Id:        id,
			Name:      name,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			OwnerId:   ownerId,
		}
	}

	list := BookList{
		Books: items,
		Count: itemCount,
	}

	return list, nil
}
