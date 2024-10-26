package models

import "time"

// DepositoModel representa um depósito em USDT
type DepositoModel struct {
	ID              uint64     `json:"id"`
	IDUsuario       uint64     `json:"id_usuario"`
	Referencia      *string    `json:"referencia"`
	Moeda           string     `json:"moeda"`
	Endereco        *string    `json:"endereco"`
	Valor           float64    `json:"valor"`
	Transacao       string     `json:"transacao"`
	Retorno         string     `json:"retorno,omitempty"`
	PixImagem       *string    `json:"pix_imagem,omitempty"`
	Status          int        `json:"status"` // 0: pendente; 1: em andamento; 2: concluído; 3: cancelado;
	StatusDescricao string     `json:"status_descricao,omitempty"`
	DataConcluido   *time.Time `json:"data_concluido,omitempty"`
	DataCadastro    time.Time  `json:"data_cadastro"`
}
