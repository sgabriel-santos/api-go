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
		IsAuthenticated: true,
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
	{
		URI:             "/timezones",
		Method:          http.MethodGet,
		Function:        controllers.ListarTimezones,
		IsAuthenticated: true,
	},
	{
		URI:             "/config_timezone",
		Method:          http.MethodPost,
		Function:        controllers.AtualizarFusoHorario,
		IsAuthenticated: true,
	},
}
