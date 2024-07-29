package models

import (
	"api/src/security"
	"errors"
	"strings"

	"github.com/badoux/checkmail"
)

// Representa um usuário no sistema
type UserModel struct {
	ID       uint64 `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// Chama os métodos para validar e formatar o usuário recebido
func (user *UserModel) Prepare(etapa string) error {
	if erro := user.validate(etapa); erro != nil {
		return erro
	}

	if erro := user.format(etapa); erro != nil {
		return erro
	}

	return nil
}

func (user *UserModel) validate(etapa string) error {
	if user.Name == "" {
		return errors.New("o campo name é obrigatório e não pode estar em branco")
	}

	if user.Email == "" {
		return errors.New("o campo email é obrigatório e não pode estar em branco")
	}

	if erro := checkmail.ValidateFormat(user.Email); erro != nil {
		return errors.New("o e-mail inserido é inválido")
	}

	if etapa == "cadastro" && user.Password == "" {
		return errors.New("a campo password é obrigatório e não pode estar em branco")
	}

	return nil
}

func (user *UserModel) format(etapa string) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	if etapa == "cadastro" {
		senhaComHash, erro := security.Hash(user.Password)
		if erro != nil {
			return erro
		}

		user.Password = string(senhaComHash)
	}

	return nil
}
