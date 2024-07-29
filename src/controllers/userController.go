package controllers

import (
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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
	name := strings.ToLower(r.URL.Query().Get("user"))
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

// Busca um usuário salvo no database pelo id
func GetUserById(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	userID, erro := strconv.ParseUint(parametros["userId"], 10, 64)
	if erro != nil {
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
	user, erro := repository.GetUserById(userID)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}
