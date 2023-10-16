package book

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

func listBooks(
	ctx context.Context,
	ownerId ulid.ULID,
) (
	BookList,
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

func findBook(
	ctx context.Context,
	id ulid.ULID,
) (
	Book,
	error,
) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return Book{}, err
	}

	item, err := findItemById(ctx, tx, id)
	if err != nil {
		return Book{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return Book{}, err
	}
	return item, err
}

func saveBook(
	ctx context.Context,
	ownerId ulid.ULID,
	name string,
) (
	id ulid.ULID,
	err error,
) {
	book, err := NewBook(name, ownerId)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	err = saveItem(ctx, tx, &book)
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

	return book.Id, nil
}

func updateBook(
	ctx context.Context,
	id ulid.ULID,
	ownerId ulid.ULID,
	name string,
) (
	res Book,
	err error,
) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	book, err := findItemById(ctx, tx, id)
	if err != nil {
		errRB := tx.Rollback(ctx)
		if errRB != nil {
			return
		}
		return
	}

	book.OwnerId = ownerId
	book.Name = name
	book.UpdatedAt = null.TimeFrom(time.Now())

	err = saveItem(ctx, tx, &book)
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

	return book, nil
}

func makeDefault(
	ctx context.Context,
	ownerId ulid.ULID,
	id ulid.ULID,
) (err error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	book, err := findItemById(ctx, tx, id)
	if err != nil {
		errRB := tx.Rollback(ctx)
		if errRB != nil {
			return
		}
		return
	}

	existingDefaultBook, err := findDefaultItemById(ctx, tx, ownerId)
	if err != nil && err != ErrNotFound {
		errRB := tx.Rollback(ctx)
		if errRB != nil {
			return
		}
		return
	}

	if err != nil && err == ErrNotFound {
		existingDefaultBook.IsDefault = false
		existingDefaultBook.UpdatedAt = null.TimeFrom(time.Now())
		err = saveItem(ctx, tx, &existingDefaultBook)
		if err != nil {
			errRB := tx.Rollback(ctx)
			if errRB != nil {
				return
			}
			return
		}
	}

	book.IsDefault = true
	book.UpdatedAt = null.TimeFrom(time.Now())
	err = saveItem(ctx, tx, &book)
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

	return nil
}
