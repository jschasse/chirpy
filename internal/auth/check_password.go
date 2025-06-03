package auth

import(
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}