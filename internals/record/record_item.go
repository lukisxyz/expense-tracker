package record

import (
	"time"

	"github.com/flukis/expt/service/internals/book"
	"github.com/flukis/expt/service/internals/category"
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Record struct {
	Id        ulid.ULID
	Note      string
	Amount    float64
	IsExpense bool
	Category  category.Category
	Book      book.Book
	CreatedAt time.Time
	UpdatedAt null.Time
}

func NewRecord(
	note string,
	amount float64,
	isExpense bool,
	categoryId ulid.ULID,
	bookId ulid.ULID,
) (
	Record,
	error,
) {
	cat := Record{
		Id:        ulid.Make(),
		Note:      note,
		IsExpense: isExpense,
		Amount:    amount,
		CreatedAt: time.Now(),
		Category: category.Category{
			Id: categoryId,
		},
		Book: book.Book{
			Id: bookId,
		},
	}
	return cat, nil
}
