package utils

import (
	"crypto/sha256"
	"fmt"
)

func HashSHA256(inp string) string {
	hash := sha256.Sum256([]byte(inp))
	return fmt.Sprintf("%x", hash)
}
