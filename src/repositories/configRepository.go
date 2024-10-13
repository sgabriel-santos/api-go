package repositories

import (
	"api/src/models"
	"database/sql"
)

// Representa um repositório de config
type Config struct {
	db *sql.DB
}

// Cria um repositório de usuários
func NewConfigRepository(db *sql.DB) *Config {
	return &Config{db}
}

func (repository Config) GetConfigById(id_config int) (models.ConfigModel, error) {
	var config models.ConfigModel

	err := repository.db.QueryRow(`SELECT * FROM config WHERE id = ?`, id_config).Scan(
		&config.ID,
		&config.ClientID,
		&config.ClientSecret,
		&config.FeeUSDTBRLConvert,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.ConfigModel{}, nil // Config não encontrado
		}
		return models.ConfigModel{}, err // Outro erro ocorreu
	}

	return config, nil
}
