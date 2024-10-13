package models

import (
	"api/src/security"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

// Representa um usuário no sistema
type UserModel struct {
	ID                uint64     `json:"id,omitempty"`
	UID               string     `json:"uid,omitempty"`
	Tipo              int        `json:"tipo,omitempty"`
	Email             string     `json:"email,omitempty"`
	EmailVerificado   int64      `json:"email_verificado,omitempty"`
	Senha             string     `json:"senha,omitempty"`
	Nome              string     `json:"nome,omitempty"`
	CpfCnpj           string     `json:"cpf_cnpj,omitempty"`
	Cnpj              string     `json:"cnpj,omitempty"`
	Cpf               string     `json:"cpf,omitempty"`
	Rg                string     `json:"rg,omitempty"`
	OrgaoEmissor      string     `json:"orgao_emissor,omitempty"`
	Telefone          string     `json:"telefone,omitempty"`
	Celular           string     `json:"celular,omitempty"`
	Pin               string     `json:"pin,omitempty"`
	Timezone          *string    `json:"timezone,omitempty"`
	Blockchain        string     `json:"blockchain,omitempty"`
	TaxaDeposito      string     `json:"taxa_deposito,omitempty"`
	TaxaSaque         string     `json:"taxa_saque,omitempty"`
	TaxaConversaoPorc string     `json:"taxa_conversao_porc,omitempty"`
	TaxaConversaoCent string     `json:"taxa_conversao_cent,omitempty"`
	Lang              string     `json:"lang,omitempty"`
	API               int        `json:"api,omitempty"`
	Status            int        `json:"status,omitempty"`
	DataCadastro      *time.Time `json:"data_cadastro,omitempty"`
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
	if user.Nome == "" {
		return errors.New("o campo nome é obrigatório e não pode estar em branco")
	}

	if user.Email == "" {
		return errors.New("o campo email é obrigatório e não pode estar em branco")
	}

	if erro := checkmail.ValidateFormat(user.Email); erro != nil {
		return errors.New("o e-mail inserido é inválido")
	}

	if etapa == "cadastro" && user.Senha == "" {
		return errors.New("a campo senha é obrigatório e não pode estar em branco")
	}

	return nil
}

func (user *UserModel) format(etapa string) error {
	user.Nome = strings.TrimSpace(user.Nome)
	user.Email = strings.TrimSpace(user.Email)

	if etapa == "cadastro" {
		// passwordHash, erro := security.Hash(user.Senha)
		// if erro != nil {
		// 	return erro
		// }

		passwordHash := security.Hash((user.Senha))

		// pinHash, erro := security.Hash(user.Pin)
		// if erro != nil {
		// 	return erro
		// }

		pinHash := security.Hash(user.Pin)

		user.Senha = string(passwordHash)
		user.Pin = string(pinHash)
		user.UID = user.generateUID()
		user.Status = 2
		user.Tipo = 1
	}

	return nil
}

func (user *UserModel) generateUID() string {
	// Obtém o timestamp em nanosegundos como uma string
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

	// Calcula o hash MD5 do timestamp
	hash := md5.New()
	hash.Write([]byte(timestamp))

	// Retorna o hash em formato hexadecimal
	return hex.EncodeToString(hash.Sum(nil))
}
