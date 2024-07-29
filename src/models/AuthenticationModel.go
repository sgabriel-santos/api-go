package models

// Contém o token e o id do usuário autenticado
type AuthenticationModel struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}
