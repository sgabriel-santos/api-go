package database

import (
	"api/src/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Driver
)

// Abre a conexão com o database de dados e a retorna
func Connect() (*sql.DB, error) {
	db, erro := sql.Open("mysql", config.DatabaseConnectionString)
	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}

	return db, nil

}
