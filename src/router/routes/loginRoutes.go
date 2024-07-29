package routes

import (
	"api/src/controllers"
	"net/http"
)

var loginRoutes = Rota{
	URI:             "/login",
	Method:          http.MethodPost,
	Function:        controllers.Login,
	IsAuthenticated: false,
}
