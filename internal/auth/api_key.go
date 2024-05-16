package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func generateAPIKey() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(b)

	return hex.EncodeToString(hash[:]), nil
}
