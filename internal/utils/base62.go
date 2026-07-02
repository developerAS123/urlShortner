package utils

import (
	"math/rand"
	"time"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateSlug generates a random base62 string of a given length
func GenerateSlug(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = base62Chars[rand.Intn(len(base62Chars))]
	}
	return string(bytes)
}
