package repositories

import (
	"api/src/models"
	"database/sql"
)

// Representa um repositório de Deposito
type Deposito struct {
	db *sql.DB
}

// Cria um repositório de usuários
func NewDepositoRepository(db *sql.DB) *Deposito {
	return &Deposito{db}
}

// SaveEnderecoUSDT salva um novo endereço USDT no banco de dados
func (repository Deposito) SaveEnderecoUSDT(userID uint64, enderecoUSDT *models.EnderecoModel) error {
	query := `
        INSERT INTO enderecos (id_usuario, moeda, endereco, rede, dados) 
        VALUES (?, ?, ?, ?, ?)
    `
	_, err := repository.db.Exec(query, userID, "USDT", enderecoUSDT.Endereco, "BEP20", enderecoUSDT.Dados)
	return err
}

// GetDepositosUSDTByUserID lista todos os depósitos em USDT do usuário
func (repository Deposito) GetDepositosUSDTByUserID(userID uint64) ([]models.DepositoModel, error) {
	var depositos []models.DepositoModel

	query := `
        SELECT id, id_usuario, referencia, moeda, endereco, valor, transacao, retorno, pix_imagem, status, data_concluido, data_cadastro 
        FROM depositos 
        WHERE id_usuario = ? AND moeda = 'USDT'
        ORDER BY id DESC
    `
	rows, err := repository.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var deposito models.DepositoModel
		err := rows.Scan(
			&deposito.ID,
			&deposito.IDUsuario,
			&deposito.Referencia,
			&deposito.Moeda,
			&deposito.Endereco,
			&deposito.Valor,
			&deposito.Transacao,
			&deposito.Retorno,
			&deposito.PixImagem,
			&deposito.Status,
			&deposito.DataConcluido,
			&deposito.DataCadastro,
		)
		if err != nil {
			return nil, err
		}
		depositos = append(depositos, deposito)
	}

	return depositos, nil
}
