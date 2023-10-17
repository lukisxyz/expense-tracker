package category

import (
	"encoding/json"
	"time"

	"github.com/flukis/expt/service/utils"
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

func (b *Category) MarshalJSON() ([]byte, error) {
	var j struct {
		Id        ulid.ULID `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	}
	j.Id = b.Id
	j.Name = b.Name
	j.CreatedAt = b.CreatedAt
	return json.Marshal(j)
}

func (b *Category) UnmarshalJSON(data []byte) error {
	var j struct {
		Id        ulid.ULID   `json:"id"`
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

	b = &Category{ //nolint:all // not implement yet
		Id:        j.Id,
		Name:      j.Name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		OwnerId:   j.OwnerId,
	}
	return nil
}
