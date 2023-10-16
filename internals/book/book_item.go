package book

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Book struct {
	Id        ulid.ULID
	Name      string
	CreatedAt time.Time
	UpdatedAt null.Time
	OwnerId   ulid.ULID
	IsDefault bool
}

func NewBook(name string, ownerId ulid.ULID) (Book, error) {
	// TODO: check if user is valid
	book := Book{
		Id:        ulid.Make(),
		Name:      name,
		OwnerId:   ownerId,
		IsDefault: false,
		CreatedAt: time.Now(),
	}
	return book, nil
}
