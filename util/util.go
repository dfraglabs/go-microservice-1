package util

import (
	"math/rand"
	"time"
)

const (
	RAND_LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func RandString(n int, t time.Time) string {
	src := rand.NewSource(t.UTC().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = RAND_LETTERS[src.Int63()%int64(len(RAND_LETTERS))]
	}
	return string(b)
}
