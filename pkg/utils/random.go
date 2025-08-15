package utils

import (
	"math/rand"
	"time"
)

func GenerateSixDigitCode() int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	code := r.Intn(900000) + 100000

	return code
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
