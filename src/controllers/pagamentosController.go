package controllers

import (
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"api/src/security"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
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

// PagamentoQRCodeDecode decodifica um QR Code de pagamento
func PagamentoQRCodeDecode(w http.ResponseWriter, r *http.Request) {
	// Ler corpo da requisição (espera-se que contenha o QR Code)
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Estrutura para ler os dados da requisição (espera-se que tenha o campo 'codigo')
	var dados struct {
		Codigo string `json:"codigo"`
	}
	if err := json.Unmarshal(requestBody, &dados); err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
		return
	}

	id := r.Header.Get("id")
	userId, err := strconv.Atoi(id)

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Conexão com o banco de dados
	db, err := database.Connect()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	// Repositório para interagir com o banco de dados
	usuarioRepository := repositories.NewUserRepository(db)
	configRepository := repositories.NewConfigRepository(db)
	pagamentoLogRepository := repositories.NewPagamentoLogRepository(db)

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

	// Obter as configurações (client_id e client_secret) do banco de dados
	config, err := configRepository.GetConfigById(1)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Autenticar no serviço externo
	loginResponse, err := authenticateWithExternalService(config.ClientID, config.ClientSecret)
	if err != nil || loginResponse.Token == "" {
		responses.Erro(w, http.StatusBadRequest, errors.New("620: Requisição inválida"))
		return
	}

	// Decodificar o QR Code utilizando o serviço externo
	decodedData, err := decodeQRCodeWithExternalService(loginResponse.Token, dados.Codigo)
	if err != nil {
		responses.Erro(w, http.StatusBadRequest, errors.New("702: Requisição inválida"))
		return
	}

	// Armazenar as informações decodificadas no banco de dados (tabela pagamentos_log)
	logData := models.PagamentoLogModel{
		IDUsuario:    uint64(userId),
		Method:       "pagamentos_qrcode",
		Data:         &decodedData.Data.Data,
		DataID:       decodedData.Data.Reference,
		Code:         time.Now().Format("20060102150405"),
		DataCadastro: time.Now(), // Define o timestamp atual para DataCadastro
	}

	if err := pagamentoLogRepository.CreatePagamentoLog(logData); err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Retornar os dados decodificados ao usuário
	responses.JSON(w, http.StatusOK, decodedData)
}

type LoginResponse struct {
	Token string `json:"token"`
}

func authenticateWithExternalService(clientID, clientSecret string) (LoginResponse, error) {
	url := "https://api-lastmile.banco.green/v1/user/login"
	requestBody, _ := json.Marshal(map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return LoginResponse{}, err
	}
	defer resp.Body.Close()

	var loginResponse LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		return LoginResponse{}, err
	}

	return loginResponse, nil
}

type QRCodeData struct {
	Error     string `json:"error,omitempty"`
	Reference string `json:"reference,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Data      string `json:"data,omitempty"`
}

// Modelo para a resposta completa da API
type DecodedQRCode struct {
	Data QRCodeData `json:"data"`
}

func decodeQRCodeWithExternalService(token, codigo string) (DecodedQRCode, error) {
	url := "https://api-lastmile.banco.green/v1/pix/qrcode_decode"
	requestBody, _ := json.Marshal(map[string]string{
		"token": token,
		"code":  codigo,
	})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return DecodedQRCode{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return DecodedQRCode{}, err
	}
	defer resp.Body.Close()

	// Ler e decodificar a resposta no modelo DecodedQRCode
	var decodedQRCode DecodedQRCode
	if err := json.NewDecoder(resp.Body).Decode(&decodedQRCode); err != nil {
		return DecodedQRCode{}, err
	}

	// Verificar se há algum erro dentro do campo 'data.error'
	if decodedQRCode.Data.Error != "" {
		return DecodedQRCode{}, errors.New(decodedQRCode.Data.Error)
	}

	// Retornar os dados decodificados
	return decodedQRCode, nil
}

// PagamentoQRCode realiza o pagamento de um PIX via copia e cola
func PagamentoQRCode(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	// Repositório para interagir com o banco de dados
	usuarioRepository := repositories.NewUserRepository(db)
	pagamentoLogRepository := repositories.NewPagamentoLogRepository(db)
	balancaRepository := repositories.NewBalancaRepository(db)
	pagamentoRepository := repositories.NewPagamentoRepository(db)

	id := r.Header.Get("id")
	userId, err := strconv.Atoi(id)

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

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

	// Ler os dados da requisição
	var dados struct {
		Timestamp   string  `json:"timestamp"`
		Reference   string  `json:"reference"`
		Value       float64 `json:"value"`
		Pin         string  `json:"pin"`
		Description string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&dados); err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
		return
	}

	// Verificar se os dados necessários estão presentes
	if dados.Timestamp == "" || dados.Reference == "" || dados.Value <= 0 || dados.Pin == "" {
		responses.Erro(w, http.StatusBadRequest, errors.New("dados incorretos"))
		return
	}

	if !security.VerifyPassword(user.Pin, dados.Pin) {
		responses.Erro(w, http.StatusUnauthorized, errors.New("PIN informado está incorreto"))
		return
	}

	// Verificar se já existe um pagamento logado para o mesmo código (timestamp) e referência
	pagamentoLog, err := pagamentoLogRepository.GetPagamentoLog(uint64(userId), dados.Timestamp, dados.Reference)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if pagamentoLog == nil {
		responses.Erro(w, http.StatusBadRequest, errors.New("erro: Dados Incorretos"))
		return
	}

	pagamento, err := pagamentoRepository.GetPagamento(uint64(userId), dados.Reference)

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	if pagamento != nil {
		responses.Erro(w, http.StatusBadRequest, errors.New("erro: Pagamento Duplicado"))
		return
	}

	// Verificar saldo do usuário
	saldo, err := balancaRepository.GetSaldoByUsuarioId(uint64(userId))
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	if saldo < dados.Value {
		responses.Erro(w, http.StatusBadRequest, errors.New("erro: Saldo insuficiente"))
		return
	}

	// Inserir pagamento
	status := new(string)
	*status = "1"

	value := new(string)
	*value = fmt.Sprintf("%.2f", dados.Value)

	novo_pagamento := models.PagamentoModel{
		IDUsuario:        userId,
		Tipo:             2, // Tipo de pagamento PIX
		Valor:            value,
		NomeBeneficiario: &user.Nome, // Supomos que o nome vem do usuário
		Reference:        &dados.Reference,
		Description:      &dados.Description,
		Status:           status, // Status inicial
		DataCadastro:     time.Now(),
	}

	if err := pagamentoRepository.CreatePagamento(novo_pagamento); err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	// Debitar saldo do usuário
	if err := balancaRepository.DebitarSaldo(uint64(userId), dados.Value); err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{"message": "Pagamento realizado com sucesso"})
}
