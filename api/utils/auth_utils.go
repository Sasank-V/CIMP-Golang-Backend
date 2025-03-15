package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashSHA256(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
