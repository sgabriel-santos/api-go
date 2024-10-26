package models

import "time"

// ConversaoModel representa uma conversão no sistema
type ConversaoModel struct {
	ID              uint64    `json:"id"`
	IDUsuario       uint64    `json:"id_usuario"`
	Simbolo         string    `json:"simbolo"`
	De              string    `json:"de"`
	Para            string    `json:"para"`
	Valor           float64   `json:"valor"`
	ValorFinal      float64   `json:"valor_final"`
	Preco           float64   `json:"preco,omitempty"`
	Bid             float64   `json:"bid"`
	Ask             float64   `json:"ask"`
	Last            float64   `json:"last"`
	Liquidar        int       `json:"liquidar,omitempty"` // 0: pendente, 1: em andamento, 2: concluído
	LiquidarRetorno *string   `json:"liquidar_retorno,omitempty"`
	Status          int       `json:"status"` // 0: pendente, 1: em andamento, 2: concluído, 3: cancelado
	StatusDescricao string    `json:"status_descricao,omitempty"`
	DataCadastro    time.Time `json:"data_cadastro"`
}
