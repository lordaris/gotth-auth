package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"
)

const (
	TokenLength = 32 // Length in bytes
)

func GenerateToken(expiry time.Duration) (string, []byte, time.Time, error) {
	tokenBytes := make([]byte, TokenLength)

	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", nil, time.Time{}, err
	}

	plaintextToken := base64.URLEncoding.EncodeToString(tokenBytes)

	// Hash the token
	hash := sha256.Sum256([]byte(plaintextToken))
	tokenHash := hash[:]

	expiryTime := time.Now().Add(expiry)

	return plaintextToken, tokenHash, expiryTime, nil
}

func ValidateTokenPlaintext(plaintextToken string) error {
	_, err := base64.URLEncoding.DecodeString(plaintextToken)
	if err != nil {
		return errors.New("invalid token format")
	}
	return nil
}

func HashToken(plaintextToken string) []byte {
	hash := sha256.Sum256([]byte(plaintextToken))
	return hash[:]
}
