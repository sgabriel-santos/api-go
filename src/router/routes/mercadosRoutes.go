package routes

import (
	"api/src/controllers"
	"net/http"
)

var mercadosRoutes = []Rota{
	{
		URI:             "/mercados",
		Method:          http.MethodGet,
		Function:        controllers.ListarMercados,
		IsAuthenticated: true,
	},
}
