package repositories

import (
	"api/src/models"
	"database/sql"
	"time"
)

// Representa um repositório de Endereco
type Endereco struct {
	db *sql.DB
}

// Cria um repositório de usuários
func NewEnderecoRepository(db *sql.DB) *Endereco {
	return &Endereco{db}
}

// GetEnderecoUSDTByUserID busca o endereço USDT do usuário no banco de dados
func (repository Endereco) GetEnderecoUSDTByUserID(userID uint64) (*models.EnderecoModel, error) {
	var endereco models.EnderecoModel

	query := `
        SELECT id_usuario, moeda, rede, endereco, dados 
        FROM enderecos 
        WHERE id_usuario = ? AND moeda = 'USDT' AND rede = 'BEP20'
    `
	err := repository.db.QueryRow(query, userID).Scan(
		&endereco.IDUsuario,
		&endereco.Moeda,
		&endereco.Rede,
		&endereco.Endereco,
		&endereco.Dados,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &endereco, nil
}

// SaveEnderecoUSDT salva um novo endereço USDT no banco de dados
func (repository Endereco) SaveEnderecoUSDT(userID uint64, enderecoUSDT *models.EnderecoModel) error {
	query := `
        INSERT INTO enderecos (id_usuario, moeda, rede, endereco, dados, data_cadastro) 
        VALUES (?, 'USDT', 'BEP20', ?, ?, ?)
    `
	_, err := repository.db.Exec(query, userID, enderecoUSDT.Endereco, enderecoUSDT.Dados, time.Now())
	return err
}
