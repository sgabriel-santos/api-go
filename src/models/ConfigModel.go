package models

type ConfigModel struct {
	ID                int     `json:"id"`
	ClientID          string  `json:"client_id"`
	ClientSecret      string  `json:"client_secret"`
	FeeUSDTBRLConvert float64 `json:"fee_usdtbrl_convert"`
}
