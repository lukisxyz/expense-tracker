package account

import (
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

type Account struct {
	Id            ulid.ULID
	Email         string
	Password      string
	DefaultBookId ulid.ULID
	CreatedAt     time.Time
	UpdatedAt     null.Time
}

func NewAccount(email, password string) (Account, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Account{}, err
	}
	acc := Account{
		Id:        ulid.Make(),
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}
	return acc, nil
}
