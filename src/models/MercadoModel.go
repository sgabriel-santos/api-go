package models

// MercadoModel representa um mercado no sistema
type MercadoModel struct {
	ID           uint64  `json:"id"`
	IDSymbol     string  `json:"id_symbol"`
	Symbol       string  `json:"symbol"`
	Base         string  `json:"base"`
	Quote        string  `json:"quote"`
	MinAmount    float64 `json:"min_amount"`
	MinDecimal   string  `json:"min_decimal"`
	DecPrecision string  `json:"dec_precision"`
	MinPrecision float64 `json:"min_precision"`
	Ask          float64 `json:"ask"`
	Bid          float64 `json:"bid"`
	Last         float64 `json:"last"`
}
