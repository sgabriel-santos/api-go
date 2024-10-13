package routes

import (
	"api/src/controllers"
	"net/http"
)

var userRoutes = []Rota{
	{
		URI:             "/cadastrar",
		Method:          http.MethodPost,
		Function:        controllers.CreateUser,
		IsAuthenticated: false,
	},
	{
		URI:             "/pagamentos/copiacola",
		Method:          http.MethodGet,
		Function:        controllers.CopiaECola,
		IsAuthenticated: false,
	},
	{
		URI:             "/verifyUser",
		Method:          http.MethodGet,
		Function:        controllers.VerifyUser,
		IsAuthenticated: true,
	},
	{
		URI:             "/users",
		Method:          http.MethodGet,
		Function:        controllers.SearchUserByName,
		IsAuthenticated: true,
	},
}
