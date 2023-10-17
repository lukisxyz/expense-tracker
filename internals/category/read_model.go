package category

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
)

type CategoryList struct {
	Categories []Category `json:"data"`
	Count      int        `json:"count"`
}

var emptyList = CategoryList{
	Categories: []Category{},
	Count:      0,
}

func findItemsByOwnerId(
	ctx context.Context,
	tx pgx.Tx,
	id ulid.ULID,
) (
	CategoryList,
	error,
) {
	var itemCount int

	row := tx.QueryRow(
		ctx,
		`SELECT COUNT(id) as c FROM category WHERE owner_id = $1`,
		id,
	)

	if err := row.Scan(&itemCount); err != nil {
		log.Warn().Err(err).Msg("cannot find a count in category")
		return emptyList, err
	}

	if itemCount == 0 {
		return emptyList, nil
	}

	log.Debug().Int("count", itemCount).Msg("found category items")
	items := make([]Category, itemCount)
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
				category
			WHERE
				owner_id = $1
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
		items[i] = Category{
			Id:        id,
			Name:      name,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			OwnerId:   ownerId,
		}
	}

	list := CategoryList{
		Categories: items,
		Count:      itemCount,
	}

	return list, nil
}
