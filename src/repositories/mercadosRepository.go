package repositories

import (
	"api/src/models"
	"database/sql"
	"errors"
)

// Representa um repositório de Mercado
type Mercado struct {
	db *sql.DB
}

// Cria um repositório de usuários
func NewMercadoRepository(db *sql.DB) *Mercado {
	return &Mercado{db}
}

// ListarTodosMercados lista todos os mercados disponíveis no banco de dados
func (repository Mercado) ListarTodosMercados() ([]models.MercadoModel, error) {
	var mercados []models.MercadoModel

	query := `
        SELECT id, id_symbol, symbol, base, quote, min_amount, min_decimal, dec_precision, 
               min_precision, ask, bid, last
        FROM mercados
        ORDER BY id_symbol ASC
    `
	rows, err := repository.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var mercado models.MercadoModel
		err := rows.Scan(
			&mercado.ID,
			&mercado.IDSymbol,
			&mercado.Symbol,
			&mercado.Base,
			&mercado.Quote,
			&mercado.MinAmount,
			&mercado.MinDecimal,
			&mercado.DecPrecision,
			&mercado.MinPrecision,
			&mercado.Ask,
			&mercado.Bid,
			&mercado.Last,
		)
		if err != nil {
			return nil, err
		}
		mercados = append(mercados, mercado)
	}

	return mercados, nil
}

// GetMercado busca as informações de mercado para a conversão de uma moeda para outra
func (repository Mercado) GetMercado(moedaDe, moedaPara string) (models.MercadoModel, error) {
	var mercado models.MercadoModel
	query := `
        SELECT id, id_symbol, symbol, base, quote, min_amount, min_decimal, dec_precision, min_precision, ask, bid, last
        FROM mercados
        WHERE base = ? AND quote = ?
    `
	err := repository.db.QueryRow(query, moedaDe, moedaPara).Scan(
		&mercado.ID,
		&mercado.IDSymbol,
		&mercado.Symbol,
		&mercado.Base,
		&mercado.Quote,
		&mercado.MinAmount,
		&mercado.MinDecimal,
		&mercado.DecPrecision,
		&mercado.MinPrecision,
		&mercado.Ask,
		&mercado.Bid,
		&mercado.Last,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return mercado, errors.New("mercado não encontrado")
		}
		return mercado, err
	}

	return mercado, nil
}
