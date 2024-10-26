package routes

import (
	"api/src/controllers"
	"net/http"
)

var depositosRoutes = []Rota{
	{
		URI:             "/depositos/usdt",
		Method:          http.MethodGet,
		Function:        controllers.ListarDepositosUSDT,
		IsAuthenticated: true,
	},
}
