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
