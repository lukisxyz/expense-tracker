package record

import (
	"encoding/json"
	"time"

	"github.com/flukis/expt/service/internals/book"
	"github.com/flukis/expt/service/internals/category"
	"github.com/flukis/expt/service/utils"
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

func (b *Record) MarshalJSON() ([]byte, error) {
	var j struct {
		Id           ulid.ULID `json:"id"`
		Note         string    `json:"note"`
		Amount       float64   `json:"amount"`
		CategoryId   ulid.ULID `json:"category_id"`
		CategoryName string    `json:"category_name"`
		BookId       ulid.ULID `json:"book_id"`
		BookName     string    `json:"book_name"`
		CreatedAt    time.Time `json:"created_at"`
	}
	j.Id = b.Id
	j.Amount = b.Amount
	j.Note = b.Note
	j.CreatedAt = b.CreatedAt
	j.CategoryId = b.Category.Id
	j.CategoryName = b.Category.Name
	j.BookId = b.Book.Id
	j.BookName = b.Book.Name
	return json.Marshal(j)
}

func (b *Record) UnmarshalJSON(data []byte) error {
	var j struct {
		Id        ulid.ULID         `json:"id"`
		Note      string            `json:"note"`
		Amount    float64           `json:"amount"`
		CreatedAt string            `json:"created_at"`
		Category  category.Category `json:"category"`
		Book      book.Book         `json:"book"`
		UpdatedAt null.String       `json:"updated_at"`
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
