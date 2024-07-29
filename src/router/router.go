package router

import (
	routes "api/src/router/routes"

	"github.com/gorilla/mux"
)

// Retorna um router com as rotas configuradas
func GenerateRoutes() *mux.Router {
	r := mux.NewRouter()
	return routes.ConfigRoutes(r)
}
