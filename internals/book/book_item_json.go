package book

import (
	"encoding/json"
	"time"

	"github.com/flukis/expt/service/utils"
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

func (b *Book) MarshalJSON() ([]byte, error) {
	var j struct {
		Id        ulid.ULID `json:"id"`
		Name      string    `json:"name"`
		IsDefault bool      `json:"is_default"`
		CreatedAt time.Time `json:"created_at"`
	}
	j.Id = b.Id
	j.Name = b.Name
	j.IsDefault = b.IsDefault
	j.CreatedAt = b.CreatedAt
	return json.Marshal(j)
}

func (b *Book) UnmarshalJSON(data []byte) error {
	var j struct {
		Id        ulid.ULID   `json:"id"`
		IsDefault bool        `json:"is_default"`
		Name      string      `json:"name"`
		CreatedAt string      `json:"created_at"`
		UpdatedAt null.String `json:"updated_at"`
		OwnerId   ulid.ULID   `json:"owner_id"`
	}
	err := json.Unmarshal(data, &j)
	if err != nil {
		return err
	}
	createdAt, err := time.Parse(time.RFC3339, j.CreatedAt)
	if err != nil {
		return err
	}
	updatedAt := utils.ParseNullStringToNullTime(j.UpdatedAt)

	b = &Book{ //nolint:all // not implement yet
		Id:        j.Id,
		IsDefault: j.IsDefault,
		Name:      j.Name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		OwnerId:   j.OwnerId,
	}
	return nil
}