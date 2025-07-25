package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	var hashedPassword string

	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", ErrorHandler(errors.New("failed to generate salt"), "Unable to add data")
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)
	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)

	hashedPassword = encodedHash

	return hashedPassword, nil
}

func VerifyPassword(password, hashPassword string) error {

	parts := strings.Split(hashPassword, ".")
	if len(parts) != 2 {
		return ErrorHandler(errors.New("invalid encoded hash format"), "invalid encoded hash format")
	}

	saltBase64 := parts[0]
	hashedPasswordBase64 := parts[1]

	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return ErrorHandler(err, "failed to decode the salt")
	}

	hashedPassword, err := base64.StdEncoding.DecodeString(hashedPasswordBase64)
	if err != nil {
		return ErrorHandler(err, "failed to decode the hashed password")
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	if len(hash) != len(hashedPassword) {
		log.Println(password, hashPassword)
		return ErrorHandler(errors.New("invalid password"), "invalid password")
	}

	if subtle.ConstantTimeCompare(hash, hashedPassword) != 1 {
		return ErrorHandler(errors.New("invalid password"), "invalid password")
	}

	return nil
}
