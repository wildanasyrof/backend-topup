package hash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func ComparePassword(hashed string, plain string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)); err != nil {
		return errors.New("invalid credentials")
	}
	return nil
}
