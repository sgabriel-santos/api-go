package models

import "time"

// PagamentoLog representa um log de pagamento decodificado
type PagamentoLog struct {
	IDUsuario    uint64    `json:"id_usuario"`
	Method       string    `json:"method"`
	Data         string    `json:"data"`
	DataID       string    `json:"data_id"`
	Code         string    `json:"code"`
	DataCadastro time.Time `json:"data_cadastro"`
}
