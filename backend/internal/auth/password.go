package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 2, 16)
	passwordHash := append(salt, hash...)

	return base64.StdEncoding.EncodeToString(passwordHash), nil
}

func VerifyPasswordHash(password, encodedHash string) (bool, error) {
	decodedHash, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false, err
	}

	salt := decodedHash[:16]
	correctHash := decodedHash[16:]

	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 2, 16)

	return subtle.ConstantTimeCompare(correctHash, hash) == 1, nil
}
