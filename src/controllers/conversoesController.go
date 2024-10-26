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
	"time"
)

// ListarConversoes lida com o endpoint /conversoes
func ListarConversoes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	origin := query.Get("origin")

	if origin == "" {
		responses.Erro(w, http.StatusConflict, errors.New("parâmetro origin não foi passado na rota"))
		return
	}

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

	// Repositório de conversões
	userRepository := repositories.NewUserRepository(db)
	conversaoRepository := repositories.NewConversaoRepository(db)

	// Verificar se o usuário está ativo
	user, err := userRepository.GetUserByID(uint64(userId))
	if err != nil || user.Status != 1 {
		responses.Erro(w, http.StatusUnauthorized, errors.New("usuário inativo ou token inválido"))
		return
	}

	// Listar todas as conversões do usuário
	conversoes, err := conversaoRepository.ListarConversoesPorUsuario(uint64(userId))
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Se nenhuma conversão foi encontrada
	if len(conversoes) == 0 {
		responses.JSON(w, http.StatusNotFound, map[string]string{"status": "Nenhuma conversão encontrada"})
		return
	}

	// Ajustar a exibição do fuso horário
	if *user.Timezone != "" {
		conversoes = ajustarConversoesTimezone(conversoes, *user.Timezone)
	}

	if origin == "mobile" {
		conversoes = ajustarStatusConversoes(conversoes)
	}

	// Responder com a lista de conversões
	responses.JSON(w, http.StatusOK, map[string]interface{}{"data": conversoes})
}

// ajustarConversoesTimezone ajusta as datas de depósito para o fuso horário do usuário
func ajustarConversoesTimezone(conversoes []models.ConversaoModel, timezone string) []models.ConversaoModel {
	loc, _ := time.LoadLocation(timezone)
	for i := range conversoes {
		conversoes[i].DataCadastro = conversoes[i].DataCadastro.In(loc)
	}
	return conversoes
}

// ajustarStatusConversoes converte o valor do status para uma descrição por extenso
func ajustarStatusConversoes(conversoes []models.ConversaoModel) []models.ConversaoModel {
	for i := range conversoes {
		switch conversoes[i].Status {
		case 0:
			conversoes[i].StatusDescricao = "Pendente"
		case 1:
			conversoes[i].StatusDescricao = "Em andamento"
		case 2:
			conversoes[i].StatusDescricao = "Concluído"
		case 3:
			conversoes[i].StatusDescricao = "Cancelado"
		default:
			conversoes[i].StatusDescricao = "Desconhecido"
		}
	}
	return conversoes
}

// ConverterMoeda lida com o endpoint /converter
func ConverterMoeda(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	origin := query.Get("origin")

	if origin == "" {
		responses.Erro(w, http.StatusConflict, errors.New("parâmetro origin não foi passado na rota"))
		return
	}

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

	// Extrair dados da requisição
	valorDe, err := strconv.ParseFloat(r.FormValue("valor_de"), 64)
	if err != nil {
		responses.Erro(w, http.StatusBadRequest, errors.New("valor de conversão inválido"))
		return
	}

	moedaDe := r.FormValue("moeda_de")
	moedaPara := r.FormValue("moeda_para")

	if moedaDe == "" {
		responses.Erro(w, http.StatusInternalServerError, fmt.Errorf("moeda_de não foi passado"))
		return
	}

	if moedaPara == "" {
		responses.Erro(w, http.StatusInternalServerError, fmt.Errorf("moeda_para não foi passado"))
		return
	}

	userRepository := repositories.NewUserRepository(db)
	mercadosRepository := repositories.NewMercadoRepository(db)
	conversaoRepository := repositories.NewConversaoRepository(db)

	// Verificar saldo do usuário
	saldo, err := userRepository.GetSaldoByUsuario(uint64(userId), moedaDe)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	if saldo < valorDe {
		responses.Erro(w, http.StatusForbidden, errors.New("saldo insuficiente para realizar a conversão"))
		return
	}

	// Buscar preço de mercado para a conversão
	mercado, err := mercadosRepository.GetMercado(moedaDe, moedaPara)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, errors.New("falha ao consultar o preço de mercado"))
		return
	}

	conversao, err := conversaoRepository.ProcessarConversao(uint64(userId), valorDe, mercado)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Responder com os dados da conversão
	responses.JSON(w, http.StatusOK, conversao)
}
