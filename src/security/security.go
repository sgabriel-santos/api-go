package security

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Recebe uma string e coloca um hash nela
// func Hash(password string) ([]byte, error) {
// 	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// }

// Recebe uma string e coloca um hash nela
func Hash(input string) string {
	// Gera o hash SHA-256 da senha
	hash := sha256.Sum256([]byte(input))

	// Converte o hash para hexadecimal
	hashHex := hex.EncodeToString(hash[:])

	// Retorna apenas os primeiros 13 caracteres do hash
	return hashHex[:13]
}

// // Compara uma senha e um hash e retorna se elas são iguais
// func VerifyPassword(hashPassword, passwordStr string) error {
// 	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(passwordStr))
// }

// Compara uma senha e um hash e retorna se elas são iguais
func VerifyPassword(hashPassword, passwordStr string) bool {
	fmt.Println(passwordStr, hashPassword, Hash(passwordStr))
	return hashPassword == Hash(passwordStr)
}
