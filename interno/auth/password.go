package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	"crypto/pbkdf2"
)

const (
	passwordAlgorithm = "pbkdf2_sha256"
	passwordIter      = 100000
	passwordSaltBytes = 16
	passwordKeyBytes  = 32
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, passwordSaltBytes)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return hashPasswordWithSalt(password, salt, passwordIter)
}

func VerifyPassword(password, encoded string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 4 || parts[0] != passwordAlgorithm {
		return false
	}

	iter, err := strconv.Atoi(parts[1])
	if err != nil || iter < 10000 {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false
	}
	actual, err := pbkdf2.Key(sha256.New, password, salt, iter, len(expected))
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(actual, expected) == 1
}

func MustHashPassword(password string) string {
	hash, err := HashPassword(password)
	if err != nil {
		panic(err)
	}
	return hash
}

func hashPasswordWithSalt(password string, salt []byte, iter int) (string, error) {
	if strings.TrimSpace(password) == "" {
		return "", errors.New("senha nao pode ser vazia")
	}
	key, err := pbkdf2.Key(sha256.New, password, salt, iter, passwordKeyBytes)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{
		passwordAlgorithm,
		strconv.Itoa(iter),
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	}, "$"), nil
}
