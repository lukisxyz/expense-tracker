package account

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

func findAccountById(
	ctx context.Context,
	id ulid.ULID,
) (
	Account,
	error,
) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return Account{}, err
	}

	item, err := findItemById(ctx, tx, id)
	if err != nil {
		return Account{}, err
	}
	item.Password = ""

	err = tx.Commit(ctx)
	if err != nil {
		return Account{}, err
	}
	return item, err
}

func saveAccount(
	ctx context.Context,
	email, password string,
) (
	id ulid.ULID,
	err error,
) {
	Account, err := NewAccount(email, password)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	err = saveItem(ctx, tx, &Account)
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

	return Account.Id, nil
}

func updateAccount(
	ctx context.Context,
	id ulid.ULID,
	email, password string,
) (
	res Account,
	err error,
) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	Account, err := findItemById(ctx, tx, id)
	if err != nil {
		errRB := tx.Rollback(ctx)
		if errRB != nil {
			return
		}
		return
	}

	Account.Email = email
	Account.Password = password
	Account.UpdatedAt = null.TimeFrom(time.Now())

	err = saveItem(ctx, tx, &Account)
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

	return Account, nil
}
