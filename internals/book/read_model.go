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
	Books []Book `json:"books"`
	Count int    `json:"count"`
}

var emptyList BookList

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
	err := row.Scan(&itemCount)
	if err != nil {
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
				is_default,
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
		var isDefault bool
		var name string
		var ownerId ulid.ULID
		var createdAt time.Time
		var updatedAt null.Time
		if !rows.Next() {
			break
		}
		if err := rows.Scan(
			&id,
			&isDefault,
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
			IsDefault: isDefault,
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
