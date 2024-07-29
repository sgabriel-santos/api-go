package security

import (
	"golang.org/x/crypto/bcrypt"
)

// Recebe uma string e coloca um hash nela
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// Compara uma senha e um hash e retorna se elas são iguais
func VerifyPassword(hashPassword, passwordStr string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(passwordStr))
}
