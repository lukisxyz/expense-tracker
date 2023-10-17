package account

import (
	"encoding/json"
	"time"

	"github.com/flukis/expt/service/utils"
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

func (b *Account) MarshalJSON() ([]byte, error) {
	var j struct {
		Id            ulid.ULID `json:"id"`
		Email         string    `json:"email"`
		DefaultBookId ulid.ULID `json:"default"`
		CreatedAt     time.Time `json:"created_at"`
	}
	j.Id = b.Id
	j.Email = b.Email
	j.DefaultBookId = b.DefaultBookId
	j.CreatedAt = b.CreatedAt
	return json.Marshal(j)
}

func (b *Account) UnmarshalJSON(data []byte) error {
	var j struct {
		Id            ulid.ULID   `json:"id"`
		Email         string      `json:"email"`
		Password      string      `json:"password"`
		DefaultBookId ulid.ULID   `json:"default"`
		CreatedAt     string      `json:"created_at"`
		UpdatedAt     null.String `json:"updated_at"`
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

	b = &Account{ //nolint:all // not implement yet
		Id:            j.Id,
		Email:         j.Email,
		Password:      j.Password,
		DefaultBookId: j.DefaultBookId,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
	return nil
}
