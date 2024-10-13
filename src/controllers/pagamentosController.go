package controllers

import (
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Data struct {
	Error     string `json:"error"`
	Reference string `json:"reference"`
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}

type Response struct {
	Data Data `json:"data"`
}

// Insere um usuário no database de dados
func CopiaECola(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	origin := query.Get("origin")

	if origin == "" {
		responses.Erro(w, http.StatusConflict, errors.New("parâmetro origin não foi passado na rota"))
		return
	}

	id := r.Header.Get("id")
	userId, err := strconv.Atoi(id)

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	usuarioRepository := repositories.NewUserRepository(db)
	pagamentoRepository := repositories.NewPagamentoRepository(db)

	user, err := usuarioRepository.GetUserByID(uint64(userId))

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if user.ID == 0 {
		responses.Erro(w, http.StatusNotFound, fmt.Errorf("usuário com id %d não encontrado", userId))
		return
	}

	var pagamentos []models.PagamentoModel

	timezone := ""
	if user.Timezone != nil {
		timezone = *user.Timezone
	}

	pagamentos, err = pagamentoRepository.GetPagamentosCopiaCola(uint64(userId), timezone, origin)

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Verifica se os pagamentos foram encontrados
	if len(pagamentos) == 0 {
		responses.JSON(w, http.StatusOK, map[string]string{"message": "Nenhum pagamento encontrado."})
		return
	}

	responses.JSON(w, http.StatusOK, pagamentos)
}
