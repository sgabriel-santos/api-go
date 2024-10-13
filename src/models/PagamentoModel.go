package models

import (
	"time"
)

// Representa um usuário no sistema
type PagamentoModel struct {
	ID               int       `json:"id"`
	IDUsuario        int       `json:"id_usuario"`
	Tipo             int       `json:"tipo"`
	NomeBeneficiario *string   `json:"nome_beneficiario"`
	Banco            *string   `json:"banco"`
	Agencia          *string   `json:"agencia"`
	CCDig            *string   `json:"cc_dig"`
	TipoConta        *string   `json:"tipo_conta"`
	TipoPessoa       *string   `json:"tipo_pessoa"`
	CPFCNPJ          *string   `json:"cpf_cnpj"`
	Valor            *string   `json:"valor"`
	Data             *string   `json:"data"`
	DataCode         *string   `json:"data_code"`
	Reference        *string   `json:"reference"`
	AddressKey       *string   `json:"addresskey"`
	Description      *string   `json:"description"`
	Response         *string   `json:"response"`
	Status           *string   `json:"status"`
	DataCadastro     time.Time `json:"data_cadastro"`
}
