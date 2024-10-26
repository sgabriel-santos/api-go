package repositories

import (
	"api/src/models"
	"database/sql"
	"fmt"
	"time"
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
		`INSERT INTO usuarios 
		(uid, tipo, email, senha, nome, cpf, cnpj, rg, orgao_emissor, telefone, celular, pin, status, data_cadastro) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
	)
	if erro != nil {
		return 0, erro
	}
	// Garantir que a instrução SQL seja fechada após o uso
	defer statement.Close()

	var cpf *string = nil
	var cnpj *string = nil
	if user.CpfCnpj == "cpf" {
		cpf = &user.Cpf
	}

	if user.CpfCnpj == "cnpj" {
		cnpj = &user.Cnpj
	}

	dataCadastro := time.Now().Format("2006-01-02 15:04:05")

	response, erro := statement.Exec(
		user.UID,
		user.Tipo,
		user.Email,
		user.Senha,
		user.Nome,
		cpf,
		cnpj,
		user.Rg,
		user.OrgaoEmissor,
		user.Telefone,
		user.Celular,
		user.Pin,
		user.Status,
		dataCadastro,
	)
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
		"select id, nome, email from usuarios where nome LIKE ?",
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
			&user.Nome,
			&user.Email,
		); erro != nil {
			return nil, erro
		}

		users = append(users, user)
	}

	return users, nil
}

// Busca usuário do database de dados pelo ID
func (repository User) GetUserByID(userID uint64) (models.UserModel, error) {
	row := repository.db.QueryRow(`
        SELECT id, email, nome, timezone, status, pin
        FROM usuarios
        WHERE id = ?`, userID)

	var user models.UserModel
	var timezone sql.NullString

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Nome,
		&timezone,
		&user.Status,
		&user.Pin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil // Retorna usuário vazio se não encontrado
		}
		return user, err
	}

	if timezone.Valid {
		user.Timezone = &timezone.String
	}

	user.DataCadastro = nil
	return user, nil
}

// Busca um usuário por email e retorna o seu id e password com hash
func (repository User) GetUserByEmail(email string) (models.UserModel, error) {
	row, erro := repository.db.Query("select id, senha from usuarios where email = ?", email)
	if erro != nil {
		return models.UserModel{}, erro
	}
	defer row.Close()

	var user models.UserModel

	if row.Next() {
		if erro = row.Scan(&user.ID, &user.Senha); erro != nil {
			return models.UserModel{}, erro
		}
	}

	return user, nil
}

// GetUniqueTimezones busca todos os timezones únicos dos usuários
func (repository User) GetUniqueTimezones() ([]string, error) {
	var timezones []string

	// Consulta SQL para obter todos os timezones únicos
	query := `
        SELECT DISTINCT timezone
        FROM usuarios
        WHERE timezone IS NOT NULL AND timezone != ''
    `

	// Executar a consulta
	rows, err := repository.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Percorrer os resultados e adicionar os timezones na lista
	for rows.Next() {
		var timezone string
		if err := rows.Scan(&timezone); err != nil {
			return nil, err
		}
		timezones = append(timezones, timezone)
	}

	return timezones, nil
}

// GetSaldoByUsuario busca o saldo disponível para uma moeda específica do usuário
func (repository User) GetSaldoByUsuario(userID uint64, moeda string) (float64, error) {
	var saldo float64
	query := `SELECT valor FROM balancas WHERE id_usuario = ? AND moeda = ?`
	err := repository.db.QueryRow(query, userID, moeda).Scan(&saldo)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Se não houver saldo, retorna 0
		}
		return 0, err // Outro erro ocorreu
	}
	return saldo, nil
}

// AtualizarTimezone atualiza o fuso horário de um usuário
func (repository User) AtualizarTimezone(userID uint64, timezone string) error {
	query := `UPDATE usuarios SET timezone = ? WHERE id = ?`
	_, err := repository.db.Exec(query, timezone, userID)
	return err
}
