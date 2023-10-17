package account

import (
	"context"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

func loginAccount(
	ctx context.Context,
	email, password string,
) (
	Account,
	error,
) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return Account{}, err
	}

	item, err := findItemByEmail(ctx, tx, email)
	if err != nil {
		return Account{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(item.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return Account{}, ErrWrongPassword
		}
		return Account{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return Account{}, err
	}

	item.Password = ""
	return item, nil
}

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
	acc, err := NewAccount(email, password)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	err = saveItem(ctx, tx, &acc)
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

	return acc.Id, nil
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	Account.Email = email
	Account.Password = string(hashedPassword)
	Account.UpdatedAt = null.TimeFrom(time.Now())

	err = saveItem(ctx, tx, &Account)
	if err != nil {
		errRB := tx.Rollback(ctx)
		if errRB != nil {
			return
		}
		return
	}
	Account.Password = ""

	err = tx.Commit(ctx)
	if err != nil {
		return
	}

	return Account, nil
}
