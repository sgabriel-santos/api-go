package security

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
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

// Função para converter a string em hexadecimal (equivalente ao bin2hex do PHP)
func StringParaHex(s string) string {
	return hex.EncodeToString([]byte(s))
}

// GerarHashComSalt gera um hash da senha usando o salt convertido para hexadecimal
func GerarHashComSalt(senha string) string {
	salt := StringParaHex("M4N4U54M") // Converte "M4N4U54M" para hexadecimal, que será usado como salt
	senhaComSalt := salt + senha

	// Gerar o hash SHA-512
	hash := sha512.New()
	hash.Write([]byte(senhaComSalt))

	// Converter o hash para string hexadecimal
	return hex.EncodeToString(hash.Sum(nil))
}

// // Compara uma senha e um hash e retorna se elas são iguais
// func VerifyPassword(hashPassword, passwordStr string) error {
// 	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(passwordStr))
// }

// Compara uma senha e um hash e retorna se elas são iguais
func VerifyPassword(hashPassword, passwordStr string) bool {
	return hashPassword == Hash(passwordStr)
}
