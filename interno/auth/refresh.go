package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

const refreshTokenBytes = 48

func NewRefreshToken() (raw string, hash string, err error) {
	data := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(data); err != nil {
		return "", "", err
	}
	raw = base64.RawURLEncoding.EncodeToString(data)
	return raw, HashRefreshToken(raw), nil
}

func HashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
