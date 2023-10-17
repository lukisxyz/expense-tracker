package account

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Account struct {
	Id        ulid.ULID
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt null.Time
}

func NewAccount(email, password string) (Account, error) {
	// TODO: check if user is valid
	acc := Account{
		Id:        ulid.Make(),
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}
	return acc, nil
}
