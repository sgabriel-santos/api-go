package repositories

import (
	"api/src/models"
	"database/sql"
	"fmt"
	"time"
)

// Representa um repositório de pagamentos
type Pagamento struct {
	db *sql.DB
}

// Fornece a lista de pagamentos do tipo 2 (QrCode) do usuário
func (repository Pagamento) GetPagamentosCopiaCola(userID uint64, timezone string, origin string) ([]models.PagamentoModel, error) {
	var rows *sql.Rows

	// Construção da consulta SQL com base no timezone e origin
	// Se possui timezone -> Realiza conversão, caso contrário apresenta o valor original do banco de dados
	// Se possui timezone e é mobile -> Os campos do pagamento “valor” e “status” são formatados para melhor visualização.
	data_cadastro := "p.data_cadastro"
	valor_formated := "CAST(p.valor as varchar(100))"
	status_text := "CAST(p.status as varchar(10))"

	if timezone != "" {
		timezoneOffset, err := getTimezoneOffset(timezone)
		if err != nil {
			return nil, err
		}
		data_cadastro = fmt.Sprintf(`CONVERT_TZ(p.data_cadastro,'+00:00','%s')`, timezoneOffset)

		if origin == "mobile" {
			valor_formated = `FORMAT(p.valor, 2, 'de_DE')`
			status_text = getStatusCaseStatement()
		}
	}

	query := fmt.Sprintf(`
		SELECT
			p.*,
			%s AS data_cadastro,
			%s AS valor_formatted,
			%s AS status_text
		FROM pagamentos p
		WHERE p.tipo = 2 AND p.id_usuario = ?
		ORDER BY p.id DESC`,
		data_cadastro, valor_formated, status_text)

	rows, err := repository.db.Query(query, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pagamentos []models.PagamentoModel

	for rows.Next() {
		var nomeBeneficiario, banco, agencia, ccDig, tipoConta, tipoPessoa, cpfCnpj, data, dataCode, reference, addressKey, description, response, statusText, valorFormatted sql.NullString
		var valor sql.NullFloat64
		var status sql.NullInt16
		var dataCadastro sql.NullTime
		var pagamento models.PagamentoModel

		// Scan dos campos retornados
		err = rows.Scan(
			&pagamento.ID,
			&pagamento.IDUsuario,
			&pagamento.Tipo,
			&nomeBeneficiario,
			&banco,
			&agencia,
			&ccDig,
			&tipoConta,
			&tipoPessoa,
			&cpfCnpj,
			&valor,
			&data,
			&dataCode,
			&reference,
			&addressKey,
			&description,
			&response,
			&status,
			&pagamento.DataCadastro,
			&dataCadastro,
			&valorFormatted,
			&statusText,
		)
		if err != nil {
			return nil, err
		}

		pagamento.DataCadastro = nullTimeToTime(dataCadastro)
		pagamento.NomeBeneficiario = nullStringToPointer(nomeBeneficiario)
		pagamento.Banco = nullStringToPointer(banco)
		pagamento.Agencia = nullStringToPointer(agencia)
		pagamento.CCDig = nullStringToPointer(ccDig)
		pagamento.TipoConta = nullStringToPointer(tipoConta)
		pagamento.TipoPessoa = nullStringToPointer(tipoPessoa)
		pagamento.CPFCNPJ = nullStringToPointer(cpfCnpj)
		pagamento.Valor = nullStringToPointer(valorFormatted)
		pagamento.Data = nullStringToPointer(data)
		pagamento.DataCode = nullStringToPointer(dataCode)
		pagamento.Reference = nullStringToPointer(reference)
		pagamento.AddressKey = nullStringToPointer(addressKey)
		pagamento.Description = nullStringToPointer(description)
		pagamento.Response = nullStringToPointer(response)
		pagamento.DataCadastro = nullTimeToTime(dataCadastro)
		pagamento.Status = nullStringToPointer(statusText)

		pagamentos = append(pagamentos, pagamento)
	}

	return pagamentos, nil
}

// Cria um repositório de usuários
func NewPagamentoRepository(db *sql.DB) *Pagamento {
	return &Pagamento{db}
}

// Função para obter o offset do timezone
func getTimezoneOffset(timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}
	_, offset := time.Now().In(loc).Zone()
	hoursOffset := offset / 3600
	minutesOffset := (offset % 3600) / 60
	return fmt.Sprintf("%+03d:%02d", hoursOffset, minutesOffset), nil
}

// Função para construir o CASE statement para o status
func getStatusCaseStatement() string {
	return `
        CASE
            WHEN p.status = 0 THEN 'Pendente'
            WHEN p.status = 1 THEN 'Processando'
            WHEN p.status = 2 THEN 'Concluído'
            WHEN p.status = 3 THEN 'Cancelado'
            WHEN p.status = 4 THEN 'Em Análise'
            ELSE 'Desconhecido'
        END`
}

// Função para converter sql.NullString para *string
func nullStringToPointer(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// Função para converter sql.NullFloat64 para *float64
func nullFloat64ToPointer(nf sql.NullFloat64) *float64 {
	if nf.Valid {
		return &nf.Float64
	}
	return nil
}

// Função para converter sql.NullTime para time.Time
func nullTimeToTime(nt sql.NullTime) time.Time {
	if nt.Valid {
		return nt.Time
	}
	return time.Time{} // Retorna o valor zero de time.Time
}

// CreatePagamento insere um novo pagamento no banco de dados
func (repository Pagamento) CreatePagamento(pagamento models.PagamentoModel) error {
	// Preparar a query de inserção
	query := `
        INSERT INTO pagamentos (id_usuario, tipo, valor, nome_beneficiario, data, data_code, reference, description, status, data_cadastro)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	// Preparar a inserção
	stmt, err := repository.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Executar a inserção com os parâmetros fornecidos
	_, err = stmt.Exec(
		pagamento.IDUsuario,
		pagamento.Tipo,
		pagamento.Valor,
		pagamento.NomeBeneficiario,
		pagamento.Data,
		pagamento.DataCode,
		pagamento.Reference,
		pagamento.Description,
		pagamento.Status,
		pagamento.DataCadastro,
	)
	if err != nil {
		return err
	}

	// Retornar nil se tudo ocorrer bem
	return nil
}

// GetPagamento busca um pagamento pelo userID e reference
func (repository Pagamento) GetPagamento(userID uint64, reference string) (*models.PagamentoModel, error) {
	var pagamento models.PagamentoModel

	query := `
        SELECT id_usuario, data, data_code, data_cadastro
        FROM pagamentos
        WHERE id_usuario = ? AND tipo = 2 AND reference = ?`

	// Executar a consulta no banco de dados
	row := repository.db.QueryRow(query, userID, reference)

	// Fazer o scan dos resultados
	err := row.Scan(
		&pagamento.IDUsuario,
		&pagamento.Data,
		&pagamento.DataCode,
		&pagamento.DataCadastro,
	)

	if err != nil {
		// Se nenhum registro for encontrado, retornar nil (sem erro)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		// Retornar o erro se ocorrer algum problema durante a consulta
		return nil, err
	}

	// Retornar o pagamento encontrado
	return &pagamento, nil
}
