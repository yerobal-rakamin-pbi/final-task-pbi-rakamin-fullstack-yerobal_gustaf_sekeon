package password

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string, saltRound int64) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), int(saltRound))
	return string(bytes), err
}

func Compare(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}