package middlewares

import (
	"api/src/authentication"
	"api/src/responses"
	"log"
	"net/http"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Realiza configurações necessárias para o log da aplicação
func ConfigureLogger() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "logs/api-go.log",
		MaxSize:    10,   // Tamanho máximo em megabytes antes da rotação do arquivo de log
		MaxBackups: 3,    // Número máximo de arquivos de log de backup
		MaxAge:     90,   // Número máximo de dias para reter os arquivos de log
		Compress:   true, // Se os backups antigos devem ser comprimidos usando gzip
	})
}

// Logger escreve informações da requisição no terminal
func Logger(nextFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.Host)
		nextFunction(w, r)
	}
}

// Verifica se o usuário fazendo a requisição está autenticado
func VerifyAuthentication(nextFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if erro := authentication.ValidateToken(r); erro != nil {
			responses.Erro(w, http.StatusUnauthorized, erro)
			return
		}
		nextFunction(w, r)
	}
}
