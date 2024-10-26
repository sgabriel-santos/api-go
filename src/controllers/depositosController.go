package controllers

import (
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// ListarDepositosUSDT lida com o endpoint /depositos/usdt
func ListarDepositosUSDT(w http.ResponseWriter, r *http.Request) {
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

	// Criar repositories
	usuarioRepository := repositories.NewUserRepository(db)
	depositoRepository := repositories.NewDepositoRepository(db)
	enderecoRepository := repositories.NewEnderecoRepository(db)
	configRepository := repositories.NewConfigRepository(db)

	user, err := usuarioRepository.GetUserByID(uint64(userId))

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if user.ID == 0 {
		responses.Erro(w, http.StatusNotFound, fmt.Errorf("usuário com id %d não encontrado", userId))
		return
	}

	if user.Status != 1 {
		responses.Erro(w, http.StatusNotFound, fmt.Errorf("usuário com id %d está inativo", userId))
		return
	}

	// Verificar se o usuário já possui um endereço USDT
	endereco, err := enderecoRepository.GetEnderecoUSDTByUserID(uint64(userId))
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Se o usuário não tem um endereço USDT, criá-lo chamando um serviço externo
	if endereco == nil {
		endereco, err = criarEnderecoUSDTExterno(uint64(userId), *configRepository, *depositoRepository)
		if err != nil {
			responses.Erro(w, http.StatusInternalServerError, err)
			return
		}
	}

	// Listar os depósitos do usuário
	depositos, err := depositoRepository.GetDepositosUSDTByUserID(uint64(userId))
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Ajustar a exibição do fuso horário
	if *user.Timezone != "" {
		depositos = ajustarDepositosTimezone(depositos, *user.Timezone)
	}

	if origin == "mobile" {
		depositos = ajustarStatusDepositos(depositos)
	}

	// Resposta final com o endereço e os depósitos
	response := map[string]interface{}{
		"endereco":  endereco,
		"depositos": depositos,
	}
	responses.JSON(w, http.StatusOK, response)
}

// criarEnderecoUSDTExterno cria um endereço USDT para o usuário chamando um serviço externo
func criarEnderecoUSDTExterno(userID uint64, configRepository repositories.Config, depositoRepository repositories.Deposito) (*models.EnderecoModel, error) {
	// Obter as configurações (client_id e client_secret) do banco de dados
	config, err := configRepository.GetConfigById(1)
	if err != nil {
		return nil, errors.New("falha ao buscar informações de client_id e client_secret do banco de dados")
	}

	// Autenticar com o serviço externo
	loginResponse, err := authenticateWithExternalService(config.ClientID, config.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// Requisição para criar o endereço USDT
	enderecoUSDT, err := criarEnderecoUSDT(loginResponse.Token, userID)
	if err != nil {
		return nil, err
	}

	// Salvar o endereço gerado no banco de dados
	err = depositoRepository.SaveEnderecoUSDT(userID, enderecoUSDT)
	if err != nil {
		return nil, err
	}

	return enderecoUSDT, nil
}

// ajustarDepositosTimezone ajusta as datas de depósito para o fuso horário do usuário
func ajustarDepositosTimezone(depositos []models.DepositoModel, timezone string) []models.DepositoModel {
	loc, _ := time.LoadLocation(timezone)
	for i := range depositos {
		depositos[i].DataCadastro = depositos[i].DataCadastro.In(loc)
	}
	return depositos
}

// ajustarStatusDepositos converte o valor do status para uma descrição por extenso
func ajustarStatusDepositos(depositos []models.DepositoModel) []models.DepositoModel {
	for i := range depositos {
		switch depositos[i].Status {
		case 0:
			depositos[i].StatusDescricao = "Pendente"
		case 1:
			depositos[i].StatusDescricao = "Em andamento"
		case 2:
			depositos[i].StatusDescricao = "Concluído"
		case 3:
			depositos[i].StatusDescricao = "Cancelado"
		default:
			depositos[i].StatusDescricao = "Desconhecido"
		}
	}
	return depositos
}

// criarEnderecoUSDT faz a requisição ao serviço externo para criar um endereço de depósito USDT
func criarEnderecoUSDT(token string, userID uint64) (*models.EnderecoModel, error) {
	url := "https://api-lastmile.banco.green/v1/usdt/address"
	data := map[string]interface{}{
		"token":      token,
		"projeto":    "brpay",
		"id_usuario": userID,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	dataResult, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("falha ao obter dados do serviço externo")
	}

	endereco := dataResult["endereco"].(string)
	dados, _ := json.Marshal(result)

	return &models.EnderecoModel{
		Endereco: endereco,
		Dados:    string(dados),
	}, nil
}
