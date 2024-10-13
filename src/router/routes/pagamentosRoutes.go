package routes

import (
	"api/src/controllers"
	"net/http"
)

var pagamentosRoutes = []Rota{
	{
		URI:             "/pagamentos/copiacola",
		Method:          http.MethodGet,
		Function:        controllers.CopiaECola,
		IsAuthenticated: false,
	},
	{
		URI:             "/pagamentos/qrcode_decode",
		Method:          http.MethodPost,
		Function:        controllers.PagamentoQRCodeDecode,
		IsAuthenticated: false,
	},
}
