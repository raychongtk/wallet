package util

import (
	"math/rand"
	"time"
)

var (
	candidates = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func AlphaNumeric(length int) string {
	rand.Seed(time.Now().UnixNano())
	runes := []rune(candidates)
	var result string
	for i := 0; i < length; i++ {
		result += string(runes[rand.Intn(len(candidates))])
	}
	return result
}

func Number(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}
