package utils

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

func ParseNullStringToNullTime(s null.String) (t null.Time) {
	if !s.Valid {
		return
	}

	ts, err := time.Parse(time.RFC3339, s.String)

	if err != nil {
		return
	}

	return null.TimeFrom(ts)
}
