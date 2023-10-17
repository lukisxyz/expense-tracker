package category

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

func listCategories(
	ctx context.Context,
	ownerId ulid.ULID,
) (
	CategoryList,
	error,
) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return emptyList, err
	}

	list, err := findItemsByOwnerId(ctx, tx, ownerId)
	if err != nil {
		return emptyList, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return emptyList, err
	}
	return list, err
}

func saveCategory(
	ctx context.Context,
	ownerId ulid.ULID,
	name string,
) (
	id ulid.ULID,
	err error,
) {
	ctg, err := NewCategory(name, ownerId)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	err = saveItem(ctx, tx, &ctg)
	if err != nil {
		errRB := tx.Rollback(ctx)
		if errRB != nil {
			return
		}
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		return
	}

	return ctg.Id, nil
}

func findCategoryById(
	ctx context.Context,
	id ulid.ULID,
) (
	Category,
	error,
) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return Category{}, err
	}

	item, err := findItemById(ctx, tx, id)
	if err != nil {
		return Category{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return Category{}, err
	}
	return item, err
}

func updateBook(
	ctx context.Context,
	id ulid.ULID,
	ownerId ulid.ULID,
	name string,
) (
	res Category,
	err error,
) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	ctg, err := findItemById(ctx, tx, id)
	if err != nil {
		errRB := tx.Rollback(ctx)
		if errRB != nil {
			return
		}
		return
	}

	ctg.OwnerId = ownerId
	ctg.Name = name
	ctg.UpdatedAt = null.TimeFrom(time.Now())

	err = saveItem(ctx, tx, &ctg)
	if err != nil {
		errRB := tx.Rollback(ctx)
		if errRB != nil {
			return
		}
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		return
	}

	return ctg, nil
}
