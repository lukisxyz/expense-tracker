package record

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type RecordCategory struct {
	Id   ulid.ULID
	Name string
}

type RecordBook struct {
	Id   ulid.ULID
	Name string
}

type Record struct {
	Id        ulid.ULID
	Note      string
	Amount    float64
	IsExpense bool
	Category  RecordCategory
	Book      RecordBook
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
		Category: RecordCategory{
			Id: categoryId,
		},
		Book: RecordBook{
			Id: bookId,
		},
	}
	return cat, nil
}
