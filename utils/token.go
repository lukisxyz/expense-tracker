package utils

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func GenerateRefreshToken() (ulid.ULID, error) {
	entropy := rand.Reader
	ms := ulid.Timestamp(time.Now())
	return ulid.New(ms, entropy)
}
