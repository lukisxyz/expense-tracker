package category

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Category struct {
	Id        ulid.ULID
	Name      string
	CreatedAt time.Time
	UpdatedAt null.Time
	OwnerId   ulid.ULID
}

func NewCategory(name string, ownerId ulid.ULID) (Category, error) {
	// TODO: check if user is valid
	cat := Category{
		Id:        ulid.Make(),
		Name:      name,
		OwnerId:   ownerId,
		CreatedAt: time.Now(),
	}
	return cat, nil
}
