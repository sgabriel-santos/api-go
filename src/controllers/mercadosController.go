package controllers

import (
	"api/src/database"
	"api/src/repositories"
	"api/src/responses"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// ListarMercados lida com o endpoint /mercados
func ListarMercados(w http.ResponseWriter, r *http.Request) {

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

	// Conectar ao banco de dados
	db, err := database.Connect()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	// Repositório de mercados
	repoMercado := repositories.NewMercadoRepository(db)
	userRepository := repositories.NewUserRepository(db)

	// Verificar se o usuário está ativo
	usuario, err := userRepository.GetUserByID(uint64(userId))
	if err != nil || usuario.Status != 1 {
		responses.Erro(w, http.StatusUnauthorized, errors.New("usuário inativo ou token inválido"))
		return
	}

	// Listar todos os mercados
	mercados, err := repoMercado.ListarTodosMercados()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Se nenhum mercado foi encontrado
	if len(mercados) == 0 {
		responses.JSON(w, http.StatusNotFound, map[string]string{"status": "Nenhum mercado encontrado"})
		return
	}

	// Responder com a lista de mercados
	responses.JSON(w, http.StatusOK, map[string]interface{}{"data": mercados})
}
