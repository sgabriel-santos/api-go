package routes

import (
	"api/src/controllers"
	"net/http"
)

var conversoesRoutes = []Rota{
	{
		URI:             "/conversoes",
		Method:          http.MethodGet,
		Function:        controllers.ListarConversoes,
		IsAuthenticated: true,
	},
	{
		URI:             "/converter",
		Method:          http.MethodPost,
		Function:        controllers.ConverterMoeda,
		IsAuthenticated: true,
	},
}
