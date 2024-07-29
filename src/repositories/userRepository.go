package repositories

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

// Representa um repositório de users
type User struct {
	db *sql.DB
}

// Cria um repositório de usuários
func NewUserRepository(db *sql.DB) *User {
	return &User{db}
}

// Insere um usuário no database de dados
func (repository User) CreateUser(user models.UserModel) (uint64, error) {
	statement, erro := repository.db.Prepare(
		"insert into users (name, email, password) values(?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	// Garantir que a instrução SQL seja fechada após o uso
	defer statement.Close()

	response, erro := statement.Exec(user.Name, user.Email, user.Password)
	if erro != nil {
		return 0, erro
	}

	lastInsertId, erro := response.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(lastInsertId), nil
}

// Busca todos os usuários que atendem um filtro de name
func (repository User) SearchUserByName(name string) ([]models.UserModel, error) {
	name = fmt.Sprintf("%%%s%%", name) // %name%

	rows, erro := repository.db.Query(
		"select id, name, email from users where name LIKE ?",
		name,
	)

	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var users []models.UserModel

	for rows.Next() {
		var user models.UserModel

		if erro = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
		); erro != nil {
			return nil, erro
		}

		users = append(users, user)
	}

	return users, nil
}

// Busca usuário do database de dados pelo ID
func (repository User) GetUserById(ID uint64) (models.UserModel, error) {
	rows, erro := repository.db.Query(
		"select id, name, email from users where id = ?",
		ID,
	)
	if erro != nil {
		return models.UserModel{}, erro
	}
	defer rows.Close()

	var user models.UserModel

	if rows.Next() {
		if erro = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
		); erro != nil {
			return models.UserModel{}, erro
		}
	}

	return user, nil
}

// Busca um usuário por email e retorna o seu id e password com hash
func (repository User) GetUserByEmail(email string) (models.UserModel, error) {
	row, erro := repository.db.Query("select id, password from users where email = ?", email)
	if erro != nil {
		return models.UserModel{}, erro
	}
	defer row.Close()

	var user models.UserModel

	if row.Next() {
		if erro = row.Scan(&user.ID, &user.Password); erro != nil {
			return models.UserModel{}, erro
		}
	}

	return user, nil

}
