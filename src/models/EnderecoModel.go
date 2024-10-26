package models

// EnderecoModel representa um endereço de depósito de criptomoeda
type EnderecoModel struct {
	IDUsuario uint64 `json:"id_usuario"`
	Moeda     string `json:"moeda"`
	Endereco  string `json:"endereco"`
	Rede      string `json:"rede"`
	Dados     string `json:"dados"`
}
