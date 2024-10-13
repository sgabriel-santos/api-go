package repositories

import (
	"api/src/models"
	"database/sql"
)

// Representa um repositório de pagamentoLog
type PagamentoLog struct {
	db *sql.DB
}

// Cria um repositório de usuários
func NewPagamentoLogRepository(db *sql.DB) *PagamentoLog {
	return &PagamentoLog{db}
}

// Armazena os logs de pagamento no banco de dados
func (repository PagamentoLog) CreatePagamentoLog(log models.PagamentoLog) error {
	stmt, err := repository.db.Prepare(`
        INSERT INTO pagamentos_log (id_usuario, method, data, data_id, code, data_cadastro)
        VALUES (?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(log.IDUsuario, log.Method, log.Data, log.DataID, log.Code, log.DataCadastro)
	return err
}
