package routes

import (
	"api/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

// Rota representa todas as rotas da API
type Rota struct {
	URI             string
	Method          string
	Function        func(http.ResponseWriter, *http.Request)
	IsAuthenticated bool
}

// Coloca todas as rotas dentro do router
func ConfigRoutes(r *mux.Router) *mux.Router {
	routes := userRoutes
	routes = append(routes, loginRoutes)

	for _, route := range routes {

		if route.IsAuthenticated {
			r.HandleFunc(
				route.URI,
				middlewares.Logger(middlewares.VerifyAuthentication(route.Function)),
			).Methods(route.Method)
		} else {
			r.HandleFunc(
				route.URI,
				middlewares.Logger(route.Function),
			).Methods(route.Method)
		}

	}

	return r
}
