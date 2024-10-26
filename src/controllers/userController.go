package controllers

import (
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Insere um usuário no database de dados
func CreateUser(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var user models.UserModel
	if erro = json.Unmarshal(corpoRequest, &user); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = user.Prepare("cadastro"); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)
	user.ID, erro = repository.CreateUser(user)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusCreated, user)
}

// Somente verifica credenciais do usuário
func VerifyUser(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, true)
}

// Busca todos os usuários salvos no database de dados
func SearchUserByName(w http.ResponseWriter, r *http.Request) {
	name := strings.ToLower(r.URL.Query().Get("name"))
	db, erro := database.Connect()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)
	users, erro := repository.SearchUserByName(name)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

// ListarTimezones lista os fusos horários disponíveis e retorna o fuso horário do usuário
func ListarTimezones(w http.ResponseWriter, r *http.Request) {
	// Conectar ao banco de dados
	db, err := database.Connect()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	// Repositório para interagir com o banco de dados
	repository := repositories.NewUserRepository(db)

	id := r.Header.Get("id")
	userId, err := strconv.Atoi(id)

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Verificar se o usuário está ativo e obter o timezone dele
	usuario, err := repository.GetUserByID(uint64(userId))
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	if usuario.Status != 1 {
		responses.Erro(w, http.StatusUnauthorized, errors.New("token de autenticação é inválido"))
		return
	}

	// Obter os timezones únicos dos usuários cadastrados
	timezones, err := repository.GetUniqueTimezones()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Calcular as diferenças de cada timezone em relação ao UTC/GMT
	timezoneDiffs := calcularDiferencaTimezones(timezones)

	// Retornar o fuso horário do usuário junto com a lista de todos os fusos horários
	response := map[string]interface{}{
		"timezone": usuario.Timezone, // O timezone atual do usuário
		"data":     timezoneDiffs,    // Lista de fusos horários com suas diferenças
	}
	responses.JSON(w, http.StatusOK, response)
}

// calcularDiferencaTimezones calcula a diferença de cada timezone em relação ao UTC/GMT
func calcularDiferencaTimezones(timezones []string) map[string]string {
	timezoneDiffs := make(map[string]string)
	currentTime := time.Now()

	// Calcular a diferença de cada timezone em relação ao UTC
	for _, zone := range timezones {
		loc, err := time.LoadLocation(zone)
		if err != nil {
			continue // Ignorar timezones inválidos
		}
		diff := currentTime.In(loc).Format("-07:00")
		timezoneDiffs[zone] = "UTC/GMT " + diff + " - " + zone
	}

	return timezoneDiffs
}

// AtualizarFusoHorario lida com o endpoint /config_timezone
func AtualizarFusoHorario(w http.ResponseWriter, r *http.Request) {
	// Conectar ao banco de dados
	db, err := database.Connect()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	// Repositório de usuários
	userRepository := repositories.NewUserRepository(db)

	id := r.Header.Get("id")

	if id == "" {
		responses.Erro(w, http.StatusInternalServerError, fmt.Errorf("id do usuário não foi passado no header"))
		return
	}

	userId, err := strconv.Atoi(id)

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Verificar se o usuário está ativo e obter o timezone dele
	usuario, err := userRepository.GetUserByID(uint64(userId))
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	if usuario.Status != 1 {
		responses.Erro(w, http.StatusUnauthorized, errors.New("usuário Inativo"))
		return
	}

	// Decodificar o JSON da requisição
	var dados struct {
		Timezone string `json:"timezone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&dados); err != nil {
		responses.Erro(w, http.StatusBadRequest, errors.New("dados inválidos"))
		return
	}

	// Atualizar o fuso horário do usuário
	err = userRepository.AtualizarTimezone(uint64(userId), dados.Timezone)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Retornar a confirmação
	responses.JSON(w, http.StatusOK, map[string]string{"status": "Fuso horário atualizado com sucesso"})
}
