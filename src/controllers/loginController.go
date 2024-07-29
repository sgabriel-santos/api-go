package controllers

import (
	"api/src/authentication"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"api/src/security"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

// Responsável por autenticar um usuário na API
func Login(w http.ResponseWriter, r *http.Request) {
	requestBody, erro := io.ReadAll(r.Body)
	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var user models.UserModel
	if erro = json.Unmarshal(requestBody, &user); erro != nil {
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
	userSaveInDatabase, erro := repository.GetUserByEmail(user.Email)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if userSaveInDatabase.ID == 0 {
		responses.Erro(w, http.StatusNotFound, errors.New("usuário não encontrado"))
		return
	}

	if erro = security.VerifyPassword(userSaveInDatabase.Password, user.Password); erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	token, erro := authentication.CreateToken(userSaveInDatabase.ID)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	userID := strconv.FormatUint(userSaveInDatabase.ID, 10)

	responses.JSON(w, http.StatusOK, models.AuthenticationModel{ID: userID, Token: token})
}
