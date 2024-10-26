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
		IsAuthenticated: true,
	},
	{
		URI:             "/pagamento/qrcode_decode",
		Method:          http.MethodPost,
		Function:        controllers.PagamentoQRCodeDecode,
		IsAuthenticated: true,
	},
	{
		URI:             "/pagamento/qrcode",
		Method:          http.MethodPost,
		Function:        controllers.PagamentoQRCode,
		IsAuthenticated: true,
	},
}
