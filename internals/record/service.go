package record

import (
	"context"

	"github.com/oklog/ulid/v2"
)

func listRecords(
	ctx context.Context,
	bookId ulid.ULID,
) (
	RecordList,
	error,
) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return emptyList, err
	}

	list, err := findItemsByBookId(ctx, tx, bookId)
	if err != nil {
		return emptyList, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return emptyList, err
	}
	return list, err
}

func saveRecord(
	ctx context.Context,
	bookId ulid.ULID,
	categoryId ulid.ULID,
	note string,
	amount float64,
	isExpense bool,
) (
	id ulid.ULID,
	err error,
) {
	ctg, err := NewRecord(note, amount, isExpense, categoryId, bookId)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}

	err = saveItem(ctx, tx, &ctg)
	if err != nil {
		errRB := tx.Rollback(ctx)
		if errRB != nil {
			return
		}
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		return
	}

	return ctg.Id, nil
}

// func findRecordById(
// 	ctx context.Context,
// 	id ulid.ULID,
// ) (
// 	Record,
// 	error,
// ) {
// 	tx, err := pool.Begin(ctx)
// 	if err != nil {
// 		return Record{}, err
// 	}

// 	item, err := findItemById(ctx, tx, id)
// 	if err != nil {
// 		return Record{}, err
// 	}

// 	err = tx.Commit(ctx)
// 	if err != nil {
// 		return Record{}, err
// 	}
// 	return item, err
// }
