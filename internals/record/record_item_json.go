package record

import (
	"encoding/json"
	"time"

	"github.com/flukis/expt/service/utils"
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

func (b *Record) MarshalJSON() ([]byte, error) {
	var j struct {
		Id        ulid.ULID      `json:"id"`
		Note      string         `json:"note"`
		Amount    float64        `json:"amount"`
		Category  RecordCategory `json:"category"`
		Book      RecordBook     `json:"book"`
		CreatedAt time.Time      `json:"created_at"`
	}
	j.Id = b.Id
	j.Amount = b.Amount
	j.Note = b.Note
	j.CreatedAt = b.CreatedAt
	j.Category.Id = b.Category.Id
	j.Category.Name = b.Category.Name
	j.Book.Id = b.Book.Id
	j.Book.Name = b.Book.Name
	return json.Marshal(j)
}

func (b *Record) UnmarshalJSON(data []byte) error {
	var j struct {
		Id        ulid.ULID      `json:"id"`
		Note      string         `json:"note"`
		Amount    float64        `json:"amount"`
		CreatedAt string         `json:"created_at"`
		Category  RecordCategory `json:"category"`
		Book      RecordBook     `json:"book"`
		UpdatedAt null.String    `json:"updated_at"`
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

	b = &Record{ //nolint:all // not implement yet
		Id:        j.Id,
		Note:      j.Note,
		Amount:    j.Amount,
		Category:  j.Category,
		Book:      j.Book,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	return nil
}
