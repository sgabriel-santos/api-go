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
func (repository PagamentoLog) CreatePagamentoLog(log models.PagamentoLogModel) error {
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

// GetPagamentoLog busca um log de pagamento pelo userID, timestamp (code) e reference
func (repository PagamentoLog) GetPagamentoLog(userID uint64, timestamp string, reference string) (*models.PagamentoLogModel, error) {
	var pagamentoLog models.PagamentoLogModel

	query := `
        SELECT id_usuario, method, data, data_id, code, data_cadastro
        FROM pagamentos_log
        WHERE id_usuario = ? AND code = ? AND data_id = ?`

	// Executar a consulta no banco de dados
	row := repository.db.QueryRow(query, userID, timestamp, reference)

	// Fazer o scan dos resultados
	err := row.Scan(
		&pagamentoLog.IDUsuario,
		&pagamentoLog.Method,
		&pagamentoLog.Data,
		&pagamentoLog.DataID,
		&pagamentoLog.Code,
		&pagamentoLog.DataCadastro,
	)

	if err != nil {
		// Se nenhum registro for encontrado, retornar nil (sem erro)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		// Retornar o erro se ocorrer algum problema durante a consulta
		return nil, err
	}

	// Retornar o log de pagamento encontrado
	return &pagamentoLog, nil
}

// GetSaldoByUsuario busca o saldo de um usuário pelo ID
func (repository PagamentoLog) GetSaldoByUsuario(userID uint64) (float64, error) {
	var saldo float64

	query := `
        SELECT valor
        FROM balancas
        WHERE id_usuario = ? AND moeda = 'BRL'`

	// Executar a consulta no banco de dados
	err := repository.db.QueryRow(query, userID).Scan(&saldo)
	if err != nil {
		if err == sql.ErrNoRows {
			// Se não houver saldo registrado para o usuário, retorna saldo 0
			return 0, nil
		}
		// Se ocorrer algum outro erro, retorna o erro
		return 0, err
	}

	// Retornar o saldo encontrado
	return saldo, nil
}
