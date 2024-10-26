package repositories

import (
	"api/src/models"
	"database/sql"
	"time"
)

// Conversao representa o repositório de conversões
type Conversao struct {
	db *sql.DB
}

// NewConversaoRepository cria um novo repositório de conversões
func NewConversaoRepository(db *sql.DB) *Conversao {
	return &Conversao{db}
}

// ListarConversoesPorUsuario lista todas as conversões realizadas pelo usuário
func (repository Conversao) ListarConversoesPorUsuario(userID uint64) ([]models.ConversaoModel, error) {
	var conversoes []models.ConversaoModel

	query := `
        SELECT id, id_usuario, simbolo, de, para, valor, valor_final, preco, bid, ask, last, 
               liquidar, liquidar_retorno, status, data_cadastro
        FROM conversoes
        WHERE id_usuario = ?
        ORDER BY data_cadastro DESC
    `
	rows, err := repository.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var conversao models.ConversaoModel
		err := rows.Scan(
			&conversao.ID,
			&conversao.IDUsuario,
			&conversao.Simbolo,
			&conversao.De,
			&conversao.Para,
			&conversao.Valor,
			&conversao.ValorFinal,
			&conversao.Preco,
			&conversao.Bid,
			&conversao.Ask,
			&conversao.Last,
			&conversao.Liquidar,
			&conversao.LiquidarRetorno,
			&conversao.Status,
			&conversao.DataCadastro,
		)
		if err != nil {
			return nil, err
		}
		conversoes = append(conversoes, conversao)
	}

	return conversoes, nil
}

// Registrar e processar a conversão
func (repository Conversao) ProcessarConversao(userID uint64, valorDe float64, mercado models.MercadoModel) (models.ConversaoModel, error) {
	tx, err := repository.db.Begin() // Iniciar uma transação
	if err != nil {
		return models.ConversaoModel{}, err
	}

	// Debitar o saldo da moeda de origem (USDT, por exemplo)
	queryDebitar := `UPDATE balancas SET valor = valor - ? WHERE id_usuario = ? AND moeda = ?`
	_, err = tx.Exec(queryDebitar, valorDe, userID, mercado.Base)
	if err != nil {
		tx.Rollback()
		return models.ConversaoModel{}, err
	}

	// Calcular o valor final após a conversão, incluindo taxas (se houver)
	valorFinal := valorDe * mercado.Ask

	// Creditar o saldo na moeda de destino (BRL, por exemplo)
	queryCreditar := `UPDATE balancas SET valor = valor + ? WHERE id_usuario = ? AND moeda = ?`
	_, err = tx.Exec(queryCreditar, valorFinal, userID, mercado.Quote)
	if err != nil {
		tx.Rollback()
		return models.ConversaoModel{}, err
	}

	// Registrar a conversão
	queryInserirConversao := `
        INSERT INTO conversoes (id_usuario, simbolo, de, para, valor, valor_final, preco, bid, ask, last, liquidar, status, data_cadastro)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, 2, NOW())
    `
	result, err := tx.Exec(queryInserirConversao, userID, mercado.IDSymbol, mercado.Base, mercado.Quote, valorDe, valorFinal, mercado.Ask, mercado.Bid, mercado.Ask, mercado.Last)
	if err != nil {
		tx.Rollback()
		return models.ConversaoModel{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return models.ConversaoModel{}, err
	}

	err = tx.Commit() // Confirmar a transação
	if err != nil {
		return models.ConversaoModel{}, err
	}

	return models.ConversaoModel{
		ID:           uint64(id),
		IDUsuario:    userID,
		Simbolo:      mercado.IDSymbol,
		De:           mercado.Base,
		Para:         mercado.Quote,
		Valor:        valorDe,
		ValorFinal:   valorFinal,
		Preco:        mercado.Ask,
		Bid:          mercado.Bid,
		Ask:          mercado.Ask,
		Last:         mercado.Last,
		Liquidar:     0, // Pendente
		Status:       2, // Concluído
		DataCadastro: time.Now(),
	}, nil
}
