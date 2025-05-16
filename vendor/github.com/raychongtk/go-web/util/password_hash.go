package util

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

func Hash(value string) string {
	iteration := Number(100, 500)
	salt := AlphaNumeric(10)
	hash := computeHash(value, iteration, salt)
	return strconv.Itoa(iteration) + ":" + salt + ":" + hash
}

func ValidateHash(hash string, value string) bool {
	parts := strings.Split(hash, ":")
	iteration, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}

	salt := parts[1]
	return computeHash(value, iteration, salt) == parts[2]
}

func computeHash(value string, iteration int, salt string) string {
	h := sha256.New()
	hash := value

	for i := 0; i < iteration; i++ {
		h.Write([]byte(hash + salt))
		hash = hex.EncodeToString(h.Sum(nil))
	}

	return hash
}
