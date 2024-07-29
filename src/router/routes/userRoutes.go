package routes

import (
	"api/src/controllers"
	"net/http"
)

var userRoutes = []Rota{
	{
		URI:             "/users",
		Method:          http.MethodPost,
		Function:        controllers.CreateUser,
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
	{
		URI:             "/users/{userId}",
		Method:          http.MethodGet,
		Function:        controllers.GetUserById,
		IsAuthenticated: true,
	},
}
