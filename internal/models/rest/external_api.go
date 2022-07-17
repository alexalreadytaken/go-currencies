package rest

type CoinPriceRequest struct {
	Currecny string `json:"currency"`
	Code     string `json:"code"`
	Meta     bool   `json:"meta"`
}

type CoinPriceResponse struct {
	Rate float64 `json:"rate"`
}

type FiatPricesResposnse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}
