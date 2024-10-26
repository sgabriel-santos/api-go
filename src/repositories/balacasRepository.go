package repositories

import (
	"database/sql"
	"errors"
)

// Representa um repositório de balanca
type Balanca struct {
	db *sql.DB
}

// Cria um repositório de usuários
func NewBalancaRepository(db *sql.DB) *Balanca {
	return &Balanca{db}
}

// GetSaldoByUsuario busca o saldo de um usuário pelo ID
func (repository Balanca) GetSaldoByUsuarioId(userID uint64) (float64, error) {
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

// DebitarSaldo debita o valor especificado do saldo do usuário
func (repository Balanca) DebitarSaldo(userID uint64, valor float64) error {
	var saldoAtual float64

	// Consulta o saldo atual do usuário
	query := `
        SELECT valor
        FROM balancas
        WHERE id_usuario = ? AND moeda = 'BRL'
    `
	err := repository.db.QueryRow(query, userID).Scan(&saldoAtual)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("usuário não possui saldo disponível")
		}
		return err
	}

	// Verificar se o saldo é suficiente para o débito
	if saldoAtual < valor {
		return errors.New("saldo insuficiente")
	}

	// Atualizar o saldo, debitando o valor
	queryUpdate := `
        UPDATE balancas
        SET valor = valor - ?
        WHERE id_usuario = ? AND moeda = 'BRL'
    `
	_, err = repository.db.Exec(queryUpdate, valor, userID)
	if err != nil {
		return err
	}

	return nil
}
